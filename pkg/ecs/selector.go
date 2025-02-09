package ecs

import (
	"iter"
	"reflect"

	"github.com/negrel/assert"
)

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
	matchedEnts *PagedMap[Entity, struct{}]
}

func (s *selectorBase) IncludeMask() ComponentBitArray256 {
	return s.includeMask
}

func (s *selectorBase) ExcludeMask() ComponentBitArray256 {
	return s.excludeMask
}

func (s *selectorBase) NumEntities() int {
	return int(s.matchedEnts.Len())
}

func (s *selectorBase) AllEntities() iter.Seq[Entity] {
	return s.matchedEnts.Keys()
}

func (s *selectorBase) addEntity(entId Entity) {
	if s.matchedEnts == nil {
		s.matchedEnts = NewPagedMap[Entity, struct{}]()
	}
	s.matchedEnts.Set(entId, struct{}{})
}

func (s *selectorBase) removeEntity(entId Entity) {
	s.matchedEnts.Delete(entId)
}

func (s *selectorBase) makeMasks(includeComponents []AnyComponentManagerPtr, excludeComponents []AnyComponentManagerPtr) {
	s.includeMask = ComponentBitArray256{}
	s.excludeMask = ComponentBitArray256{}

	for _, mng := range includeComponents {
		s.includeMask.Set(mng.getId())
	}

	for _, mng := range excludeComponents {
		s.excludeMask.Set(mng.getId())
	}
}

// Selector can be used in component controller to track entities which have
// provided components set. Don't forget to register selector in world using
// `world.RegisterSelector(&selector)` method
//
// Note that fields in selector's struct must be of reference type.
//
//	rbSelector ecs.Selector[struct{
//	  RigidBody *components.RigidBody
//	  Mass      *components.Mass
//	  Position  *components.Position
//	  Rotation  *components.Rotation
//	}]
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
	includeManagers := []AnyComponentManagerPtr{}
	excludeManagers := []AnyComponentManagerPtr{}

	for fldIdx := range tTyp.NumField() {
		fld := tTyp.Field(fldIdx)

		if fld.Type.Kind() == reflect.Struct && fld.Type.Implements(reflect.TypeFor[exclude]()) {
			tFld, ok := fld.Type.FieldByName("t")
			assert.True(ok, "type Exclude[T] doesn't contains field 't', someone removed/renamed it?")

			excludedType := tFld.Type.Elem()
			found := false
			for _, mng := range world.components {
				if excludedType == mng.getComponentType() {
					excludeManagers = append(excludeManagers, mng)
					found = true
					break
				}
			}

			assert.True(found, "component type %s not found in world's component managers", fld.Type.String())

			continue
		}

		assert.Equal(reflect.Pointer, fld.Type.Kind(), "field in Selector type argument must be pointer to component type")
		assert.True(fld.IsExported(), "field in Selector must be exported")

		compTyp := fld.Type.Elem()
		found := false

		for _, mng := range world.components {
			if compTyp == mng.getComponentType() {
				includeManagers = append(includeManagers, mng)
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

	s.makeMasks(includeManagers, excludeManagers)
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
		for entId := range s.matchedEnts.Keys() {
			var bundle T
			s.pullComponentInstances(entId, &bundle)
			if !yield(bundle) {
				break
			}
		}
	}
}

// Add field with this type to selector's struct to exclude entities
// from selecting them.
//
// Note that field with Exclude[T] must be of value type, not reference.
//
// Example:
//
//	camSelector ecs.Selector[struct{
//	  Camera *components.CameraOrtho
//	  ecs.Exclude[components.Disabled]
//	}]
type Exclude[T any] struct {
	t [0]T // for reflect
}

// for reflect
func (Exclude[T]) isExclude() {}

// for reflect
type exclude interface {
	isExclude()
}
