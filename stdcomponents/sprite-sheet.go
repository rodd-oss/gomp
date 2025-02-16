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
)

type SpriteSheet struct {
	Texture     *rl.Texture2D
	Frame       rl.Rectangle
	Origin      rl.Vector2
	NumOfFrames int32
	FPS         int32
	Vertical    bool
}

type SpriteSheetComponentManager = ecs.ComponentManager[SpriteSheet]

func NewSpriteSheetComponentManager(world *ecs.World) *ecs.ComponentManager[SpriteSheet] {
	return ecs.NewComponentManager[SpriteSheet](world, SPRITE_SHEET_ID)
}
