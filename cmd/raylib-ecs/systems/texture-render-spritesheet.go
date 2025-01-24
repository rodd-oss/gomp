/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TextureRenderSprite is a system that prepares SpriteSheet to be rendered
type trSpriteSheetController struct{}

func (s *trSpriteSheetController) Init(world *ecs.World)        {}
func (s *trSpriteSheetController) FixedUpdate(world *ecs.World) {}
func (s *trSpriteSheetController) Update(world *ecs.World) {
	// Get component managers
	spriteSheets := components.SpriteSheetService.GetManager(world)
	textureRenders := components.TextureRenderService.GetManager(world)

	// Update sprites and spriteRenders
	spriteSheets.AllParallel(func(entity ecs.Entity, spriteSheet *components.SpriteSheet) bool {
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
				Tint:    spriteSheet.Tint,
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
			tr.Frame = spriteSheet.Frame
			tr.Origin = spriteSheet.Origin
			tr.Tint = spriteSheet.Tint
		}
		return true
	})
}
func (s *trSpriteSheetController) Destroy(world *ecs.World) {}
