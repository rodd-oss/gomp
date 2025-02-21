/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package ecs

import (
	"reflect"
)

type World[C, S any] struct {
	Entities   EntityManager
	Components C
	Systems    S
}

func NewWorld[C AnyComponentList, S AnySystemList](componentList C, systemList S) World[C, S] {
	return World[C, S]{
		Entities:   NewEntityManager(),
		Components: componentList,
		Systems:    systemList,
	}
}

func (w *World[C, S]) Init() {
	w.injectComponentsToSystems()
	w.injectEntityManagerToComponents()
	w.Entities.init()
}

func (w *World[C, S]) Destroy() {
	w.Entities.Destroy()
	//w.Components.Destroy()
	//w.Systems.Destroy()
}

func (w *World[C, S]) injectEntityManagerToComponents() {
	componentList := &w.Components
	entityManager := &w.Entities

	reflectedComponentList := reflect.ValueOf(componentList).Elem()
	componentListLen := reflectedComponentList.NumField()

	for k := range componentListLen {
		component := reflectedComponentList.Field(k)
		componentManager, ok := component.Addr().Interface().(AnyComponentManagerPtr)
		if !ok {
			continue
		}
		entityManager.registerComponent(componentManager)
		componentManager.registerEntityManager(entityManager)
	}
}

// injectToSystems
func (w *World[C, S]) injectComponentsToSystems() {
	systemList := &w.Systems
	componentList := &w.Components
	entityManager := &w.Entities

	reflectedSystemList := reflect.ValueOf(systemList).Elem()
	systemsLen := reflectedSystemList.NumField()

	reflectedComponentList := reflect.ValueOf(componentList).Elem()
	componentsLen := reflectedComponentList.NumField()

	entityManagerType := reflect.TypeOf(entityManager)

	for i := range systemsLen {
		system := reflectedSystemList.Field(i)
		systemLen := system.NumField()

		for j := range systemLen {
			systemField := system.Field(j)
			systemFieldType := systemField.Type()

			if systemFieldType.Kind() != reflect.Ptr {
				continue
			}

			if systemFieldType == entityManagerType {
				system.Field(j).Set(reflect.ValueOf(entityManager))
				continue
			}

			// TODO: refactor to component list indexed map to speed up assignment
			for k := range componentsLen {
				component := reflectedComponentList.Field(k)
				componentType := component.Type()

				if systemFieldType.Elem() == componentType {
					system.Field(j).Set(component.Addr())
					break
				}
			}
		}
	}
}
