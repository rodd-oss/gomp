package game

import (
	"log"
	"math/rand"
	"sync"
	"time"
	"tomb_mates/internal/protos"

	"google.golang.org/protobuf/proto"
)

// Game represents game state
type Game struct {
	Mx       sync.Mutex
	IsClient bool

	Areas           map[uint32]*protos.Area
	PatchedAreas    map[uint32]*protos.PatchArea
	DeletedAreasIds map[uint32]*protos.Empty
	CreatedAreas    map[uint32]*protos.Area

	Units           map[uint32]*protos.Unit
	PatchedUnits    map[uint32]*protos.PatchUnit
	DeletedUnitsIds map[uint32]*protos.Empty
	CreatedUnits    map[uint32]*protos.Unit

	StateSerialized *[]byte

	MyID            uint32
	UnhandledEvents []*protos.Event
	Broadcast       chan []byte
	lastPlayerID    uint32
	lastAreaID      uint32
	MaxPlayers      int32
}

func New(isClient bool, units map[uint32]*protos.Unit) *Game {
	world := &Game{
		IsClient:        isClient,
		Areas:           make(map[uint32]*protos.Area),
		Units:           units,
		PatchedUnits:    make(map[uint32]*protos.PatchUnit),
		CreatedUnits:    make(map[uint32]*protos.Unit),
		DeletedUnitsIds: make(map[uint32]*protos.Empty),
		Broadcast:       make(chan []byte, 1),
		lastPlayerID:    0,
		lastAreaID:      0,
		MaxPlayers:      1000,
	}

	if !world.IsClient {
		world.AddArea()
	}

	return world
}

func (world *Game) AddPlayer() *protos.Unit {
	world.Mx.Lock()
	defer world.Mx.Unlock()

	id := world.lastPlayerID
	world.lastPlayerID++
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
	delete(world.DeletedUnitsIds, unit.Id)
	world.CreatedUnits[id] = unit
	world.Units[id] = unit

	return unit
}

func (world *Game) RemovePlayer(unit *protos.Unit) {
	world.Mx.Lock()
	defer world.Mx.Unlock()

	world.DeletedUnitsIds[unit.Id] = &protos.Empty{}
	delete(world.CreatedUnits, unit.Id)
	delete(world.Units, unit.Id)
}

