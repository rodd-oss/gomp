package ecs

import (
	"iter"
	"maps"
)

// Selector can be used in component controller to track entities which have
// provided components set.
type Selector struct {
	includeMask ComponentBitArray256
	excludeMask ComponentBitArray256
	matchedEnts map[Entity]struct{}
}

func newSelector(includeAllMask, excludeAnyMask ComponentBitArray256) *Selector {
	return &Selector{
		includeMask: includeAllMask,
		excludeMask: excludeAnyMask,
		matchedEnts: map[Entity]struct{}{},
	}
}

func (s *Selector) NumEntities() int {
	return len(s.matchedEnts)
}

func (s *Selector) EnumEntities() iter.Seq[Entity] {
	return maps.Keys(s.matchedEnts)
}

func (s *Selector) addEntity(entId Entity) {
	s.matchedEnts[entId] = struct{}{}
}

func (s *Selector) removeEntity(entId Entity) {
	delete(s.matchedEnts, entId)
}
