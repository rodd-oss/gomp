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
	"gomp/vectors"
	"image/color"
)

type RLTexturePro struct {
	Texture  *rl.Texture2D
	Frame    rl.Rectangle
	Origin   rl.Vector2
	Tint     color.RGBA
	Dest     rl.Rectangle
	Rotation float32
}

func (t *RLTexturePro) Rect() vectors.Rectangle {
	return vectors.Rectangle{
		X:      t.Dest.X - t.Origin.X,
		Y:      t.Dest.Y - t.Origin.Y,
		Width:  t.Dest.Width,
		Height: t.Dest.Height,
	}
}

type RLTextureProComponentManager = ecs.ComponentManager[RLTexturePro]

func NewRlTextureProComponentManager() RLTextureProComponentManager {
	return ecs.NewComponentManager[RLTexturePro](RLTextureProComponentId)
}
