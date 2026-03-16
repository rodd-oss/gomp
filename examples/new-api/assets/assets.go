/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package assets

import (
	"embed"
	"gomp"
	"image/png"
	"log"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/negrel/assert"
)

//go:embed *.png
//go:embed *.wav
var fs embed.FS

var Textures = gomp.CreateAssetLibrary(
	func(path string) rl.Texture2D {
		assert.True(rl.IsWindowReady(), "Window is not initialized")

		file, err := fs.Open(path)
		if err != nil {
			log.Panic("Error opening file")
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			log.Panic("Error decoding file")
		}

		rlImg := rl.NewImageFromImage(img)
		return rl.LoadTextureFromImage(rlImg)
	},
	func(path string, asset *rl.Texture2D) {
		assert.True(rl.IsWindowReady(), "Window is not initialized")
		rl.UnloadTexture(*asset)
	},
)

var Audio = gomp.CreateAssetLibrary(
	func(path string) rl.Sound {
		file, err := fs.ReadFile(path)

		if err != nil {
			log.Panic("Error opening file")
		}

		fileTypeIndex := strings.LastIndex(path, ".")

		fileType := ".wav"

		if fileTypeIndex != -1 {
			fileType = path[fileTypeIndex:]
		}

		wave := rl.LoadWaveFromMemory(fileType, file, int32(len(file)))

		sound := rl.LoadSoundFromWave(wave)
		rl.UnloadWave(wave)

		assert.True(rl.IsSoundValid(sound), "Error loading sound")

		// Mute sound by default. Later it controlled by Aduio systems
		rl.SetSoundVolume(sound, 0)

		return sound
	},
	func(path string, asset *rl.Sound) {
		rl.UnloadSound(*asset)
	},
)
