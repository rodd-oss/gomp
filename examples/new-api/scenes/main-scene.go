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
	scene := new(MainScene)

	scene.EntityManager = ecs.NewEntityManager()
	scene.ComponentList = instances.NewComponentList(&scene.EntityManager)
	scene.SystemList = instances.NewSystemList(&scene.EntityManager, &scene.ComponentList)

	return scene
}

type MainScene struct {
	Game          *gomp.Game
	EntityManager ecs.EntityManager
	ComponentList instances.ComponentList
	SystemList    instances.SystemList
}

func (s *MainScene) Id() gomp.SceneId {
	return MainSceneId
}

func (s *MainScene) Init() {
	// Network receive
	s.SystemList.Network.Init()
	s.SystemList.NetworkReceive.Init()

	// Scenes
	s.SystemList.Player.Init()

	s.SystemList.Velocity.Init()

	// Network patches
	s.SystemList.NetworkSend.Init()

	// Animation
	s.SystemList.AnimationSpriteMatrix.Init()
	s.SystemList.AnimationPlayer.Init()

	// Prerender init
	s.SystemList.TextureRenderSprite.Init()
	s.SystemList.TextureRenderSpriteSheet.Init()
	s.SystemList.TextureRenderMatrix.Init()

	// Prerender fill
	s.SystemList.TextureRenderAnimation.Init()
	s.SystemList.TextureRenderFlip.Init()
	s.SystemList.TextureRenderPosition.Init()
	s.SystemList.TextureRenderRotation.Init()
	s.SystemList.TextureRenderScale.Init()
	s.SystemList.TextureRenderTint.Init()

	// Render
	s.SystemList.Render.Init()
	s.SystemList.Debug.Init()
	s.SystemList.AssetLib.Init()
}

func (s *MainScene) Update(dt time.Duration) gomp.SceneId {

	// Network receive
	s.SystemList.Network.Run(dt)
	s.SystemList.NetworkReceive.Run(dt)

	s.SystemList.Player.Run(dt)

	s.SystemList.Debug.Run(dt)
	return MainSceneId
}

func (s *MainScene) FixedUpdate(dt time.Duration) {
	// Network send
	s.SystemList.NetworkSend.Run(dt)
}

func (s *MainScene) Render(dt time.Duration) {
	s.SystemList.Velocity.Run(dt)

	// Animation
	s.SystemList.AnimationSpriteMatrix.Run(dt)
	s.SystemList.AnimationPlayer.Run(dt)

	// Prerender init
	s.SystemList.TextureRenderSprite.Run(dt)
	s.SystemList.TextureRenderSpriteSheet.Run(dt)
	s.SystemList.TextureRenderMatrix.Run(dt)

	// Prerender fill
	s.SystemList.TextureRenderAnimation.Run(dt)
	s.SystemList.TextureRenderFlip.Run(dt)
	s.SystemList.TextureRenderPosition.Run(dt)
	s.SystemList.TextureRenderRotation.Run(dt)
	s.SystemList.TextureRenderScale.Run(dt)
	s.SystemList.TextureRenderTint.Run(dt)

	// Render
	s.SystemList.AssetLib.Run(dt)
	shouldContinue := s.SystemList.Render.Run(dt)
	if !shouldContinue {
		s.Game.SetShouldDestroy(true)
		return
	}
}

func (s *MainScene) Destroy() {
	// Network intents
	s.SystemList.Network.Destroy()
	s.SystemList.NetworkReceive.Destroy()

	s.SystemList.Player.Destroy()

	// Network patches
	s.SystemList.NetworkSend.Destroy()

	// Animation
	s.SystemList.AnimationSpriteMatrix.Destroy()
	s.SystemList.AnimationPlayer.Destroy()

	// Prerender init
	s.SystemList.TextureRenderSprite.Destroy()
	s.SystemList.TextureRenderSpriteSheet.Destroy()
	s.SystemList.TextureRenderMatrix.Destroy()

	// Prerender fill
	s.SystemList.TextureRenderAnimation.Destroy()
	s.SystemList.TextureRenderFlip.Destroy()
	s.SystemList.TextureRenderPosition.Destroy()
	s.SystemList.TextureRenderRotation.Destroy()
	s.SystemList.TextureRenderScale.Destroy()
	s.SystemList.TextureRenderTint.Destroy()

	// Render
	s.SystemList.Debug.Destroy()
	s.SystemList.AssetLib.Destroy()
	s.SystemList.Render.Destroy()
}

func (s *MainScene) OnEnter() {

}

func (s *MainScene) OnExit() {

}

var _ gomp.AnyScene = (*MainScene)(nil)
