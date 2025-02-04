package ecs

import (
	"iter"
	"maps"
	"reflect"

	"github.com/negrel/assert"
)

// Selector can be used in component controller to track entities which have
// provided components set.
type AnySelector interface {
	IncludeMask() ComponentBitArray256
	ExcludeMask() ComponentBitArray256
	NumEntities() int
	AllEntities() iter.Seq[Entity]
}

type selectorBackdoor interface {
	AnySelector
	initInWorld(world *World)
	addEntity(entId Entity)
	removeEntity(entId Entity)
}

type selectorBase struct {
	includeMask ComponentBitArray256
	excludeMask ComponentBitArray256
	matchedEnts map[Entity]struct{}
}

func (s *selectorBase) IncludeMask() ComponentBitArray256 {
	return s.includeMask
}

func (s *selectorBase) ExcludeMask() ComponentBitArray256 {
	return s.excludeMask
}

func (s *selectorBase) NumEntities() int {
	return len(s.matchedEnts)
}

func (s *selectorBase) AllEntities() iter.Seq[Entity] {
	return maps.Keys(s.matchedEnts) // FIXME(?): map iterates in random order every time
}

func (s *selectorBase) addEntity(entId Entity) {
	if s.matchedEnts == nil {
		s.matchedEnts = map[Entity]struct{}{}
	}
	s.matchedEnts[entId] = struct{}{}
}

func (s *selectorBase) removeEntity(entId Entity) {
	delete(s.matchedEnts, entId)
}

func (s *selectorBase) makeMasks(includeComponents ...AnyComponentManagerPtr) {
	s.includeMask = ComponentBitArray256{}
	s.excludeMask = ComponentBitArray256{}

	for _, mng := range includeComponents {
		s.includeMask.Set(mng.getId())
	}
}

type Selector[T any] struct {
	selectorBase
	meta []selectorMeta
}

type selectorMeta struct {
	fld reflect.StructField
	mng AnyComponentManagerPtr
}

func (s *Selector[T]) initInWorld(world *World) {
	tTyp := reflect.TypeFor[T]()
	managers := []AnyComponentManagerPtr{}

	for fldIdx := range tTyp.NumField() {
		fld := tTyp.Field(fldIdx)
		assert.Equal(reflect.Pointer, fld.Type.Kind(), "field in GSelector type argument must be pointer to component type")
		assert.True(fld.IsExported(), "field in GSelector must be exported")

		compTyp := fld.Type.Elem()
		found := false

		for _, mng := range world.components {
			if compTyp == mng.getComponentType() {
				managers = append(managers, mng)
				s.meta = append(s.meta, selectorMeta{
					fld: fld,
					mng: mng,
				})
				found = true
				break
			}
		}

		assert.True(found, "component type %s not found in world's component managers", fld.Type.String())
	}

	s.makeMasks(managers...)
}

func (s *Selector[T]) pullComponentInstances(entId Entity, dst *T) {
	dstVal := reflect.ValueOf(dst).Elem()
	for _, it := range s.meta {
		comp := it.mng.GetComponent(entId)
		dstVal.FieldByIndex(it.fld.Index).Set(reflect.ValueOf(comp))
	}
}

func (s *Selector[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for entId := range s.matchedEnts {
			var bundle T
			s.pullComponentInstances(entId, &bundle)
			if !yield(bundle) {
				break
			}
		}
	}
}
