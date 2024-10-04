package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
	"tomb_mates/internal/components"
	"tomb_mates/internal/protos"

	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"google.golang.org/protobuf/proto"
)

type GameEntities struct {
	Areas map[ecs.Entity]*ecs.Entry
	Units map[ecs.Entity]*ecs.Entry
}

// Game represents game state
type Game struct {
	Mx sync.Mutex

	World ecs.World // main world of entities
	Space *cp.Space // main physical space

	Entities       *GameEntities
	NetworkManager *components.NetworkManagerData

	// Client
	IsClient bool

	// Server
	StateSerialized *[]byte
	lastPlayerID    uint32
	lastAreaID      uint32
	Broadcast       chan []byte

	// Config
	MaxPlayers int32
}

func New(isClient bool) *Game {
	game := &Game{
		World: ecs.NewWorld(),
		Space: cp.NewSpace(),

		IsClient: isClient,

		Entities: &GameEntities{
			Units: make(map[ecs.Entity]*ecs.Entry),
			Areas: make(map[ecs.Entity]*ecs.Entry),
		},

		NetworkManager: &components.NetworkManagerData{
			Units:           make(map[uint32]*protos.Unit),
			UnitToEntity:    make(map[uint32]ecs.Entity),
			PatchedUnits:    make(map[uint32]*protos.PatchUnit),
			CreatedUnits:    make(map[uint32]*protos.Unit),
			DeletedUnitsIds: make(map[uint32]*protos.Empty),

			Areas:           make(map[uint32]*protos.Area),
			AreaToEntity:    make(map[uint32]ecs.Entity),
			PatchedAreas:    make(map[uint32]*protos.PatchArea),
			CreatedAreas:    make(map[uint32]*protos.Area),
			DeletedAreasIds: make(map[uint32]*protos.Empty),
		},

		Broadcast:    make(chan []byte, 1),
		lastPlayerID: 0,
		lastAreaID:   0,
		MaxPlayers:   1000,
	}

	// if !game.IsClient {
	// 	game.AddArea()
	// }

	return game
}

func (g *Game) CreatePlayer() *protos.Unit {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	id := g.lastPlayerID
	g.lastPlayerID++
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	unit := &protos.Unit{
		Id: id,
		Position: &protos.Position{
			X: rnd.Float64()*300 + 10,
			Y: rnd.Float64()*220 + 10,
		},
		Frame:  int32(rnd.Intn(4)),
		Skin:   protos.Skin(rnd.Intn(len(protos.Skin_name))),
		Action: protos.Action_idle,
		Velocity: &protos.Velocity{
			Direction: protos.Direction_left,
			Speed:     100,
		},
		Hp: 100,
	}
	delete(g.NetworkManager.DeletedUnitsIds, unit.Id)

	playerEntityId := g.World.Create(components.Transform, components.NetworkUnit, components.Physics, components.Render)
	playerEntity := g.World.Entry(playerEntityId)

	components.NetworkUnit.SetValue(playerEntity, components.NetworkUnitData{
		Unit: unit,
	})

	components.Transform.SetValue(playerEntity, components.TransformData{
		LocalPosition: math.Vec2{X: 1, Y: 2},
		LocalRotation: 0,
		LocalScale: math.Vec2{
			X: 1,
			Y: 1,
		},
	})

	body := cp.NewKinematicBody()
	body = g.Space.AddBody(body)
	body.SetPosition(cp.Vector{
		X: unit.Position.X,
		Y: unit.Position.Y,
	})

	shape := g.Space.AddShape(cp.NewCircle(body, 8, cp.Vector{}))
	shape.SetElasticity(0)
	shape.SetFriction(0)

	components.Physics.SetValue(playerEntity, components.PhysicsData{
		Body: body,
	})

	g.Entities.Units[playerEntityId] = playerEntity
	g.NetworkManager.UnitToEntity[unit.Id] = playerEntityId
	g.NetworkManager.Units[unit.Id] = unit
	g.NetworkManager.CreatedUnits[id] = unit

	return unit
}

func (g *Game) InsertPlayer(unit *protos.Unit) error {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	delete(g.NetworkManager.DeletedUnitsIds, unit.Id)

	playerEntityId := g.World.Create(components.Transform, components.NetworkUnit, components.Physics, components.Render)
	playerEntity := g.World.Entry(playerEntityId)

	components.NetworkUnit.SetValue(playerEntity, components.NetworkUnitData{
		Unit: unit,
	})

	components.Transform.SetValue(playerEntity, components.TransformData{
		LocalPosition: math.Vec2{X: unit.Position.X, Y: unit.Position.Y},
		LocalRotation: 0,
		LocalScale: math.Vec2{
			X: 1,
			Y: 1,
		},
	})

	body := cp.NewKinematicBody()
	body = g.Space.AddBody(body)
	body.SetPosition(cp.Vector{
		X: unit.Position.X,
		Y: unit.Position.Y,
	})

	shape := g.Space.AddShape(cp.NewCircle(body, 8, cp.Vector{}))
	shape.SetElasticity(0)
	shape.SetFriction(0)

	components.Physics.SetValue(playerEntity, components.PhysicsData{
		Body: body,
	})

	g.Entities.Units[playerEntityId] = playerEntity
	g.NetworkManager.Units[unit.Id] = unit
	g.NetworkManager.UnitToEntity[unit.Id] = playerEntityId

	return nil
}

