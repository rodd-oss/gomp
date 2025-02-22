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

func NewMainScene() MainScene {
	return MainScene{
		World: ecs.NewWorld(instances.NewComponentList(), instances.NewSystemList()),
	}
}

type MainScene struct {
	Game  *gomp.Game
	World instances.World
}

func (s *MainScene) Id() gomp.SceneId {
	return MainSceneId
}

func (s *MainScene) Init() {
	s.World.Init()

	// Network receive
	s.World.Systems.Network.Init()
	s.World.Systems.NetworkReceive.Init()

	// Scenes
	s.World.Systems.Player.Init()

	s.World.Systems.Velocity.Init()

	// Network patches
	s.World.Systems.NetworkSend.Init()

	// Animation
	s.World.Systems.AnimationSpriteMatrix.Init()
	s.World.Systems.AnimationPlayer.Init()

	// Prerender init
	s.World.Systems.TextureRenderSprite.Init()
	s.World.Systems.TextureRenderSpriteSheet.Init()
	s.World.Systems.TextureRenderMatrix.Init()

	// Prerender fill
	s.World.Systems.TextureRenderAnimation.Init()
	s.World.Systems.TextureRenderFlip.Init()
	s.World.Systems.TextureRenderPosition.Init()
	s.World.Systems.TextureRenderRotation.Init()
	s.World.Systems.TextureRenderScale.Init()
	s.World.Systems.TextureRenderTint.Init()

	// Render
	s.World.Systems.Render.Init()
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
	s.World.Systems.ViewPosition.Run()
	s.World.Systems.NetworkSend.Run(dt)
}

func (s *MainScene) Render(interpolation float32) {
	// Animation
	s.World.Systems.AnimationSpriteMatrix.Run()
	s.World.Systems.AnimationPlayer.Run()

	// Prerender init
	s.World.Systems.TextureRenderSprite.Run()
	s.World.Systems.TextureRenderSpriteSheet.Run()
	s.World.Systems.TextureRenderMatrix.Run()

	// Prerender fill
	s.World.Systems.TextureRenderAnimation.Run()
	s.World.Systems.TextureRenderFlip.Run()
	s.World.Systems.TextureRenderPosition.Run(interpolation)
	s.World.Systems.TextureRenderRotation.Run(interpolation)
	s.World.Systems.TextureRenderScale.Run(interpolation)
	s.World.Systems.TextureRenderTint.Run(interpolation)

	// Render
	s.World.Systems.Debug.Run()

	s.World.Systems.AssetLib.Run()
	shouldContinue := s.World.Systems.Render.Run()
	if !shouldContinue {
		s.Game.SetShouldDestroy(true)
		return
	}
}

func (s *MainScene) Destroy() {
	s.World.Destroy()
	// Network intents
	s.World.Systems.Network.Destroy()
	s.World.Systems.NetworkReceive.Destroy()

	s.World.Systems.Player.Destroy()

	// Network patches
	s.World.Systems.NetworkSend.Destroy()

	// Animation
	s.World.Systems.AnimationSpriteMatrix.Destroy()
	s.World.Systems.AnimationPlayer.Destroy()

	// Prerender init
	s.World.Systems.TextureRenderSprite.Destroy()
	s.World.Systems.TextureRenderSpriteSheet.Destroy()
	s.World.Systems.TextureRenderMatrix.Destroy()

	// Prerender fill
	s.World.Systems.TextureRenderAnimation.Destroy()
	s.World.Systems.TextureRenderFlip.Destroy()
	s.World.Systems.TextureRenderPosition.Destroy()
	s.World.Systems.TextureRenderRotation.Destroy()
	s.World.Systems.TextureRenderScale.Destroy()
	s.World.Systems.TextureRenderTint.Destroy()

	// Render
	s.World.Systems.Debug.Destroy()
	s.World.Systems.AssetLib.Destroy()
	s.World.Systems.Render.Destroy()
}

func (s *MainScene) OnEnter() {

}

func (s *MainScene) OnExit() {

}

var _ gomp.AnyScene = (*MainScene)(nil)
