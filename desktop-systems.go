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
)

func NewDesktopSystems(world *ecs.World, components *DesktopComponents) DesktopSystems {
	return DesktopSystems{
		Debug: stdsystems.NewDebugSystem(),

		Velocity: stdsystems.NewVelocitySystem(components.Velocity, components.Position),

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
	Debug *stdsystems.DebugSystem

	Velocity *stdsystems.VelocitySystem

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
