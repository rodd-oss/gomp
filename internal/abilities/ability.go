package abilities

import (
	"time"
	"tomb_mates/internal/effects"
	"tomb_mates/internal/protos"
)

type Position struct {
	X uint32
	Y uint32
}

type Ability struct {
	Caster      *protos.Unit
	Name        string
	Description string
	Type        abilityType
	Cooldown    time.Duration
	Cost        AbilityCost
	Target      AbilityTarget
	Effects     []*effects.Effect
	CastRange   uint32

	PreviousCastAt *time.Time
}

// =============================================================================
type abilityCostType int

const (
	CostTypeNone abilityCostType = iota
	CostTypeMana
	CostTypeHp
)

type AbilityCost struct {
	Type  abilityCostType
	Value uint
}

// =============================================================================
type abilityType int

const (
	TypePassive abilityType = iota
	TypeActive
	TypeToggle
)

// =============================================================================

type abilityTargetType int

const (
	TargetTypeNone abilityTargetType = iota
	TargetTypeUnit
	TargetTypeArea
	TargetTypePoint
	TargetTypeSkillshot
	TargetTypeGlobal
)

type TargetArea struct {
	Position Position
	Range    uint32
	Radius   uint32
}

type AbilityTarget struct {
	Type abilityTargetType
	Unit *protos.Unit
	Area *TargetArea
}
