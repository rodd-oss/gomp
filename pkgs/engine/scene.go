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

	Contoller SceneContorller
	Systems   []*System

	ShouldRender bool

	currentTick uint
	syncPeriod  uint // in ticks
}

type SceneContorller interface {
	Load(scene *Scene) (s []*System, e []Entity)
	Update(dt float64)
	Unload(scene *Scene)
}

func CreateScene(controller SceneContorller) (scene Scene) {
	scene.World = ecs.NewWorld()
	scene.Space = cp.NewSpace()
	scene.Space.Iterations = 1
	scene.Contoller = controller
	scene.currentTick = 0
	scene.syncPeriod = 3
	scene.ShouldRender = false

	return scene
}

func (s *Scene) Update(dt float64) {
	if s.Engine.Debug {
		log.Println("Scene Updating:", s.Name)
		defer log.Println("Scene Updated:", s.Name)
	}

	for i := range s.Systems {
		s.Systems[i].Update(dt)
	}

	s.Contoller.Update(dt)

	// needToSync := s.currentTick%s.syncPeriod == 0

	if s.currentTick%s.syncPeriod == 0 {
		//send s.Patch
	}

	// network sync
	s.currentTick++
}
