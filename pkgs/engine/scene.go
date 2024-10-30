/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import (
	"log"

	capnp "capnproto.org/go/capnp/v3"
	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
)

type Scene struct {
	Name   string
	Engine *Engine

	World ecs.World
	Space *cp.Space

	Entities   []ecs.Entity
	Components []Component[capnp.Struct]

	Contoller SceneContorller

	currentTick uint
	syncPeriod  uint // in ticks
}

type SceneContorller interface {
	Update(dt float64)
	OnLoad(scene *Scene)
	OnUnload(scene *Scene)
}

func NewScene(controller SceneContorller) Scene {
	scene := Scene{}

	scene.World = ecs.NewWorld()
	scene.Space = cp.NewSpace()
	scene.Contoller = controller
	scene.currentTick = 0
	scene.syncPeriod = 3

	return scene
}

func (s *Scene) Update(dt float64) {
	log.Println("Scene update:", s.Name)
	s.Contoller.Update(dt)
	needToSync := s.currentTick%s.syncPeriod == 0

	for i := range s.Components {
		s.Components[i].System.Each(s.World, func(e *ecs.Entry) {
			comp := s.Components[i].System.GetValue(e)
			comp.Controller.Update(dt)

			if needToSync {
				comp.Controller.OnStateRequest(comp.State).Message().Marshal()
			}
		})
	}

	s.Space.Step(dt)

	if s.currentTick%s.syncPeriod == 0 {
		//send s.Patch
	}

	// network sync
	s.currentTick++
}

// type EntityState struct {
// 	Id         uint
// 	Components []ComponentState
// }

// type SceneState map[uint]EntityState