func (world *Game) AddArea() *protos.Area {
	world.Mx.Lock()
	defer world.Mx.Unlock()

	id := world.lastAreaID
	world.lastAreaID++

	aoe := &protos.Area{
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

	world.Areas[id] = aoe

	return aoe
}

func (world *Game) RemoveArea(area *protos.Area) {
	world.Mx.Lock()
	defer world.Mx.Unlock()

	world.DeletedUnitsIds[area.Id] = &protos.Empty{}
	delete(world.CreatedUnits, area.Id)
	delete(world.Units, area.Id)
}

func (world *Game) RegisterEvent(event *protos.Event) {
	world.Mx.Lock()
	defer world.Mx.Unlock()

	world.UnhandledEvents = append(world.UnhandledEvents, event)
}

func (world *Game) HandleEvent(event *protos.Event) {
	if event == nil {
		return
	}

	etype := event.GetType()
	switch etype {
	case protos.EventType_init:
		if world.IsClient {
			world.MyID = event.PlayerId
		}

	case protos.EventType_move:
		data := event.GetMove()
		unit := world.Units[event.PlayerId]
		if unit == nil {
			return
		}
		unit.Action = protos.Action_run
		unit.Velocity.Direction = data.Direction

		if !world.IsClient {
			world.PatchedUnits[event.PlayerId] = &protos.PatchUnit{
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
		unit := world.Units[event.PlayerId]
		if unit == nil {
			return
		}
		unit.Action = protos.Action_idle

		if !world.IsClient {
			world.PatchedUnits[event.PlayerId] = &protos.PatchUnit{
				Id:     event.PlayerId,
				Action: &unit.Action,
				Position: &protos.Position{
					X: unit.Position.X,
					Y: unit.Position.Y,
				},
			}
		}

	case protos.EventType_state:
		if !world.IsClient {
			return
		}

		data := event.GetState()
		units := data.GetUnits()
		if units != nil {
			world.Units = units
		}
		areas := data.GetAreas()
		if areas != nil {
			world.Areas = areas
		}

	case protos.EventType_state_patch:
		if !world.IsClient {
			return
		}

		data := event.GetStatePatch()
		units := data.GetUnits()
		if units != nil {
			for _, unit := range units {
				if unit == nil {
					continue
				}

				wu := world.Units[unit.Id]
				if wu == nil {
					continue
				}

				if unit.Action != nil {
					wu.Action = *unit.Action
				}

				if unit.Velocity != nil {
					wu.Velocity = unit.Velocity
				}

				if unit.Position != nil {
					wu.Position = unit.Position
				}

				if unit.Side != nil {
					wu.Side = *unit.Side
				}
			}
		}
		createdUnits := data.GetCreatedUnits()
		if createdUnits != nil {
			for _, unit := range createdUnits {
				world.Units[unit.Id] = unit
			}
		}
		deletedUnitsIds := data.GetDeletedUnitsIds()
		for id := range deletedUnitsIds {
			delete(world.Units, id)
		}

	default:
		log.Println("UNKNOWN EVENT: ", event)
	}
}

const (
	patchRate     = time.Second / 20
	lazyPatchRate = time.Minute * 5
)

func (world *Game) Run(tickRate time.Duration) {
	ticker := time.NewTicker(tickRate)
	lastEvolveTime := time.Now()

	patchTicker := time.NewTicker(patchRate)
	defer patchTicker.Stop()

	lazyPatchTicker := time.NewTicker(lazyPatchRate)
	defer patchTicker.Stop()

	for {
		select {
		case <-ticker.C:
			world.Update(lastEvolveTime)
			lastEvolveTime = time.Now()

		case <-patchTicker.C:
			world.SendPatch()

		case <-lazyPatchTicker.C:
			world.Broadcast <- *world.StateSerialized
		}
	}
}

func (world *Game) Update(lastUpdateAt time.Time) {
	world.Mx.Lock()
	defer world.Mx.Unlock()

	for _, event := range world.UnhandledEvents {
		world.HandleEvent(event)
	}

	world.UnhandledEvents = make([]*protos.Event, 0)

	dt := time.Now().Sub(lastUpdateAt).Seconds()
	world.HandlePhysics(dt)

	if world.IsClient == false {
		cachedUnits := make(map[uint32]*protos.Unit, len(world.Units))
		for key, value := range world.Units {
			cachedUnits[key] = value
		}

		cachedAreas := make(map[uint32]*protos.Area, len(world.Areas))
		for key, value := range world.Areas {
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
			return
		}

		world.StateSerialized = &s
	}

	return
}

func (world *Game) SendPatch() {
	world.Mx.Lock()
	defer world.Mx.Unlock()

	if len(world.PatchedUnits) == 0 && len(world.CreatedUnits) == 0 && len(world.DeletedUnitsIds) == 0 {
		return
	}

	statePatchEvent := &protos.Event{
		Type: protos.EventType_state_patch,
		Data: &protos.Event_StatePatch{
			StatePatch: &protos.GameStatePatche{
				Units:           world.PatchedUnits,
				CreatedUnits:    world.CreatedUnits,
				DeletedUnitsIds: world.DeletedUnitsIds,
			},
		},
	}

	m, err := proto.Marshal(statePatchEvent)
	if err != nil {
		return
	}
	world.Broadcast <- m
	world.PatchedUnits = make(map[uint32]*protos.PatchUnit)
	world.CreatedUnits = make(map[uint32]*protos.Unit)
	world.DeletedUnitsIds = make(map[uint32]*protos.Empty)
}

func (world *Game) HandlePhysics(dt float64) {
	if world.IsClient {
		world.Mx.Lock()
		defer world.Mx.Unlock()
	}

	for i := range world.Units {
		if world.Units[i].Action == protos.Action_run {
			switch world.Units[i].Velocity.Direction {
			case protos.Direction_left:
				world.Units[i].Position.X -= world.Units[i].Velocity.Speed * dt
				world.Units[i].Side = protos.Direction_left
			case protos.Direction_right:
				world.Units[i].Position.X += world.Units[i].Velocity.Speed * dt
				world.Units[i].Side = protos.Direction_right
			case protos.Direction_up:
				world.Units[i].Position.Y -= world.Units[i].Velocity.Speed * dt
			case protos.Direction_down:
				world.Units[i].Position.Y += world.Units[i].Velocity.Speed * dt
			default:
				log.Println("UNKNOWN DIRECTION: ", world.Units[i].Velocity.Direction)
			}
		}
	}
	// for  := range world.Areas {

	// }
}
