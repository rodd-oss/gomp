/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/raylib-ecs/components"
	ecs2 "gomp/pkgs/ecs"
)

// TextureRenderSprite is a system that prepares SpriteSheet to be rendered
type trSpriteSheetController struct{}

func (s *trSpriteSheetController) Init(world *ecs2.World)        {}
func (s *trSpriteSheetController) FixedUpdate(world *ecs2.World) {}
func (s *trSpriteSheetController) Update(world *ecs2.World) {
	// Get component managers
	spriteSheets := components.SpriteSheetService.GetManager(world)
	textureRenders := components.TextureRenderService.GetManager(world)

	// Update sprites and spriteRenders
	spriteSheets.AllParallel(func(entity ecs2.Entity, spriteSheet *components.SpriteSheet) bool {
		if spriteSheet == nil {
			return true
		}

		tr := textureRenders.Get(entity)
		if tr == nil {
			// Create new spriteRender
			newRender := components.TextureRender{
				Texture: spriteSheet.Texture,
				Frame:   spriteSheet.Frame,
				Origin:  spriteSheet.Origin,
				Dest: rl.NewRectangle(
					0,
					0,
					spriteSheet.Frame.Width,
					spriteSheet.Frame.Height,
				),
			}

			textureRenders.Create(entity, newRender)
		} else {
			// Update spriteRender
			tr.Texture = spriteSheet.Texture
			tr.Frame = spriteSheet.Frame
			tr.Origin = spriteSheet.Origin
		}
		return true
	})
}
func (s *trSpriteSheetController) Destroy(world *ecs2.World) {}
