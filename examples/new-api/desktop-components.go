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
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
)

func NewDesktopComponents(world *ecs.World) *desktopComponents {
	return &desktopComponents{
		position:        stdcomponents.NewPositionComponentManager(world),
		rotation:        stdcomponents.NewRotationComponentManager(world),
		scale:           stdcomponents.NewScaleComponentManager(world),
		flip:            stdcomponents.NewFlipComponentManager(world),
		sprite:          stdcomponents.NewSpriteComponentManager(world),
		spriteSheet:     stdcomponents.NewSpriteSheetComponentManager(world),
		spriteMatrix:    stdcomponents.NewSpriteMatrixComponentManager(world),
		tint:            stdcomponents.NewTintComponentManager(world),
		animationPlayer: stdcomponents.NewAnimationPlayerComponentManager(world),
		animationState:  stdcomponents.NewAnimationStateComponentManager(world),
		textureRender:   stdcomponents.NewTextureRenderComponentManager(world),
		network:         stdcomponents.NewNetworkComponentManager(world),
	}
}

type desktopComponents struct {
	position        *stdcomponents.PositionComponentManager
	rotation        *stdcomponents.RotationComponentManager
	scale           *stdcomponents.ScaleComponentManager
	flip            *stdcomponents.FlipComponentManager
	sprite          *stdcomponents.SpriteComponentManager
	spriteSheet     *stdcomponents.SpriteSheetComponentManager
	spriteMatrix    *stdcomponents.SpriteMatrixComponentManager
	tint            *stdcomponents.TintComponentManager
	animationPlayer *stdcomponents.AnimationPlayerComponentManager
	animationState  *stdcomponents.AnimationStateComponentManager
	textureRender   *stdcomponents.TextureRenderComponentManager
	network         *stdcomponents.NetworkComponentManager
}
