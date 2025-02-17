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

package scenes

import (
	"gomp"
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/systems"
	"gomp/pkg/ecs"
	"time"
)

func NewMainScene(world *ecs.World, components *components.GameComponents) *MainScene {
	return &MainScene{
		playerSystem: systems.NewPlayerSystem(world, components.SpriteMatrix, components.Position, components.Rotation, components.Scale, components.Velocity, components.AnimationPlayer, components.AnimationState, components.Tint, components.Flip),
	}
}

type MainScene struct {
	playerSystem *systems.PlayerSystem
}

func (s *MainScene) Init() {
	s.playerSystem.Init()
}

func (s *MainScene) Destroy() {
	s.playerSystem.Destroy()
}

func (s *MainScene) Update(dt time.Duration) gomp.SceneId {
	s.playerSystem.Run(dt)
	return MainSceneId
}

func (s *MainScene) FixedUpdate(dt time.Duration) {

}

func (s *MainScene) OnEnter() {

}

func (s *MainScene) OnExit() {

}

var _ gomp.AnyScene = (*MainScene)(nil)
