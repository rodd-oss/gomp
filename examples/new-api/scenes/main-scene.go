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
	"gomp/pkg/ecs"
	"time"
)

func NewMainScene() *MainScene {
	world := ecs.CreateWorld("Main")
	comps := instances.NewComponentList(&world)
	sys := instances.NewSystemList(&world, &comps)

	scene := MainScene{
		World:      world,
		Components: comps,
		Systems:    sys,
	}

	return &scene
}

type MainScene struct {
	Game       *gomp.Game
	World      ecs.World
	Components instances.ComponentList
	Systems    instances.SystemList
}

func (s *MainScene) Id() gomp.SceneId {
	return MainSceneId
}

func (s *MainScene) Init() {
	// Network receive
	s.Systems.Network.Init()
	s.Systems.NetworkReceive.Init()

	// Scenes
	s.Systems.Player.Init()

	s.Systems.Velocity.Init()

	// Network patches
	s.Systems.NetworkSend.Init()

	// Animation
	s.Systems.AnimationSpriteMatrix.Init()
	s.Systems.AnimationPlayer.Init()

	// Prerender init
	s.Systems.TextureRenderSprite.Init()
	s.Systems.TextureRenderSpriteSheet.Init()
	s.Systems.TextureRenderMatrix.Init()

	// Prerender fill
	s.Systems.TextureRenderAnimation.Init()
	s.Systems.TextureRenderFlip.Init()
	s.Systems.TextureRenderPosition.Init()
	s.Systems.TextureRenderRotation.Init()
	s.Systems.TextureRenderScale.Init()
	s.Systems.TextureRenderTint.Init()

	// Render
	s.Systems.Render.Init()
	s.Systems.Debug.Init()
	s.Systems.AssetLib.Init()
}

func (s *MainScene) Update(dt time.Duration) gomp.SceneId {

	// Network receive
	s.Systems.Network.Run(dt)
	s.Systems.NetworkReceive.Run(dt)

	s.Systems.Player.Run(dt)

	// Network patches
	s.Systems.NetworkSend.Run(dt)

	s.Systems.Debug.Run(dt)
	return MainSceneId
}

func (s *MainScene) FixedUpdate(dt time.Duration) {
	// Network send
	s.Systems.NetworkSend.Run(dt)
}

func (s *MainScene) Render(dt time.Duration) {
	s.Systems.Velocity.Run(dt)

	// Animation
	s.Systems.AnimationSpriteMatrix.Run(dt)
	s.Systems.AnimationPlayer.Run(dt)

	// Prerender init
	s.Systems.TextureRenderSprite.Run(dt)
	s.Systems.TextureRenderSpriteSheet.Run(dt)
	s.Systems.TextureRenderMatrix.Run(dt)

	// Prerender fill
	s.Systems.TextureRenderAnimation.Run(dt)
	s.Systems.TextureRenderFlip.Run(dt)
	s.Systems.TextureRenderPosition.Run(dt)
	s.Systems.TextureRenderRotation.Run(dt)
	s.Systems.TextureRenderScale.Run(dt)
	s.Systems.TextureRenderTint.Run(dt)

	// Render
	s.Systems.AssetLib.Run(dt)
	shouldContinue := s.Systems.Render.Run(dt)
	if !shouldContinue {
		s.Game.SetShouldDestroy(true)
		return
	}
}

func (s *MainScene) Destroy() {
	// Network intents
	s.Systems.Network.Destroy()
	s.Systems.NetworkReceive.Destroy()

	s.Systems.Player.Destroy()

	// Network patches
	s.Systems.NetworkSend.Destroy()

	// Animation
	s.Systems.AnimationSpriteMatrix.Destroy()
	s.Systems.AnimationPlayer.Destroy()

	// Prerender init
	s.Systems.TextureRenderSprite.Destroy()
	s.Systems.TextureRenderSpriteSheet.Destroy()
	s.Systems.TextureRenderMatrix.Destroy()

	// Prerender fill
	s.Systems.TextureRenderAnimation.Destroy()
	s.Systems.TextureRenderFlip.Destroy()
	s.Systems.TextureRenderPosition.Destroy()
	s.Systems.TextureRenderRotation.Destroy()
	s.Systems.TextureRenderScale.Destroy()
	s.Systems.TextureRenderTint.Destroy()

	// Render
	s.Systems.Debug.Destroy()
	s.Systems.AssetLib.Destroy()
	s.Systems.Render.Destroy()

}

func (s *MainScene) OnEnter() {

}

func (s *MainScene) OnExit() {

}

var _ gomp.AnyScene = (*MainScene)(nil)
