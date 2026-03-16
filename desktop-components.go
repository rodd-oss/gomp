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
	"gomp/stdcomponents"
)

func NewDesktopComponents() DesktopComponents {
	return DesktopComponents{
		Position:        stdcomponents.NewPositionComponentManager(),
		Rotation:        stdcomponents.NewRotationComponentManager(),
		Scale:           stdcomponents.NewScaleComponentManager(),
		Velocity:        stdcomponents.NewVelocityComponentManager(),
		Flip:            stdcomponents.NewFlipComponentManager(),
		Sprite:          stdcomponents.NewSpriteComponentManager(),
		SpriteSheet:     stdcomponents.NewSpriteSheetComponentManager(),
		SpriteMatrix:    stdcomponents.NewSpriteMatrixComponentManager(),
		Tint:            stdcomponents.NewTintComponentManager(),
		AnimationPlayer: stdcomponents.NewAnimationPlayerComponentManager(),
		AnimationState:  stdcomponents.NewAnimationStateComponentManager(),
		RlTexturePro:    stdcomponents.NewRlTextureProComponentManager(),
		Network:         stdcomponents.NewNetworkComponentManager(),
	}
}

type DesktopComponents struct {
	Position        stdcomponents.PositionComponentManager
	Rotation        stdcomponents.RotationComponentManager
	Scale           stdcomponents.ScaleComponentManager
	Velocity        stdcomponents.VelocityComponentManager
	Flip            stdcomponents.FlipComponentManager
	Sprite          stdcomponents.SpriteComponentManager
	SpriteSheet     stdcomponents.SpriteSheetComponentManager
	SpriteMatrix    stdcomponents.SpriteMatrixComponentManager
	Tint            stdcomponents.TintComponentManager
	AnimationPlayer stdcomponents.AnimationPlayerComponentManager
	AnimationState  stdcomponents.AnimationStateComponentManager
	RlTexturePro    stdcomponents.RLTextureProComponentManager
	Network         stdcomponents.NetworkComponentManager
}
