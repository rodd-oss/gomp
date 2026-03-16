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

package stdcomponents

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
	"image/color"
)

type Sprite struct {
	Texture *rl.Texture2D
	Frame   rl.Rectangle
	Origin  rl.Vector2
	Tint    color.RGBA
}

type SpriteComponentManager = ecs.ComponentManager[Sprite]

func NewSpriteComponentManager() SpriteComponentManager {
	return ecs.NewComponentManager[Sprite](SpriteComponentId)
}
