package tinyrpg

import (
	"log"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

// World represents game state
type World struct {
	Replica bool
	Units   map[string]*Unit
	MyID    string
}

func (world *World) AddPlayer() string {
	skins := []string{"big_demon", "big_zombie", "elf_f"}
	id := uuid.NewV4().String()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	unit := &Unit{
		Id:     id,
		X:      rnd.Float64()*300 + 10,
		Y:      rnd.Float64()*220 + 10,
		Frame:  int32(rnd.Intn(4)),
		Skin:   skins[rnd.Intn(len(skins))],
		Action: "idle",
		Speed:  1,
	}
	world.Units[id] = unit

	return id
}

func (world *World) HandleEvent(event *Event) {
	log.Println(event.GetType())
	log.Println(event.GetData())

	switch event.GetType() {
	case Event_type_connect:
		data := event.GetConnect()
		world.Units[data.Unit.Id] = data.Unit

	case Event_type_init:
		data := event.GetInit()
		if world.Replica {
			world.MyID = data.PlayerId
			world.Units = data.Units
		}

	case Event_type_exit:
		data := event.GetExit()
		delete(world.Units, data.PlayerId)

	case Event_type_move:
		data := event.GetMove()
		unit := world.Units[data.PlayerId]
		unit.Action = UnitActionMove
		unit.Direction = data.Direction

	case Event_type_idle:
		data := event.GetIdle()
		unit := world.Units[data.PlayerId]
		unit.Action = UnitActionIdle

	default:
		log.Println("UNKNOWN EVENT: ", event)
	}
}

func (world *World) Evolve() {
	ticker := time.NewTicker(time.Second / 60)

	for {
		select {
		case <-ticker.C:
			for _, unit := range world.Units {
				if unit.Action == UnitActionMove {
					switch unit.Direction {
					case Direction_left:
						unit.X -= unit.Speed
						unit.Side = Direction_left
					case Direction_right:
						unit.X += unit.Speed
						unit.Side = Direction_right
					case Direction_up:
						unit.Y -= unit.Speed
					case Direction_down:
						unit.Y += unit.Speed
					default:
						log.Println("UNKNOWN DIRECTION: ", unit.Direction)
					}
				}
			}
		}
	}
}

const UnitActionMove = "run"
const UnitActionIdle = "idle"
