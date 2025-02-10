package ecs

import (
	"fmt"
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
	if s.matchedEnts == nil {
		return 0
	}
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

func pullComponentManager[T any](world *World, dst **ComponentManager[T]) error {
	tTyp := reflect.TypeFor[T]()

	for _, mng := range world.components {
		if mng.getComponentType() == tTyp {
			*dst = mng.(*ComponentManager[T])
			return nil
		}
	}

	return fmt.Errorf("world doesn't contains component %s", tTyp.String())
}

type Selector2[T1, T2 any] struct {
	selectorBase
	t1Manager *ComponentManager[T1]
	t2Manager *ComponentManager[T2]
}

func (s *Selector2[T1, T2]) initInWorld(world *World) {
	pullComponentManager(world, &s.t1Manager)
	pullComponentManager(world, &s.t2Manager)
	s.makeMasks([]AnyComponentManagerPtr{
		s.t1Manager,
		s.t2Manager,
	},
		[]AnyComponentManagerPtr{},
	)
}

type Selector2Elements[T1 any, T2 any] struct {
	C1 *T1
	C2 *T2
}

func (s *Selector2[T1, T2]) All() iter.Seq[Selector2Elements[T1, T2]] {
	return func(yield func(Selector2Elements[T1, T2]) bool) {
		for entId := range s.matchedEnts.Keys() {
			var elems Selector2Elements[T1, T2]
			elems.C1 = s.t1Manager.Get(entId)
			elems.C2 = s.t2Manager.Get(entId)
			if !yield(elems) {
				break
			}
		}
	}
}

type Selector3[T1, T2, T3 any] struct {
	selectorBase
	t1Manager *ComponentManager[T1]
	t2Manager *ComponentManager[T2]
	t3Manager *ComponentManager[T3]
}

func (s *Selector3[T1, T2, T3]) initInWorld(world *World) {
	pullComponentManager(world, &s.t1Manager)
	pullComponentManager(world, &s.t2Manager)
	pullComponentManager(world, &s.t3Manager)
	s.makeMasks([]AnyComponentManagerPtr{
		s.t1Manager,
		s.t2Manager,
		s.t3Manager,
	},
		[]AnyComponentManagerPtr{},
	)
}

type Selector3Elements[T1, T2, T3 any] struct {
	C1 *T1
	C2 *T2
	C3 *T3
}

func (s *Selector3[T1, T2, T3]) All() iter.Seq[Selector3Elements[T1, T2, T3]] {
	return func(yield func(Selector3Elements[T1, T2, T3]) bool) {
		for entId := range s.matchedEnts.Keys() {
			var elems Selector3Elements[T1, T2, T3]
			elems.C1 = s.t1Manager.Get(entId)
			elems.C2 = s.t2Manager.Get(entId)
			elems.C3 = s.t3Manager.Get(entId)
			if !yield(elems) {
				break
			}
		}
	}
}

type Selector4[T1, T2, T3, T4 any] struct {
	selectorBase
	t1Manager *ComponentManager[T1]
	t2Manager *ComponentManager[T2]
	t3Manager *ComponentManager[T3]
	t4Manager *ComponentManager[T4]
}

func (s *Selector4[T1, T2, T3, T4]) initInWorld(world *World) {
	pullComponentManager(world, &s.t1Manager)
	pullComponentManager(world, &s.t2Manager)
	pullComponentManager(world, &s.t3Manager)
	pullComponentManager(world, &s.t4Manager)
	s.makeMasks([]AnyComponentManagerPtr{
		s.t1Manager,
		s.t2Manager,
		s.t3Manager,
		s.t4Manager,
	},
		[]AnyComponentManagerPtr{},
	)
}

type Selector4Elements[T1, T2, T3, T4 any] struct {
	C1 *T1
	C2 *T2
	C3 *T3
	C4 *T4
}

func (s *Selector4[T1, T2, T3, T4]) All() iter.Seq[Selector4Elements[T1, T2, T3, T4]] {
	return func(yield func(Selector4Elements[T1, T2, T3, T4]) bool) {
		for entId := range s.matchedEnts.Keys() {
			var elems Selector4Elements[T1, T2, T3, T4]
			elems.C1 = s.t1Manager.Get(entId)
			elems.C2 = s.t2Manager.Get(entId)
			elems.C3 = s.t3Manager.Get(entId)
			elems.C4 = s.t4Manager.Get(entId)
			if !yield(elems) {
				break
			}
		}
	}
}

type Selector5[T1, T2, T3, T4, T5 any] struct {
	selectorBase
	t1Manager *ComponentManager[T1]
	t2Manager *ComponentManager[T2]
	t3Manager *ComponentManager[T3]
	t4Manager *ComponentManager[T4]
	t5Manager *ComponentManager[T5]
}

func (s *Selector5[T1, T2, T3, T4, T5]) initInWorld(world *World) {
	pullComponentManager(world, &s.t1Manager)
	pullComponentManager(world, &s.t2Manager)
	pullComponentManager(world, &s.t3Manager)
	pullComponentManager(world, &s.t4Manager)
	pullComponentManager(world, &s.t5Manager)
	s.makeMasks([]AnyComponentManagerPtr{
		s.t1Manager,
		s.t2Manager,
		s.t3Manager,
		s.t4Manager,
		s.t5Manager,
	},
		[]AnyComponentManagerPtr{},
	)
}

type Selector5Elements[T1, T2, T3, T4, T5 any] struct {
	C1 *T1
	C2 *T2
	C3 *T3
	C4 *T4
	C5 *T5
}

func (s *Selector5[T1, T2, T3, T4, T5]) All() iter.Seq[Selector5Elements[T1, T2, T3, T4, T5]] {
	return func(yield func(Selector5Elements[T1, T2, T3, T4, T5]) bool) {
		for entId := range s.matchedEnts.Keys() {
			var elems Selector5Elements[T1, T2, T3, T4, T5]
			elems.C1 = s.t1Manager.Get(entId)
			elems.C2 = s.t2Manager.Get(entId)
			elems.C3 = s.t3Manager.Get(entId)
			elems.C4 = s.t4Manager.Get(entId)
			elems.C5 = s.t5Manager.Get(entId)
			if !yield(elems) {
				break
			}
		}
	}
}

type Selector6[T1, T2, T3, T4, T5, T6 any] struct {
	selectorBase
	t1Manager *ComponentManager[T1]
	t2Manager *ComponentManager[T2]
	t3Manager *ComponentManager[T3]
	t4Manager *ComponentManager[T4]
	t5Manager *ComponentManager[T5]
	t6Manager *ComponentManager[T6]
}

func (s *Selector6[T1, T2, T3, T4, T5, T6]) initInWorld(world *World) {
	pullComponentManager(world, &s.t1Manager)
	pullComponentManager(world, &s.t2Manager)
	pullComponentManager(world, &s.t3Manager)
	pullComponentManager(world, &s.t4Manager)
	pullComponentManager(world, &s.t5Manager)
	pullComponentManager(world, &s.t6Manager)
	s.makeMasks([]AnyComponentManagerPtr{
		s.t1Manager,
		s.t2Manager,
		s.t3Manager,
		s.t4Manager,
		s.t5Manager,
		s.t6Manager,
	},
		[]AnyComponentManagerPtr{},
	)
}

type Selector6Elements[T1, T2, T3, T4, T5, T6 any] struct {
	C1 *T1
	C2 *T2
	C3 *T3
	C4 *T4
	C5 *T5
	C6 *T6
}

func (s *Selector6[T1, T2, T3, T4, T5, T6]) All() iter.Seq[Selector6Elements[T1, T2, T3, T4, T5, T6]] {
	return func(yield func(Selector6Elements[T1, T2, T3, T4, T5, T6]) bool) {
		for entId := range s.matchedEnts.Keys() {
			var elems Selector6Elements[T1, T2, T3, T4, T5, T6]
			elems.C1 = s.t1Manager.Get(entId)
			elems.C2 = s.t2Manager.Get(entId)
			elems.C3 = s.t3Manager.Get(entId)
			elems.C4 = s.t4Manager.Get(entId)
			elems.C5 = s.t5Manager.Get(entId)
			elems.C6 = s.t6Manager.Get(entId)
			if !yield(elems) {
				break
			}
		}
	}
}

type Selector7[T1, T2, T3, T4, T5, T6, T7 any] struct {
	selectorBase
	t1Manager *ComponentManager[T1]
	t2Manager *ComponentManager[T2]
	t3Manager *ComponentManager[T3]
	t4Manager *ComponentManager[T4]
	t5Manager *ComponentManager[T5]
	t6Manager *ComponentManager[T6]
	t7Manager *ComponentManager[T7]
}

func (s *Selector7[T1, T2, T3, T4, T5, T6, T7]) initInWorld(world *World) {
	pullComponentManager(world, &s.t1Manager)
	pullComponentManager(world, &s.t2Manager)
	pullComponentManager(world, &s.t3Manager)
	pullComponentManager(world, &s.t4Manager)
	pullComponentManager(world, &s.t5Manager)
	pullComponentManager(world, &s.t6Manager)
	pullComponentManager(world, &s.t7Manager)
	s.makeMasks([]AnyComponentManagerPtr{
		s.t1Manager,
		s.t2Manager,
		s.t3Manager,
		s.t4Manager,
		s.t5Manager,
		s.t6Manager,
		s.t7Manager,
	},
		[]AnyComponentManagerPtr{},
	)
}

type Selector7Elements[T1, T2, T3, T4, T5, T6, T7 any] struct {
	C1 *T1
	C2 *T2
	C3 *T3
	C4 *T4
	C5 *T5
	C6 *T6
	C7 *T7
}

func (s *Selector7[T1, T2, T3, T4, T5, T6, T7]) All() iter.Seq[Selector7Elements[T1, T2, T3, T4, T5, T6, T7]] {
	return func(yield func(Selector7Elements[T1, T2, T3, T4, T5, T6, T7]) bool) {
		for entId := range s.matchedEnts.Keys() {
			var elems Selector7Elements[T1, T2, T3, T4, T5, T6, T7]
			elems.C1 = s.t1Manager.Get(entId)
			elems.C2 = s.t2Manager.Get(entId)
			elems.C3 = s.t3Manager.Get(entId)
			elems.C4 = s.t4Manager.Get(entId)
			elems.C5 = s.t5Manager.Get(entId)
			elems.C6 = s.t6Manager.Get(entId)
			elems.C7 = s.t7Manager.Get(entId)
			if !yield(elems) {
				break
			}
		}
	}
}

type Selector8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	selectorBase
	t1Manager *ComponentManager[T1]
	t2Manager *ComponentManager[T2]
	t3Manager *ComponentManager[T3]
	t4Manager *ComponentManager[T4]
	t5Manager *ComponentManager[T5]
	t6Manager *ComponentManager[T6]
	t7Manager *ComponentManager[T7]
	t8Manager *ComponentManager[T8]
}

