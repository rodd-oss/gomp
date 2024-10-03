package components

import (
	"tomb_mates/internal/protos"

	ecs "github.com/yohamta/donburi"
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
}
