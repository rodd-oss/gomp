/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"fmt"
	ecs2 "gomp/pkg/ecs"
	"image/color"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type spawnSystem struct {
	transformComponent *ecs2.ComponentManager[transform]
	healthComponent    *ecs2.ComponentManager[health]
	colorComponent     *ecs2.ComponentManager[color.RGBA]
	movementComponent  *ecs2.ComponentManager[movement]
}

const (
	minHpPercentage = 20
	minMaxHp        = 500
	maxMaxHp        = 2000
)

var entityCount = 0
var pprofEnabled = false

func (s *spawnSystem) Init(world *ecs2.World) {
	s.transformComponent = transformComponentType.GetManager(world)
	s.healthComponent = healthComponentType.GetManager(world)
	s.colorComponent = colorComponentType.GetManager(world)
	s.movementComponent = movementComponentType.GetManager(world)
}
func (s *spawnSystem) Run(world *ecs2.World) {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		for range rand.Intn(1000) {

			newCreature := world.CreateEntity("Creature")

			t := transform{
				x: rand.Int31n(1000),
				y: rand.Int31n(1000),
			}

			s.transformComponent.Create(newCreature, t)

			maxHp := minMaxHp + rand.Int31n(maxMaxHp-minMaxHp)
			hp := int32(float32(maxHp) * float32(minHpPercentage+rand.Int31n(100-minHpPercentage)) / 100)

			h := health{
				hp:    hp,
				maxHp: maxHp,
			}

			s.healthComponent.Create(newCreature, h)

			c := color.RGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 0,
			}

			s.colorComponent.Create(newCreature, c)

			m := movement{
				goToX: t.x,
				goToY: t.y,
			}

			s.movementComponent.Create(newCreature, m)

			entityCount++
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF9) {
		if *cpuprofile != "" {
			if pprofEnabled {
				pprof.StopCPUProfile()
				fmt.Println("CPU Profile Stopped")
			} else {
				f, err := os.Create(*cpuprofile)
				if err != nil {
					log.Fatal(err)
				}
				pprof.StartCPUProfile(f)
				fmt.Println("CPU Profile Started")
			}

			pprofEnabled = !pprofEnabled
		}
	}
}
func (s *spawnSystem) Destroy(world *ecs2.World) {}