func (s *Selector8[T1, T2, T3, T4, T5, T6, T7, T8]) initInWorld(world *World) {
	pullComponentManager(world, &s.t1Manager)
	pullComponentManager(world, &s.t2Manager)
	pullComponentManager(world, &s.t3Manager)
	pullComponentManager(world, &s.t4Manager)
	pullComponentManager(world, &s.t5Manager)
	pullComponentManager(world, &s.t6Manager)
	pullComponentManager(world, &s.t7Manager)
	pullComponentManager(world, &s.t8Manager)
	s.makeMasks([]AnyComponentManagerPtr{
		s.t1Manager,
		s.t2Manager,
		s.t3Manager,
		s.t4Manager,
		s.t5Manager,
		s.t6Manager,
		s.t7Manager,
		s.t8Manager,
	},
		[]AnyComponentManagerPtr{},
	)
}

type Selector8Elements[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	C1 *T1
	C2 *T2
	C3 *T3
	C4 *T4
	C5 *T5
	C6 *T6
	C7 *T7
	C8 *T8
}

func (s *Selector8[T1, T2, T3, T4, T5, T6, T7, T8]) All() iter.Seq[Selector8Elements[T1, T2, T3, T4, T5, T6, T7, T8]] {
	return func(yield func(Selector8Elements[T1, T2, T3, T4, T5, T6, T7, T8]) bool) {
		for entId := range s.matchedEnts.Keys() {
			var elems Selector8Elements[T1, T2, T3, T4, T5, T6, T7, T8]
			elems.C1 = s.t1Manager.Get(entId)
			elems.C2 = s.t2Manager.Get(entId)
			elems.C3 = s.t3Manager.Get(entId)
			elems.C4 = s.t4Manager.Get(entId)
			elems.C5 = s.t5Manager.Get(entId)
			elems.C6 = s.t6Manager.Get(entId)
			elems.C7 = s.t7Manager.Get(entId)
			elems.C8 = s.t8Manager.Get(entId)
			if !yield(elems) {
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
