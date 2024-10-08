package game

import (
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

	// Config
	MaxPlayers int32
}

func New(isClient bool) *Game {
	world := ecs.NewWorld()
	space := cp.NewSpace()

	game := &Game{
		World: world,
		Space: space,

		IsClient: isClient,

		Entities: &GameEntities{
			Units: make(map[ecs.Entity]*ecs.Entry),
			Areas: make(map[ecs.Entity]*ecs.Entry),
		},

		NetworkManager: &components.NetworkManagerData{
			NetworkIdToEntityId: make(map[uint32]ecs.Entity),
			EntityIdToNetworkId: make(map[ecs.Entity]uint32),
			NetworkEntities:     make(map[uint32]*components.NetworkEntityData),
			Broadcast:           make(chan []byte, 1),
			World:               world,
			Space:               space,
		},
		lastPlayerID: 0,
		MaxPlayers:   1000,
	}

	// if !game.IsClient {
	// 	game.AddArea()
	// }

	return game
}

func (g *Game) GeneratePlayerId() uint32 {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	g.lastPlayerID++
	return g.lastPlayerID
}

func (g *Game) CreatePlayer(id uint32) *protos.Unit {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	log.Println("CreatePlayer", id)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	unit := &protos.Unit{
		Id: id,
		Position: &protos.Position{
			X: rnd.Float64()*300 + 10,
			Y: rnd.Float64()*220 + 10,
		},
		Frame:  int32(rnd.Intn(4)),
		Action: protos.Action_idle,
		Hp:     100,
	}

	playerEntityId := g.World.Create(components.Transform, components.Physics, components.NetworkEntity, components.Render)
	playerEntity := g.World.Entry(playerEntityId)

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

	g.NetworkManager.Register(playerEntityId, id)

	return unit
}

func (g *Game) RemovePlayer(id uint32) {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	log.Println("remove player: ", id)

	playerEntityId := g.NetworkManager.NetworkIdToEntityId[id]
	playerEntity := g.World.Entry(playerEntityId)

	physics := components.Physics.GetValue(playerEntity)

	g.NetworkManager.Deregister(playerEntityId)

	g.Space.RemoveBody(physics.Body)
	g.World.Remove(playerEntityId)
}

func (g *Game) RemoveAllPlayer() {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	for id, entry := range g.Entities.Units {
		physics := components.Physics.GetValue(entry)

		g.NetworkManager.Deregister(id)

		g.Space.RemoveBody(physics.Body)
		g.World.Remove(id)
	}

	g.Entities.Units = make(map[ecs.Entity]*ecs.Entry)
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
		log.Println("init event. PlayerId: ", event.PlayerId)
		if g.IsClient {
			g.NetworkManager.MyID = &event.PlayerId
		}

	case protos.EventType_move:
		data := event.GetMove()

		entityId := g.NetworkManager.NetworkIdToEntityId[event.PlayerId]
		if g.World.Valid(entityId) == false {
			return
		}

		entityEntry := g.World.Entry(entityId)
		physicsComponent := components.Physics.GetValue(entityEntry)

		physicsComponent.Body.SetVelocity(data.Direction.X*100, -data.Direction.Y*100)

	case protos.EventType_state:
		if !g.IsClient {
			return
		}

		g.RemoveAllPlayer()

		entities := event.GetState().GetEntities()
		if entities != nil {
			for id := range entities {
				g.CreatePlayer(id)
			}
		} else {
			log.Println("No units")
		}

		log.Println("State event", event)

	case protos.EventType_state_patch:
		if !g.IsClient {
			return
		}

		data := event.GetStatePatch()
		if data == nil {
			return
		}

		g.NetworkManager.IncomingPatch = data

		entities := g.NetworkManager.IncomingPatch.GetEntities()
		if entities != nil {
			for _, e := range entities {
				if e.Physics != nil {
					log.Println(e.Physics)
				}
			}
		}

		createdEntities := data.GetCreatedEntities()
		if createdEntities != nil {
			for id := range createdEntities {
				err := g.CreatePlayer(id)
				if err != nil {
					log.Println(err)
				}
			}
		}

		deletedEentities := data.GetDeletedEntities()
		if deletedEentities != nil {
			for id := range deletedEentities {
				g.RemovePlayer(id)
			}
		}

	default:
		log.Println("UNKNOWN EVENT: ", event)
	}
}

func (g *Game) HandleEvents() {
	for _, event := range g.NetworkManager.UnhandledEvents {
		g.HandleEvent(event)
	}

	g.NetworkManager.UnhandledEvents = make([]*protos.Event, 0)
}

const (
	patchRate     = time.Second / 2
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
			dt := time.Now().Sub(lastEvolveTime).Seconds()
			g.Update(dt)
			lastEvolveTime = time.Now()

		// case <-patchTicker.C:
		// 	g.Mx.Lock()
		// 	g.NetworkManager.SendPatch()
		// 	g.Mx.Unlock()

		case <-lazyPatchTicker.C:
			// g.Broadcast <- *g.StateSerialized
		}
	}
}

func (g *Game) CacheGameState() error {
	if !g.IsClient {
		cachedEntities := make(map[uint32]*protos.NetworkEntity, len(g.NetworkManager.NetworkEntities))
		for key, value := range g.NetworkManager.NetworkEntities {
			cachedEntities[key] = &protos.NetworkEntity{
				Id:        value.Id,
				Transform: value.Transform,
				Physics:   value.Physics,
				Skin:      *value.Skin,
			}
		}

		stateEvent := &protos.Event{
			Type: protos.EventType_state,
			Data: &protos.Event_State{
				State: &protos.GameState{
					Entities: cachedEntities,
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

func (g *Game) Update(dt float64) {
	g.Mx.Lock()
	defer g.Mx.Unlock()

	g.Space.Step(dt)
	components.Physics.Each(g.World, func(e *ecs.Entry) {
		err := components.Physics.GetValue(e).Update(dt, e, g.IsClient)
		if err != nil {
			panic(err)
		}
	})

	if !g.IsClient {
		g.HandleEvents()
	}

	g.NetworkManager.Update(dt, g.IsClient)

	if !g.IsClient {
		err := g.CacheGameState()
		if err != nil {
			panic(err)
		}
	}
}
