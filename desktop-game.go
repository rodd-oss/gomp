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

package gomp

import (
	"github.com/negrel/assert"
	"time"
)

func NewDesktopGame(systems *DesktopSystems, scenes map[SceneId]AnyScene) DesktopGame {
	game := DesktopGame{
		Systems: systems,
		Scenes:  scenes,
	}

	return game
}

type DesktopGame struct {
	Scenes         map[SceneId]AnyScene
	CurrentSceneId SceneId

	Systems       *DesktopSystems
	shouldDestroy bool
}

func (g *DesktopGame) Init() {
	// Network receive
	g.Systems.Network.Init()
	g.Systems.NetworkReceive.Init()

	// Scenes
	scene, ok := g.Scenes[g.CurrentSceneId]
	assert.True(ok, "Scene not found")
	scene.Init()

	// Network patches
	g.Systems.NetworkSend.Init()

	// Animation
	g.Systems.AnimationSpriteMatrix.Init()
	g.Systems.AnimationPlayer.Init()

	// Prerender init
	g.Systems.TextureRenderSprite.Init()
	g.Systems.TextureRenderSpriteSheet.Init()
	g.Systems.TextureRenderMatrix.Init()

	// Prerender fill
	g.Systems.TextureRenderAnimation.Init()
	g.Systems.TextureRenderFlip.Init()
	g.Systems.TextureRenderPosition.Init()
	g.Systems.TextureRenderRotation.Init()
	g.Systems.TextureRenderScale.Init()
	g.Systems.TextureRenderTint.Init()

	// Render
	g.Systems.Render.Init()
	g.Systems.Debug.Init()
	g.Systems.AssetLib.Init()
}

func (g *DesktopGame) Update(dt time.Duration) {
	// Network receive
	g.Systems.Network.Run(dt)
	g.Systems.NetworkReceive.Run(dt)

	// Scenes
	scene, ok := g.Scenes[g.CurrentSceneId]
	assert.True(ok, "Scene not found")
	g.CurrentSceneId = scene.Update(dt)

	// Network patches
	g.Systems.NetworkSend.Run(dt)

	g.Systems.Debug.Run(dt)
}

func (g *DesktopGame) FixedUpdate(dt time.Duration) {
	// Scenes
	scene, ok := g.Scenes[g.CurrentSceneId]
	assert.True(ok, "Scene not found")
	scene.FixedUpdate(dt)

	// Network send
	g.Systems.NetworkSend.Run(dt)
}

func (g *DesktopGame) Render(dt time.Duration) {
	g.Systems.Velocity.Run(dt)

	// Animation
	g.Systems.AnimationSpriteMatrix.Run(dt)
	g.Systems.AnimationPlayer.Run(dt)

	// Prerender init
	g.Systems.TextureRenderSprite.Run(dt)
	g.Systems.TextureRenderSpriteSheet.Run(dt)
	g.Systems.TextureRenderMatrix.Run(dt)

	// Prerender fill
	g.Systems.TextureRenderAnimation.Run(dt)
	g.Systems.TextureRenderFlip.Run(dt)
	g.Systems.TextureRenderPosition.Run(dt)
	g.Systems.TextureRenderRotation.Run(dt)
	g.Systems.TextureRenderScale.Run(dt)
	g.Systems.TextureRenderTint.Run(dt)

	// Render
	g.Systems.AssetLib.Run(dt)
	if err := g.Systems.Render.Run(dt); err != nil {
		g.shouldDestroy = true
	}
}

func (g *DesktopGame) Destroy() {
	// Network intents
	g.Systems.Network.Destroy()
	g.Systems.NetworkReceive.Destroy()

	// Scenes
	scene, ok := g.Scenes[g.CurrentSceneId]
	assert.True(ok, "Scene not found")
	scene.Destroy()

	// Network patches
	g.Systems.NetworkSend.Destroy()

	// Animation
	g.Systems.AnimationSpriteMatrix.Destroy()
	g.Systems.AnimationPlayer.Destroy()

	// Prerender init
	g.Systems.TextureRenderSprite.Destroy()
	g.Systems.TextureRenderSpriteSheet.Destroy()
	g.Systems.TextureRenderMatrix.Destroy()

	// Prerender fill
	g.Systems.TextureRenderAnimation.Destroy()
	g.Systems.TextureRenderFlip.Destroy()
	g.Systems.TextureRenderPosition.Destroy()
	g.Systems.TextureRenderRotation.Destroy()
	g.Systems.TextureRenderScale.Destroy()
	g.Systems.TextureRenderTint.Destroy()

	// Render
	g.Systems.Debug.Destroy()
	g.Systems.AssetLib.Destroy()
	g.Systems.Render.Destroy()
}

func (g *DesktopGame) ShouldDestroy() bool {
	return g.shouldDestroy
}
