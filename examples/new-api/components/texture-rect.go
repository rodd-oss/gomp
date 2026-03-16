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

type TextureRect struct {
	Dest     rl.Rectangle
	Origin   rl.Vector2
	Rotation float32
	Color    color.RGBA
}

type TextureRectComponentManager = ecs.ComponentManager[TextureRect]

func NewTextureRectComponentManager() TextureRectComponentManager {
	return ecs.NewComponentManager[TextureRect](TextureRectComponentId)
}
