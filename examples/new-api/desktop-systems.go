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

package main

import (
	"gomp/examples/new-api/assets"
	"gomp/examples/new-api/systems"
	"gomp/pkg/ecs"
	"gomp/stdsystems"
	"time"
)

func NewDesktopSystems(world *ecs.World, components *desktopComponents) *desktopSystems {
	return &desktopSystems{
		player: systems.NewPlayerSystem(world, components.spriteMatrix, components.position, components.rotation, components.scale, components.animationPlayer, components.animationState, components.tint, components.flip),

		debug: stdsystems.NewDebugSystem(),

		network:        stdsystems.NewNetworkSystem(),
		networkReceive: stdsystems.NewNetworkReceiveSystem(),
		networkSend:    stdsystems.NewNetworkSendSystem(world, components.position, components.rotation, components.flip),

		animationSpriteMatrix: stdsystems.NewAnimationSpriteMatrixSystem(world, components.animationPlayer, components.animationState, components.spriteMatrix),
		animationPlayer:       stdsystems.NewAnimationPlayerSystem(components.animationPlayer),

		textureRenderSprite:      stdsystems.NewTextureRenderSpriteSystem(components.sprite, components.textureRender),
		textureRenderSpriteSheet: stdsystems.NewTextureRenderSpriteSheetSystem(components.spriteSheet, components.textureRender),
		textureRenderMatrix:      stdsystems.NewTextureRenderMatrixSystem(components.spriteMatrix, components.textureRender, components.animationState),

		textureRenderAnimation: stdsystems.NewTextureRenderAnimationSystem(components.animationPlayer, components.textureRender),
		textureRenderFlip:      stdsystems.NewTextureRenderFlipSystem(components.flip, components.textureRender),
		textureRenderPosition:  stdsystems.NewTextureRenderPositionSystem(components.position, components.textureRender),
		textureRenderRotation:  stdsystems.NewTextureRenderRotationSystem(components.rotation, components.textureRender),
		textureRenderScale:     stdsystems.NewTextureRenderScaleSystem(components.scale, components.textureRender),
		textureRenderTint:      stdsystems.NewTextureRenderTintSystem(components.tint, components.textureRender),

		assetLib: stdsystems.NewAssetLibSystem([]ecs.AnyAssetLibrary{assets.Textures}),
		render:   stdsystems.NewRenderSystem(world, components.textureRender),
	}
}

type desktopSystems struct {
	player *systems.PlayerSystem

	debug *stdsystems.DebugSystem
	// Network
	network        *stdsystems.NetworkSystem
	networkReceive *stdsystems.NetworkReceiveSystem
	networkSend    *stdsystems.NetworkSendSystem
	// Animation
	animationSpriteMatrix *stdsystems.AnimationSpriteMatrixSystem
	animationPlayer       *stdsystems.AnimationPlayerSystem
	// Prerender init
	textureRenderSprite      *stdsystems.TextureRenderSpriteSystem
	textureRenderSpriteSheet *stdsystems.TextureRenderSpriteSheetSystem
	textureRenderMatrix      *stdsystems.TextureRenderMatrixSystem
	// Prerender fill
	textureRenderAnimation *stdsystems.TextureRenderAnimationSystem
	textureRenderFlip      *stdsystems.TextureRenderFlipSystem
	textureRenderPosition  *stdsystems.TextureRenderPositionSystem
	textureRenderRotation  *stdsystems.TextureRenderRotationSystem
	textureRenderScale     *stdsystems.TextureRenderScaleSystem
	textureRenderTint      *stdsystems.TextureRenderTintSystem
	// Render
	assetLib *stdsystems.AssetLibSystem
	render   *stdsystems.RenderSystem
}

func (s *desktopSystems) Init() {
	// Network receive
	s.network.Init()
	s.networkReceive.Init()

	// Network patches
	s.networkSend.Init()

	// Scenes
	s.player.Init()

	// Animation
	s.animationSpriteMatrix.Init()
	s.animationPlayer.Init()

	// Prerender init
	s.textureRenderSprite.Init()
	s.textureRenderSpriteSheet.Init()
	s.textureRenderMatrix.Init()

	// Prerender fill
	s.textureRenderAnimation.Init()
	s.textureRenderFlip.Init()
	s.textureRenderPosition.Init()
	s.textureRenderRotation.Init()
	s.textureRenderScale.Init()
	s.textureRenderTint.Init()

	// Render
	s.render.Init()
	s.debug.Init()
	s.assetLib.Init()
}

func (s *desktopSystems) Update(dt time.Duration) {
	// Network receive
	s.network.Run(dt)
	s.networkReceive.Run(dt)

	// Scenes
	s.player.Run(dt)

	// Animation
	s.animationSpriteMatrix.Run(dt)
	s.animationPlayer.Run(dt)

	// Prerender init
	s.textureRenderSprite.Run(dt)
	s.textureRenderSpriteSheet.Run(dt)
	s.textureRenderMatrix.Run(dt)

	// Prerender fill
	s.textureRenderAnimation.Run(dt)
	s.textureRenderFlip.Run(dt)
	s.textureRenderPosition.Run(dt)
	s.textureRenderRotation.Run(dt)
	s.textureRenderScale.Run(dt)
	s.textureRenderTint.Run(dt)

	// Render
	s.assetLib.Run(dt)
	s.debug.Run(dt)
	s.render.Run(dt)
}

func (s *desktopSystems) FixedUpdate(dt time.Duration) {
	// Scenes
	s.player.Run(dt)

	// Network send
	s.networkSend.Run(dt)
}

func (s *desktopSystems) Destroy() {
	// Network intents
	s.network.Destroy()
	s.networkReceive.Destroy()

	// Scenes
	s.player.Destroy()

	// Network patches
	s.networkSend.Destroy()

	// Animation
	s.animationSpriteMatrix.Destroy()
	s.animationPlayer.Destroy()

	// Prerender init
	s.textureRenderSprite.Destroy()
	s.textureRenderSpriteSheet.Destroy()
	s.textureRenderMatrix.Destroy()

	// Prerender fill
	s.textureRenderAnimation.Destroy()
	s.textureRenderFlip.Destroy()
	s.textureRenderPosition.Destroy()
	s.textureRenderRotation.Destroy()
	s.textureRenderScale.Destroy()
	s.textureRenderTint.Destroy()

	// Render
	s.debug.Destroy()
	s.assetLib.Destroy()
	s.render.Destroy()
}
