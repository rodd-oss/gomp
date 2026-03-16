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

package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
	"image/color"
)

type TextureCircle struct {
	CenterX  float32
	CenterY  float32
	Radius   float32
	Rotation float32
	Origin   rl.Vector2
	Color    color.RGBA
}

type PrimitiveCircleComponentManager = ecs.ComponentManager[TextureCircle]

func NewTextureCircleComponentManager() PrimitiveCircleComponentManager {
	return ecs.NewComponentManager[TextureCircle](TextureCircleComponentId)
}
