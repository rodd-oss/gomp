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
	Replica         bool
	Units           map[uint32]*protos.Unit
	PatchedUnits    map[uint32]*protos.PatchUnit
	UnitsSerialized *[]byte
	MyID            uint32
	UnhandledEvents []*protos.Event
	Broadcast       chan []byte
	lastPlayerID    uint32
}

func New(isReplica bool, units map[uint32]*protos.Unit) *Game {
	world := &Game{
		Replica:      isReplica,
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

	skins := []string{"big_demon", "big_zombie", "elf_f"}
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
		Skin:   skins[rnd.Intn(len(skins))],
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
		data := event.GetInit()

		if world.Replica {
			world.MyID = data.PlayerId
		}

	case protos.EventType_exit:
		data := event.GetExit()
		delete(world.Units, data.PlayerId)

	case protos.EventType_move:
		data := event.GetMove()
		unit := world.Units[data.PlayerId]
		if unit == nil {
			return
		}
		unit.Action = protos.Action_run
		unit.Velocity.Direction = data.Direction

		if !world.Replica {
			world.PatchedUnits[data.PlayerId] = &protos.PatchUnit{
				Id:     data.PlayerId,
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
		data := event.GetStop()
		unit := world.Units[data.PlayerId]
		if unit == nil {
			return
		}
		unit.Action = protos.Action_idle

		if !world.Replica {
			world.PatchedUnits[data.PlayerId] = &protos.PatchUnit{
				Id:     data.PlayerId,
				Action: &unit.Action,
				Position: &protos.Position{
					X: unit.Position.X,
					Y: unit.Position.Y,
				},
			}
		}

	case protos.EventType_state:
		data := event.GetState()
		units := data.GetUnits()
		if units != nil {
			world.Units = units
		}

	case protos.EventType_state_patch:
		data := event.GetStatePatch()
		units := data.GetUnits()
		if units != nil {
			for _, unit := range units {
				if unit.Action != nil {
					world.Units[unit.Id].Action = *unit.Action
				}

				if unit.Velocity != nil {
					world.Units[unit.Id].Velocity = unit.Velocity
				}

				if unit.Position != nil {
					world.Units[unit.Id].Position = unit.Position
				}

				if unit.Side != nil {
					world.Units[unit.Id].Side = *unit.Side
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

			if world.Replica == false {
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
