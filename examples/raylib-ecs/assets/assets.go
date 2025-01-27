/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package assets

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkgs/ecs"
)

var Textures = ecs.CreateAssetLibrary(
	func(path string) rl.Texture2D {
		return rl.LoadTexture(path)
	},
	func(path string, asset *rl.Texture2D) {
		rl.UnloadTexture(*asset)
	},
)
