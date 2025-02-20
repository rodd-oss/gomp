/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package assets

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/negrel/assert"
	"gomp"
)

var Textures = gomp.CreateAssetLibrary(
	func(path string) rl.Texture2D {
		assert.True(rl.IsWindowReady(), "Window is not initialized")
		return rl.LoadTexture(path)
	},
	func(path string, asset *rl.Texture2D) {
		assert.True(rl.IsWindowReady(), "Window is not initialized")
		rl.UnloadTexture(*asset)
	},
)
