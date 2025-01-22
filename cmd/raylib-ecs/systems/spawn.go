/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"fmt"
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type spawnController struct {
	pprofEnabled bool
}

const (
	minHpPercentage = 20
	minMaxHp        = 500
	maxMaxHp        = 2000
)

type resources struct {
	textures map[string]rl.Texture2D
}

func (r *resources) GetTexture(name string) rl.Texture2D {
	texture, ok := r.textures[name]
	if !ok {
		texture = rl.LoadTexture(name)
		r.textures[name] = texture
	}

	return texture
}

var Resources = &resources{
	textures: make(map[string]rl.Texture2D),
}

func (s *spawnController) Init(world *ecs.World) {}
func (s *spawnController) Update(world *ecs.World) {
	sprites := components.SpriteService.GetManager(world)
	healths := components.HealthService.GetManager(world)
	positions := components.PositionService.GetManager(world)
	rotations := components.RotationService.GetManager(world)
	scales := components.ScaleService.GetManager(world)

	if rl.IsKeyDown(rl.KeySpace) {
		for range rand.Intn(10000) {
			if world.Size() > 100_000_000 {
				break
			}

			newCreature := world.CreateEntity("Creature")

			// Adding position component
			t := components.Position{
				X: float32(rand.Int31n(800)),
				Y: float32(rand.Int31n(600)),
			}
			positions.Create(newCreature, t)

			// Adding rotation component
			rotation := components.Rotation{
				Angle: float32(rand.Int31n(360)),
			}
			rotations.Create(newCreature, rotation)

			// Adding scale component
			scale := components.Scale{
				X: 2,
				Y: 2,
			}
			scales.Create(newCreature, scale)

			// Adding HP component
			maxHp := minMaxHp + rand.Int31n(maxMaxHp-minMaxHp)
			hp := int32(float32(maxHp) * float32(minHpPercentage+rand.Int31n(100-minHpPercentage)) / 100)
			h := components.Health{
				Hp:    hp,
				MaxHp: maxHp,
			}
			healths.Create(newCreature, h)

			texture := Resources.GetTexture("assets/star.png")

			// Adding sprite component
			c := components.Sprite{
				Origin:  rl.Vector2{X: 0.5, Y: 0.5},
				Texture: texture,
				Frame:   rl.Rectangle{X: 0, Y: 0, Width: float32(texture.Width), Height: float32(texture.Height)},
			}

			sprites.Create(newCreature, c)
		}
	}

	if rl.IsKeyPressed(rl.KeyF9) {
		if s.pprofEnabled {
			pprof.StopCPUProfile()
			fmt.Println("CPU Profile Stopped")
		} else {
			f, err := os.Create("cpu.out")
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			fmt.Println("CPU Profile Started")
		}

		s.pprofEnabled = !s.pprofEnabled
	}
}
func (s *spawnController) FixedUpdate(world *ecs.World) {
}

func (s *spawnController) Destroy(world *ecs.World) {}
