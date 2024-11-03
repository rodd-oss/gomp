/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import (
	"log"

	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
)

type Scene struct {
	Name   string
	Engine *Engine

	World ecs.World
	Space *cp.Space

	Systems  []System
	entities []Entity

	ShouldRender bool

	currentTick uint
	syncPeriod  uint // in ticks
}

type sceneFactoryEntities struct {
	scene *Scene
}

func (f sceneFactoryEntities) AddEntities(ent ...Entity) sceneFactorySystems {
	f.scene.entities = ent
	return sceneFactorySystems(f)
}

type sceneFactorySystems struct {
	scene *Scene
}

func (f sceneFactorySystems) AddSystems(sys ...System) Scene {
	f.scene.Systems = sys
	return *f.scene
}

func CreateScene(name string) sceneFactoryEntities {
	scene := Scene{}

	scene.Name = name
	scene.currentTick = 0
	scene.syncPeriod = 3
	scene.ShouldRender = false

	factory := sceneFactoryEntities{
		scene: &scene,
	}

	return factory
}

func (s *Scene) Load() {
	s.World = ecs.NewWorld()
	s.Space = cp.NewSpace()
	s.Space.Iterations = 1

	for i := range s.entities {
		s.World.Create(s.entities[i]...)
	}

	for i := range s.Systems {
		s.Systems[i].Init(s)
	}
}

func (s *Scene) Unload() {}

func (s *Scene) Update(dt float64) {
	if s.Engine.Debug {
		log.Println("Scene Updating:", s.Name)
		defer log.Println("Scene Updated:", s.Name)
	}

	for i := range s.Systems {
		s.Systems[i].Update(dt)
	}

	// needToSync := s.currentTick%s.syncPeriod == 0

	if s.currentTick%s.syncPeriod == 0 {
		//send s.Patch
	}

	// network sync
	s.currentTick++
}
