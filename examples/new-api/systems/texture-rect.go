/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/negrel/assert"
	"gomp/examples/new-api/components"
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"gomp/stdcomponents"
	"runtime"
	"time"
)

func NewTextureRectSystem() TextureRectSystem {
	return TextureRectSystem{}
}

type TextureRectSystem struct {
	TextureRect *components.TextureRectComponentManager
	Textures    *stdcomponents.RLTextureProComponentManager
	texture     rl.RenderTexture2D
	numWorkers  int
	Engine      *core.Engine
}

func (s *TextureRectSystem) Init() {
	var texture = rl.LoadRenderTexture(1, 1) // :)
	rl.BeginTextureMode(texture)
	rl.ClearBackground(rl.White)
	rl.EndTextureMode()
	s.texture = texture
	s.numWorkers = runtime.NumCPU() - 2
}

func (s *TextureRectSystem) Run(dt time.Duration) {
	// Create shallow copy of texture to draw rectangles
	s.TextureRect.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		rect := s.TextureRect.GetUnsafe(entity)
		assert.NotNil(rect, "rect is nil; entity: %d", entity)
		texture := s.Textures.GetUnsafe(entity)
		assert.NotNil(texture, "texture is nil; entity: %d", entity)

		texture.Texture = &s.texture.Texture
		texture.Dest = rect.Dest
		texture.Rotation = rect.Rotation
		texture.Origin = rect.Origin
		texture.Tint = rect.Color
	})
}

func (s *TextureRectSystem) Destroy() {
	rl.UnloadRenderTexture(s.texture)
}
