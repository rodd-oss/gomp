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
	Mx              sync.Mutex
	IsServer        bool
	Units           map[uint32]*protos.Unit
	PatchedUnits    map[uint32]*protos.PatchUnit
	UnitsSerialized *[]byte
	MyID            uint32
	UnhandledEvents []*protos.Event
	Broadcast       chan []byte
	lastPlayerID    uint32
}

func New(isServer bool, units map[uint32]*protos.Unit) *Game {
	world := &Game{
		IsServer:     isServer,
		Units:        units,
		PatchedUnits: make(map[uint32]*protos.PatchUnit),
		Broadcast:    make(chan []byte, 1),
		lastPlayerID: 0,
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
	}
	world.Units[id] = unit

	return unit
}

func (world *Game) RemovePlayer(unit *protos.Unit) {
	world.Mx.Lock()
	defer world.Mx.Unlock()

	delete(world.Units, unit.Id)
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
	case protos.EventType_connect:
		data := event.GetConnect()
		world.Units[data.Unit.Id] = data.Unit

	case protos.EventType_init:
		if world.IsServer {
			world.MyID = event.PlayerId
		}

	case protos.EventType_exit:
		delete(world.Units, event.PlayerId)

	case protos.EventType_move:
		data := event.GetMove()
		unit := world.Units[event.PlayerId]
		if unit == nil {
			return
		}
		unit.Action = protos.Action_run
		unit.Velocity.Direction = data.Direction

		if !world.IsServer {
			if event.PlayerId == unit.Id {
				return
			}

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

		if !world.IsServer {
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
		if !world.IsServer {
			return
		}

		data := event.GetState()
		units := data.GetUnits()
		if units != nil {
			world.Units = units
		}

	case protos.EventType_state_patch:
		if !world.IsServer {
			return
		}

		data := event.GetStatePatch()
		units := data.GetUnits()
		if units != nil {
			for _, unit := range units {
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

	default:
		log.Println("UNKNOWN EVENT: ", event)
	}
}

func (world *Game) ProccessEvents() error {
	world.Mx.Lock()
	defer world.Mx.Unlock()

	for _, event := range world.UnhandledEvents {
		world.HandleEvent(event)
	}

	world.UnhandledEvents = make([]*protos.Event, 0)

	return nil
}

const (
	patchRate     = time.Second / 20
	lazyPatchRate = time.Second * 30
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
			world.ProccessEvents()

			dt := time.Now().Sub(lastEvolveTime).Seconds()
			world.HandlePhysics(dt)
			lastEvolveTime = time.Now()

			if world.IsServer == false {
				world.Mx.Lock()
				cachedUnits := make(map[uint32]*protos.Unit, len(world.Units))
				for key, value := range world.Units {
					cachedUnits[key] = value
				}
				world.Mx.Unlock()

				stateEvent := &protos.Event{
					Type: protos.EventType_state,
					Data: &protos.Event_State{
						State: &protos.GameState{
							Units: cachedUnits,
						},
					},
				}
				s, err := proto.Marshal(stateEvent)
				if err != nil {
					continue
				}

				world.UnitsSerialized = &s
			}
		case <-patchTicker.C:
			if len(world.PatchedUnits) == 0 {
				continue
			}
			statePatchEvent := &protos.Event{
				Type: protos.EventType_state_patch,
				Data: &protos.Event_StatePatch{
					StatePatch: &protos.GameStatePatche{
						Units: world.PatchedUnits,
					},
				},
			}

			m, err := proto.Marshal(statePatchEvent)
			if err != nil {
				continue
			}
			world.Broadcast <- m
			world.PatchedUnits = make(map[uint32]*protos.PatchUnit)

		case <-lazyPatchTicker.C:
			world.Broadcast <- *world.UnitsSerialized
		}
	}
}

func (world *Game) HandlePhysics(dt float64) {
	world.Mx.Lock()
	defer world.Mx.Unlock()

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
}
