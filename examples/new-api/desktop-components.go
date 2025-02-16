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
	rotation        *ecs.ComponentManager[stdcomponents.Rotation]
	scale           *ecs.ComponentManager[stdcomponents.Scale]
	flip            *ecs.ComponentManager[stdcomponents.Flip]
	sprite          *ecs.ComponentManager[stdcomponents.Sprite]
	spriteSheet     *ecs.ComponentManager[stdcomponents.SpriteSheet]
	spriteMatrix    *ecs.ComponentManager[stdcomponents.SpriteMatrix]
	tint            *ecs.ComponentManager[stdcomponents.Tint]
	animationPlayer *ecs.ComponentManager[stdcomponents.AnimationPlayer]
	animationState  *ecs.ComponentManager[stdcomponents.AnimationState]
	textureRender   *ecs.ComponentManager[stdcomponents.TextureRender]
	network         *ecs.ComponentManager[stdcomponents.Network]
}
