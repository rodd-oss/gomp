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
	"gomp/examples/new-api/assets"
	"gomp/pkg/ecs"
	"gomp/stdsystems"
	"time"
)

func NewDesktopSystems(world *ecs.World, components *DesktopComponents, scenes []Scene) *DesktopSystems {
	return &DesktopSystems{
		Scenes: scenes,
		Debug:  stdsystems.NewDebugSystem(),

		Network:        stdsystems.NewNetworkSystem(),
		NetworkReceive: stdsystems.NewNetworkReceiveSystem(),
		NetworkSend:    stdsystems.NewNetworkSendSystem(world, components.Position, components.Rotation, components.Flip),

		AnimationSpriteMatrix: stdsystems.NewAnimationSpriteMatrixSystem(world, components.AnimationPlayer, components.AnimationState, components.SpriteMatrix),
		AnimationPlayer:       stdsystems.NewAnimationPlayerSystem(components.AnimationPlayer),

		TextureRenderSprite:      stdsystems.NewTextureRenderSpriteSystem(components.Sprite, components.TextureRender),
		TextureRenderSpriteSheet: stdsystems.NewTextureRenderSpriteSheetSystem(components.SpriteSheet, components.TextureRender),
		TextureRenderMatrix:      stdsystems.NewTextureRenderMatrixSystem(components.SpriteMatrix, components.TextureRender, components.AnimationState),

		TextureRenderAnimation: stdsystems.NewTextureRenderAnimationSystem(components.AnimationPlayer, components.TextureRender),
		TextureRenderFlip:      stdsystems.NewTextureRenderFlipSystem(components.Flip, components.TextureRender),
		TextureRenderPosition:  stdsystems.NewTextureRenderPositionSystem(components.Position, components.TextureRender),
		TextureRenderRotation:  stdsystems.NewTextureRenderRotationSystem(components.Rotation, components.TextureRender),
		TextureRenderScale:     stdsystems.NewTextureRenderScaleSystem(components.Scale, components.TextureRender),
		TextureRenderTint:      stdsystems.NewTextureRenderTintSystem(components.Tint, components.TextureRender),

		AssetLib: stdsystems.NewAssetLibSystem([]ecs.AnyAssetLibrary{assets.Textures}),
		Render:   stdsystems.NewRenderSystem(world, components.TextureRender),
	}
}

type DesktopSystems struct {
	Scenes []Scene
	Debug  *stdsystems.DebugSystem
	// Network
	Network        *stdsystems.NetworkSystem
	NetworkReceive *stdsystems.NetworkReceiveSystem
	NetworkSend    *stdsystems.NetworkSendSystem
	// Animation
	AnimationSpriteMatrix *stdsystems.AnimationSpriteMatrixSystem
	AnimationPlayer       *stdsystems.AnimationPlayerSystem
	// Prerender init
	TextureRenderSprite      *stdsystems.TextureRenderSpriteSystem
	TextureRenderSpriteSheet *stdsystems.TextureRenderSpriteSheetSystem
	TextureRenderMatrix      *stdsystems.TextureRenderMatrixSystem
	// Prerender fill
	TextureRenderAnimation *stdsystems.TextureRenderAnimationSystem
	TextureRenderFlip      *stdsystems.TextureRenderFlipSystem
	TextureRenderPosition  *stdsystems.TextureRenderPositionSystem
	TextureRenderRotation  *stdsystems.TextureRenderRotationSystem
	TextureRenderScale     *stdsystems.TextureRenderScaleSystem
	TextureRenderTint      *stdsystems.TextureRenderTintSystem
	// Render
	AssetLib *stdsystems.AssetLibSystem
	Render   *stdsystems.RenderSystem
}

func (s *DesktopSystems) Init() {
	// Network receive
	s.Network.Init()
	s.NetworkReceive.Init()

	// Network patches
	s.NetworkSend.Init()

	// Scenes
	for i := range s.Scenes {
		if s.Scenes[i].Enabled {
			s.Scenes[i].Init()
		}
	}

	// Animation
	s.AnimationSpriteMatrix.Init()
	s.AnimationPlayer.Init()

	// Prerender init
	s.TextureRenderSprite.Init()
	s.TextureRenderSpriteSheet.Init()
	s.TextureRenderMatrix.Init()

	// Prerender fill
	s.TextureRenderAnimation.Init()
	s.TextureRenderFlip.Init()
	s.TextureRenderPosition.Init()
	s.TextureRenderRotation.Init()
	s.TextureRenderScale.Init()
	s.TextureRenderTint.Init()

	// Render
	s.Render.Init()
	s.Debug.Init()
	s.AssetLib.Init()
}

func (s *DesktopSystems) Update(dt time.Duration) {
	// Network receive
	s.Network.Run(dt)
	s.NetworkReceive.Run(dt)

	// Scenes
	for i := range s.Scenes {
		if s.Scenes[i].Enabled {
			s.Scenes[i].Update(dt)
		}
	}

	// Animation
	s.AnimationSpriteMatrix.Run(dt)
	s.AnimationPlayer.Run(dt)

	// Prerender init
	s.TextureRenderSprite.Run(dt)
	s.TextureRenderSpriteSheet.Run(dt)
	s.TextureRenderMatrix.Run(dt)

	// Prerender fill
	s.TextureRenderAnimation.Run(dt)
	s.TextureRenderFlip.Run(dt)
	s.TextureRenderPosition.Run(dt)
	s.TextureRenderRotation.Run(dt)
	s.TextureRenderScale.Run(dt)
	s.TextureRenderTint.Run(dt)

	// Render
	s.AssetLib.Run(dt)
	s.Debug.Run(dt)
	s.Render.Run(dt)
}

func (s *DesktopSystems) FixedUpdate(dt time.Duration) {
	// Scenes
	for i := range s.Scenes {
		if s.Scenes[i].Enabled {
			s.Scenes[i].FixedUpdate(dt)
		}
	}

	// Network send
	s.NetworkSend.Run(dt)
}

func (s *DesktopSystems) Destroy() {
	// Network intents
	s.Network.Destroy()
	s.NetworkReceive.Destroy()

	// Scenes
	for i := range s.Scenes {
		if s.Scenes[i].IsInitialized {
			s.Scenes[i].Destroy()
		}
	}

	// Network patches
	s.NetworkSend.Destroy()

	// Animation
	s.AnimationSpriteMatrix.Destroy()
	s.AnimationPlayer.Destroy()

	// Prerender init
	s.TextureRenderSprite.Destroy()
	s.TextureRenderSpriteSheet.Destroy()
	s.TextureRenderMatrix.Destroy()

	// Prerender fill
	s.TextureRenderAnimation.Destroy()
	s.TextureRenderFlip.Destroy()
	s.TextureRenderPosition.Destroy()
	s.TextureRenderRotation.Destroy()
	s.TextureRenderScale.Destroy()
	s.TextureRenderTint.Destroy()

	// Render
	s.Debug.Destroy()
	s.AssetLib.Destroy()
	s.Render.Destroy()
}
