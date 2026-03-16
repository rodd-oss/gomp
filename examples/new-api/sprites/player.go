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

package sprites

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/new-api/assets"
	"gomp/stdcomponents"
)

var PlayerSpriteMatrix = stdcomponents.SpriteMatrix{
	Texture: assets.Textures.Get("milansheet.png"),
	Origin:  rl.Vector2{X: 0.5, Y: 0.5},
	Dest:    rl.Rectangle{X: 0, Y: 0, Width: 96, Height: 128},
	FPS:     12,
	Animations: []stdcomponents.SpriteMatrixAnimation{
		{
			Name:        "idle",
			Frame:       rl.Rectangle{X: 0, Y: 0, Width: 96, Height: 128},
			NumOfFrames: 1,
			Vertical:    false,
			Loop:        true,
		},
		{
			Name:        "walk",
			Frame:       rl.Rectangle{X: 0, Y: 512, Width: 96, Height: 128},
			NumOfFrames: 8,
			Vertical:    false,
			Loop:        true,
		},
		{
			Name:        "jump",
			Frame:       rl.Rectangle{X: 96, Y: 0, Width: 96, Height: 128},
			NumOfFrames: 1,
			Vertical:    false,
			Loop:        false,
		},
	},
}
