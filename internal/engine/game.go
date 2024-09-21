package engine

import (
	"log"
	"math/rand"
	"sync"
	"time"
	"tomb_mates/internal/protos"

	uuid "github.com/satori/go.uuid"
)

// World represents game state
type World struct {
	Mx      sync.Mutex
	Replica bool
	Units   map[string]*protos.Unit
	MyID    string
}

func (world *World) AddPlayer() string {
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
		Speed:  2,
	}
	world.Units[id] = unit

	return id
}

func (world *World) HandleEvent(event *protos.Event) {
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

func (world *World) Evolve() {
	ticker := time.NewTicker(time.Second / 30)

	for {
		select {
		case <-ticker.C:
			world.Mx.Lock()
			for i := range world.Units {
				if world.Units[i].Action == UnitActionMove {
					switch world.Units[i].Direction {
					case protos.Direction_left:
						world.Units[i].X -= world.Units[i].Speed
						world.Units[i].Side = protos.Direction_left
					case protos.Direction_right:
						world.Units[i].X += world.Units[i].Speed
						world.Units[i].Side = protos.Direction_right
					case protos.Direction_up:
						world.Units[i].Y -= world.Units[i].Speed
					case protos.Direction_down:
						world.Units[i].Y += world.Units[i].Speed
					default:
						log.Println("UNKNOWN DIRECTION: ", world.Units[i].Direction)
					}
				}
			}
			world.Mx.Unlock()
		}
	}
}

const UnitActionMove = "run"
const UnitActionIdle = "idle"