func (g *Game) RemovePlayer(id uint32) {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	playerEntityId := g.NetworkManager.UnitToEntity[id]
	playerEntity := g.World.Entry(playerEntityId)
	physics := components.Physics.GetValue(playerEntity)

	g.Space.RemoveBody(physics.Body)
	delete(g.NetworkManager.CreatedUnits, id)
	delete(g.NetworkManager.UnitToEntity, id)
	delete(g.NetworkManager.Units, id)
	g.NetworkManager.DeletedUnitsIds[id] = &protos.Empty{}

	g.World.Remove(playerEntityId)
}

func (g *Game) RemoveAllPlayer() {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	for _, playerEntityId := range g.NetworkManager.UnitToEntity {
		g.World.Remove(playerEntityId)
	}

	g.Entities.Units = make(map[ecs.Entity]*ecs.Entry)
	g.NetworkManager.UnitToEntity = make(map[uint32]ecs.Entity)
}

func (g *Game) AddArea() *protos.Area {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	id := g.lastAreaID
	g.lastAreaID++

	area := &protos.Area{
		Id: id,
		Position: &protos.Position{
			X: 0,
			Y: 0,
		},
		Size: &protos.Vector2{
			X: 100,
			Y: 100,
		},
		Skin:            "area",
		Frame:           0,
		AffectedUnitIds: make(map[uint32]*protos.Empty),
	}

	delete(g.NetworkManager.DeletedAreasIds, area.Id)
	g.NetworkManager.CreatedAreas[id] = area

	areaEntityId := g.World.Create(components.Transform, components.NetworkArea, components.Physics)
	areaEntity := g.World.Entry(areaEntityId)

	components.NetworkArea.SetValue(areaEntity, components.NetworkAreaData{
		Area: area,
	})

	components.Transform.SetValue(areaEntity, components.TransformData{
		LocalPosition: math.Vec2{X: area.Position.X, Y: 2},
		LocalRotation: 0,
		LocalScale: math.Vec2{
			X: 1,
			Y: 1,
		},
	})

	body := cp.NewKinematicBody()
	body = g.Space.AddBody(body)
	body.SetPosition(cp.Vector{
		X: area.Position.X,
		Y: area.Position.Y,
	})

	shape := g.Space.AddShape(cp.NewCircle(body, 8, cp.Vector{}))
	shape.SetElasticity(0)
	shape.SetFriction(0)

	components.Physics.SetValue(areaEntity, components.PhysicsData{
		Body: body,
	})

	g.Entities.Areas[areaEntityId] = areaEntity
	g.NetworkManager.AreaToEntity[area.Id] = areaEntityId

	return area
}

func (g *Game) RemoveArea(area *protos.Area) {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	areaEntityId := g.NetworkManager.AreaToEntity[area.Id]

	g.NetworkManager.DeletedAreasIds[area.Id] = &protos.Empty{}
	delete(g.NetworkManager.CreatedUnits, area.Id)
	g.World.Remove(areaEntityId)
	delete(g.NetworkManager.UnitToEntity, area.Id)
}

func (g *Game) InsertArea(area *protos.Area) error {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	delete(g.NetworkManager.DeletedAreasIds, area.Id)

	areaEntityId := g.World.Create(components.Transform, components.NetworkArea, components.Physics)
	areaEntity := g.World.Entry(areaEntityId)

	components.NetworkArea.SetValue(areaEntity, components.NetworkAreaData{
		Area: area,
	})

	components.Transform.SetValue(areaEntity, components.TransformData{
		LocalPosition: math.Vec2{X: area.Position.X, Y: area.Position.Y},
		LocalRotation: 0,
		LocalScale: math.Vec2{
			X: 1,
			Y: 1,
		},
	})

	body := cp.NewKinematicBody()

	body = g.Space.AddBody(body)
	shape := g.Space.AddShape(cp.NewBox(body, area.Size.X, area.Size.Y, 0))

	body.SetPosition(cp.Vector{
		X: area.Position.X,
		Y: area.Position.Y,
	})
	shape.SetElasticity(0)
	shape.SetFriction(0)

	components.Physics.SetValue(areaEntity, components.PhysicsData{
		Body: body,
	})

	g.Entities.Units[areaEntityId] = areaEntity
	g.NetworkManager.UnitToEntity[area.Id] = areaEntityId

	return nil
}

