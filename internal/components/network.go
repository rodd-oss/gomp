package components

import (
	"tomb_mates/internal/protos"

	ecs "github.com/yohamta/donburi"
	"google.golang.org/protobuf/proto"
)

type NetworkUnitData struct {
	Unit *protos.Unit
}

var NetworkUnit = ecs.NewComponentType[NetworkUnitData]()

type NetworkAreaData struct {
	Area *protos.Area
}

var NetworkArea = ecs.NewComponentType[NetworkAreaData]()

type NetworkManagerData struct {
	MyID *uint32

	UnhandledEvents []*protos.Event

	AreaToEntity    map[uint32]ecs.Entity
	Areas           map[uint32]*protos.Area
	PatchedAreas    map[uint32]*protos.PatchArea
	DeletedAreasIds map[uint32]*protos.Empty
	CreatedAreas    map[uint32]*protos.Area

	UnitToEntity    map[uint32]ecs.Entity
	Units           map[uint32]*protos.Unit
	PatchedUnits    map[uint32]*protos.PatchUnit
	DeletedUnitsIds map[uint32]*protos.Empty
	CreatedUnits    map[uint32]*protos.Unit

	Broadcast chan []byte
}

func (n *NetworkManagerData) SendPatch() {
	if len(n.PatchedUnits) == 0 && len(n.CreatedUnits) == 0 && len(n.DeletedUnitsIds) == 0 {
		return
	}

	statePatchEvent := &protos.Event{
		Type: protos.EventType_state_patch,
		Data: &protos.Event_StatePatch{
			StatePatch: &protos.GameStatePatche{
				Units:           n.PatchedUnits,
				CreatedUnits:    n.CreatedUnits,
				DeletedUnitsIds: n.DeletedUnitsIds,
			},
		},
	}

	message, err := proto.Marshal(statePatchEvent)
	if err != nil {
		return
	}
	n.Broadcast <- message
	n.PatchedUnits = make(map[uint32]*protos.PatchUnit)
	n.CreatedUnits = make(map[uint32]*protos.Unit)
	n.DeletedUnitsIds = make(map[uint32]*protos.Empty)
}
