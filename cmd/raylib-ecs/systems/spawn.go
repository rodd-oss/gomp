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
	"image/color"
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

func (s *spawnController) Init(world *ecs.World) {}
func (s *spawnController) Run(world *ecs.World) {
	colors := components.ColorService.GetManager(world)
	healths := components.HealthService.GetManager(world)
	transforms := components.TransformService.GetManager(world)

	if rl.IsKeyPressed(rl.KeySpace) {
		for range rand.Intn(10000) {
			if world.Size() > 100_000_000 {
				break
			}

			newCreature := world.CreateEntity("Creature")

			t := components.Transform{
				X: rand.Int31n(1000),
				Y: rand.Int31n(1000),
			}

			transforms.Set(newCreature, t)

			maxHp := minMaxHp + rand.Int31n(maxMaxHp-minMaxHp)
			hp := int32(float32(maxHp) * float32(minHpPercentage+rand.Int31n(100-minHpPercentage)) / 100)

			h := components.Health{
				Hp:    hp,
				MaxHp: maxHp,
			}

			healths.Set(newCreature, h)

			c := color.RGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 0,
			}

			colors.Set(newCreature, c)
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

func (s *spawnController) Destroy(world *ecs.World) {}
