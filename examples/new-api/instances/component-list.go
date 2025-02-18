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

package instances

import (
	"gomp/examples/new-api/components"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
)

type ComponentList struct {
	Position        *stdcomponents.PositionComponentManager
	Rotation        *stdcomponents.RotationComponentManager
	Scale           *stdcomponents.ScaleComponentManager
	Velocity        *stdcomponents.VelocityComponentManager
	Flip            *stdcomponents.FlipComponentManager
	Sprite          *stdcomponents.SpriteComponentManager
	SpriteSheet     *stdcomponents.SpriteSheetComponentManager
	SpriteMatrix    *stdcomponents.SpriteMatrixComponentManager
	Tint            *stdcomponents.TintComponentManager
	AnimationPlayer *stdcomponents.AnimationPlayerComponentManager
	AnimationState  *stdcomponents.AnimationStateComponentManager
	TextureRender   *stdcomponents.TextureRenderComponentManager
	Network         *stdcomponents.NetworkComponentManager
	Health          *components.HealthComponentManager
}

func NewComponentList(world *ecs.World) ComponentList {
	return ComponentList{
		Position:        stdcomponents.NewPositionComponentManager(world),
		Rotation:        stdcomponents.NewRotationComponentManager(world),
		Scale:           stdcomponents.NewScaleComponentManager(world),
		Velocity:        stdcomponents.NewVelocityComponentManager(world),
		Flip:            stdcomponents.NewFlipComponentManager(world),
		Sprite:          stdcomponents.NewSpriteComponentManager(world),
		SpriteSheet:     stdcomponents.NewSpriteSheetComponentManager(world),
		SpriteMatrix:    stdcomponents.NewSpriteMatrixComponentManager(world),
		Tint:            stdcomponents.NewTintComponentManager(world),
		AnimationPlayer: stdcomponents.NewAnimationPlayerComponentManager(world),
		AnimationState:  stdcomponents.NewAnimationStateComponentManager(world),
		TextureRender:   stdcomponents.NewTextureRenderComponentManager(world),
		Network:         stdcomponents.NewNetworkComponentManager(world),
		Health:          components.NewHealthComponentManager(world),
	}
}
