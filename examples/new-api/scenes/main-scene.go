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
	"gomp/examples/new-api/instances"
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"time"
)

func NewMainScene() MainScene {
	return MainScene{}
}

type MainScene struct {
	Game  core.AnyGame
	World *instances.World
}

func (s *MainScene) Id() gomp.SceneId {
	return MainSceneId
}

func (s *MainScene) Init(world ecs.AnyWorld) {
	s.World = world.(*instances.World)

	// Network receive
	s.World.Systems.Network.Init()
	s.World.Systems.NetworkReceive.Init()

	// Scenes
	s.World.Systems.Player.Init()
	s.World.Systems.Velocity.Init()

	s.World.Systems.CollisionDetectionGrid.Init()

	// Network patches
	s.World.Systems.NetworkSend.Init()

	// Animation
	s.World.Systems.AnimationSpriteMatrix.Init()
	s.World.Systems.AnimationPlayer.Init()

	s.World.Systems.SpriteMatrix.Init()
	s.World.Systems.YSort.Init()

	// RenderAssterodd
	s.World.Systems.RenderBogdan.Init()
	s.World.Systems.Debug.Init()
	s.World.Systems.AssetLib.Init()
}

func (s *MainScene) Update(dt time.Duration) gomp.SceneId {
	// Network receive
	s.World.Systems.NetworkReceive.Run(dt)
	s.World.Systems.Player.Run()

	return MainSceneId
}

func (s *MainScene) FixedUpdate(dt time.Duration) {
	// Network send
	s.World.Systems.Network.Run(dt)

	s.World.Systems.Velocity.Run(dt)
	s.World.Systems.CollisionDetectionGrid.Run(dt)
	s.World.Systems.CollisionHandler.Run(dt)
	s.World.Systems.NetworkSend.Run(dt)
}

func (s *MainScene) Render(dt time.Duration) {
	// Animation
	s.World.Systems.AnimationSpriteMatrix.Run()
	s.World.Systems.AnimationPlayer.Run()

	s.World.Systems.SpriteMatrix.Run()
	s.World.Systems.Debug.Run()
	s.World.Systems.AssetLib.Run()
	s.World.Systems.YSort.Run()
}

func (s *MainScene) Destroy() {
	s.World.Destroy()
	// Network intents
	s.World.Systems.Network.Destroy()
	s.World.Systems.NetworkReceive.Destroy()

	s.World.Systems.Player.Destroy()

	s.World.Systems.CollisionDetectionGrid.Destroy()

	// Network patches
	s.World.Systems.NetworkSend.Destroy()

	// Animation
	s.World.Systems.AnimationSpriteMatrix.Destroy()
	s.World.Systems.AnimationPlayer.Destroy()

	s.World.Systems.SpriteMatrix.Destroy()
	s.World.Systems.YSort.Destroy()

	// RenderAssterodd
	s.World.Systems.Debug.Destroy()
	s.World.Systems.AssetLib.Destroy()
	s.World.Systems.RenderBogdan.Destroy()
}

func (s *MainScene) OnEnter() {

}

func (s *MainScene) OnExit() {

}

var _ gomp.AnyScene = (*MainScene)(nil)