func (g *Game) RemoveAllAreas() {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	for _, areaEntityId := range g.NetworkManager.AreaToEntity {
		g.World.Remove(areaEntityId)
	}

	g.Entities.Areas = make(map[ecs.Entity]*ecs.Entry)
	g.NetworkManager.UnitToEntity = make(map[uint32]ecs.Entity)
}

func (g *Game) RegisterEvent(event *protos.Event) {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	g.NetworkManager.UnhandledEvents = append(g.NetworkManager.UnhandledEvents, event)
}

func (g *Game) HandleEvent(event *protos.Event) {
	if event == nil {
		return
	}

	etype := event.GetType()
	switch etype {
	case protos.EventType_init:
		fmt.Println("init event")
		if g.IsClient {
			g.NetworkManager.MyID = &event.PlayerId
		}

	case protos.EventType_move:
		data := event.GetMove()

		unitId := g.NetworkManager.UnitToEntity[event.PlayerId]
		if g.World.Valid(unitId) == false {
			return
		}

		unitEntry := g.World.Entry(unitId)
		unitComponent := components.NetworkUnit.GetValue(unitEntry)
		unit := unitComponent.Unit

		unit.Action = protos.Action_run
		unit.Velocity.Direction = data.Direction

		components.NetworkUnit.SetValue(unitEntry, unitComponent)
		physicsComponent := components.Physics.GetValue(unitEntry)
		switch unit.Velocity.Direction {
		case protos.Direction_left:
			physicsComponent.Body.SetVelocity(-unit.Velocity.Speed, 0)
		case protos.Direction_right:
			physicsComponent.Body.SetVelocity(unit.Velocity.Speed, 0)
		case protos.Direction_up:
			physicsComponent.Body.SetVelocity(0, -unit.Velocity.Speed)
		case protos.Direction_down:
			physicsComponent.Body.SetVelocity(0, unit.Velocity.Speed)
		}
		if !g.IsClient {

			g.NetworkManager.PatchedUnits[event.PlayerId] = &protos.PatchUnit{
				Id:     event.PlayerId,
				Action: &unit.Action,
				Velocity: &protos.Velocity{
					Direction: unit.Velocity.Direction,
					Speed:     unit.Velocity.Speed,
				},
				Position: &protos.Position{
					X: unit.Position.X,
					Y: unit.Position.Y,
				},
			}
		}

	case protos.EventType_stop:
		unitId := g.NetworkManager.UnitToEntity[event.PlayerId]
		if g.World.Valid(unitId) == false {
			return
		}

		unitEntry := g.World.Entry(unitId)
		unitComponent := components.NetworkUnit.GetValue(unitEntry)
		unit := unitComponent.Unit

		unit.Action = protos.Action_idle

		components.NetworkUnit.SetValue(unitEntry, unitComponent)
		physicsComponent := components.Physics.GetValue(unitEntry)
		physicsComponent.Body.SetVelocity(0, 0)

		if !g.IsClient {
			g.NetworkManager.PatchedUnits[event.PlayerId] = &protos.PatchUnit{
				Id:     event.PlayerId,
				Action: &unit.Action,
				Position: &protos.Position{
					X: unit.Position.X,
					Y: unit.Position.Y,
				},
				Velocity: &protos.Velocity{
					Direction: unit.Velocity.Direction,
					Speed:     0,
				},
			}
		}

	case protos.EventType_state:
		if !g.IsClient {
			return
		}

		g.RemoveAllPlayer()
		g.RemoveAllAreas()

		units := event.GetState().GetUnits()
		if units != nil {
			for _, unit := range units {
				err := g.InsertPlayer(unit)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			log.Println("No units")
		}

		areas := event.GetState().GetAreas()
		if areas != nil {
			for _, area := range areas {
				err := g.InsertArea(area)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			log.Println("No areas")
		}

		fmt.Println("State event")

	case protos.EventType_state_patch:
		if !g.IsClient {
			return
		}

		data := event.GetStatePatch()
		units := data.GetUnits()
		if units != nil {
			for _, unit := range units {
				if unit == nil {
					continue
				}

				unitId := g.NetworkManager.UnitToEntity[unit.Id]
				if g.World.Valid(unitId) == false {
					return
				}

				unitEntry := g.World.Entry(unitId)
				unitComponent := components.NetworkUnit.GetValue(unitEntry)
				wu := unitComponent.Unit

				physicsComponent := components.Physics.GetValue(unitEntry)

				if unit.Action != nil {
					wu.Action = *unit.Action
				}

				if unit.Velocity != nil {
					wu.Velocity = unit.Velocity
					switch unit.Velocity.Direction {
					case protos.Direction_left:
						physicsComponent.Body.SetVelocity(-unit.Velocity.Speed, 0)
					case protos.Direction_right:
						physicsComponent.Body.SetVelocity(unit.Velocity.Speed, 0)
					case protos.Direction_up:
						physicsComponent.Body.SetVelocity(0, -unit.Velocity.Speed)
					case protos.Direction_down:
						physicsComponent.Body.SetVelocity(0, unit.Velocity.Speed)
					default:
						physicsComponent.Body.SetVelocity(0, 0)
					}
				}

				if unit.Position != nil {
					wu.Position = unit.Position
					physicsComponent.Body.SetPosition(cp.Vector{
						X: unit.Position.X,
						Y: unit.Position.Y,
					})
				}

				if unit.Side != nil {
					wu.Side = *unit.Side
				}

				components.NetworkUnit.SetValue(unitEntry, unitComponent)
			}
		}

		createdUnits := data.GetCreatedUnits()
		if createdUnits != nil {
			for _, unit := range createdUnits {
				g.InsertPlayer(unit)
			}
		}

		deletedUnitsIds := data.GetDeletedUnitsIds()
		for id := range deletedUnitsIds {
			g.RemovePlayer(id)
		}

	default:
		log.Println("UNKNOWN EVENT: ", event)
	}
}

func (g *Game) HandleEvents() {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	for _, event := range g.NetworkManager.UnhandledEvents {
		g.HandleEvent(event)
	}

	g.NetworkManager.UnhandledEvents = make([]*protos.Event, 0)
}

const (
	patchRate     = time.Second / 20
	lazyPatchRate = time.Minute * 5
)

func (g *Game) Run(tickRate time.Duration) {
	ticker := time.NewTicker(tickRate)
	lastEvolveTime := time.Now()

	patchTicker := time.NewTicker(patchRate)
	defer patchTicker.Stop()

	lazyPatchTicker := time.NewTicker(lazyPatchRate)
	defer patchTicker.Stop()

	for {
		select {
		case <-ticker.C:
			g.Update(lastEvolveTime)
			lastEvolveTime = time.Now()

		case <-patchTicker.C:
			g.SendPatch()

		case <-lazyPatchTicker.C:
			g.Broadcast <- *g.StateSerialized
		}
	}
}

func (g *Game) CacheGameState() error {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	if !g.IsClient {
		cachedUnits := make(map[uint32]*protos.Unit, len(g.NetworkManager.Units))
		for key, value := range g.NetworkManager.Units {
			cachedUnits[key] = value
		}

		cachedAreas := make(map[uint32]*protos.Area, len(g.NetworkManager.Areas))
		for key, value := range g.NetworkManager.Areas {
			cachedAreas[key] = value
		}

		stateEvent := &protos.Event{
			Type: protos.EventType_state,
			Data: &protos.Event_State{
				State: &protos.GameState{
					Units: cachedUnits,
					Areas: cachedAreas,
				},
			},
		}
		s, err := proto.Marshal(stateEvent)
		if err != nil {
			return err
		}

		g.StateSerialized = &s
	}

	return nil
}

func (g *Game) Update(lastUpdateAt time.Time) {
	dt := time.Now().Sub(lastUpdateAt).Seconds()

	g.HandlePhysics(dt)

	g.HandleEvents()

	err := g.CacheGameState()
	if err != nil {
		fmt.Println(err)
	}

	return
}

func (g *Game) SendPatch() {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	if len(g.NetworkManager.PatchedUnits) == 0 && len(g.NetworkManager.CreatedUnits) == 0 && len(g.NetworkManager.DeletedUnitsIds) == 0 {
		return
	}

	statePatchEvent := &protos.Event{
		Type: protos.EventType_state_patch,
		Data: &protos.Event_StatePatch{
			StatePatch: &protos.GameStatePatche{
				Units:           g.NetworkManager.PatchedUnits,
				CreatedUnits:    g.NetworkManager.CreatedUnits,
				DeletedUnitsIds: g.NetworkManager.DeletedUnitsIds,
			},
		},
	}

	m, err := proto.Marshal(statePatchEvent)
	if err != nil {
		return
	}
	g.Broadcast <- m
	g.NetworkManager.PatchedUnits = make(map[uint32]*protos.PatchUnit)
	g.NetworkManager.CreatedUnits = make(map[uint32]*protos.Unit)
	g.NetworkManager.DeletedUnitsIds = make(map[uint32]*protos.Empty)
}

func (g *Game) HandlePhysics(dt float64) {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	g.Space.Step(dt)

	components.Physics.Each(g.World, func(e *ecs.Entry) {
		err := components.Physics.GetValue(e).Update(dt, e)
		if err != nil {
			fmt.Println(err)
		}
	})
}
