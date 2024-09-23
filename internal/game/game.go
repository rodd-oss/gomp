package game

import (
	"log"
	"math/rand"
	"sync"
	"time"
	"tomb_mates/internal/protos"

	uuid "github.com/satori/go.uuid"
)

// Game represents game state
type Game struct {
	Mx      sync.Mutex
	Replica bool
	Units   map[string]*protos.Unit
	MyID    string
}

func New(isReplica bool, units map[string]*protos.Unit, tickRate time.Duration) *Game {
	world := &Game{
		Replica: isReplica,
		Units:   units,
	}

	go world.evolve(tickRate)

	return world
}

func (world *Game) AddPlayer() string {
	skins := []string{"big_demon", "big_zombie", "elf_f"}
	id := uuid.NewV4().String()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	unit := &protos.Unit{
		Id:     id,
		X:      rnd.Float64()*300 + 10,
		Y:      rnd.Float64()*220 + 10,
		Frame:  int32(rnd.Intn(4)),
		Skin:   skins[rnd.Intn(len(skins))],
		Action: "idle",
		Speed:  200,
	}
	world.Units[id] = unit

	return id
}

func (world *Game) HandleEvent(event *protos.Event) {
	world.Mx.Lock()
	defer world.Mx.Unlock()

	switch event.GetType() {
	case protos.Event_type_connect:
		data := event.GetConnect()
		world.Units[data.Unit.Id] = data.Unit

	case protos.Event_type_init:
		data := event.GetInit()
		if world.Replica {
			world.MyID = data.PlayerId
			world.Units = data.Units
		}

	case protos.Event_type_exit:
		data := event.GetExit()
		delete(world.Units, data.PlayerId)

	case protos.Event_type_move:
		data := event.GetMove()
		unit := world.Units[data.PlayerId]
		if unit == nil {
			return
		}
		unit.Action = UnitActionMove
		unit.Direction = data.Direction

	case protos.Event_type_idle:
		data := event.GetIdle()
		unit := world.Units[data.PlayerId]
		if unit == nil {
			return
		}
		unit.Action = UnitActionIdle

	default:
		log.Println("UNKNOWN EVENT: ", event)
	}
}

func (world *Game) evolve(tickRate time.Duration) {
	ticker := time.NewTicker(tickRate)
	lastEvolveTime := time.Now()

	for {
		select {
		case <-ticker.C:
			dt := time.Now().Sub(lastEvolveTime).Seconds()
			for i := range world.Units {
				if world.Units[i].Action == UnitActionMove {
					switch world.Units[i].Direction {
					case protos.Direction_left:
						world.Units[i].X -= world.Units[i].Speed * dt
						world.Units[i].Side = protos.Direction_left
					case protos.Direction_right:
						world.Units[i].X += world.Units[i].Speed * dt
						world.Units[i].Side = protos.Direction_right
					case protos.Direction_up:
						world.Units[i].Y -= world.Units[i].Speed * dt
					case protos.Direction_down:
						world.Units[i].Y += world.Units[i].Speed * dt
					default:
						log.Println("UNKNOWN DIRECTION: ", world.Units[i].Direction)
					}
				}
			}
			lastEvolveTime = time.Now()
		}
	}
}

const UnitActionMove = "run"
const UnitActionIdle = "idle"
