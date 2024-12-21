/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"gomp_game/pkgs/gomp/ecs"
	"image/color"
	"math/rand"
)

type spawnSystem struct {
	transformComponent ecs.WorldComponents[transform]
	healthComponent    ecs.WorldComponents[health]
	colorComponent     ecs.WorldComponents[color.RGBA]
	movementComponent  ecs.WorldComponents[movement]
}

const (
	minHpPercentage = 20
	minMaxHp        = 500
	maxMaxHp        = 2000
)

var entityCount = 0

func (s *spawnSystem) Init(world *ecs.World) {
	s.transformComponent = transformComponentType.Instances(world)
	s.healthComponent = healthComponentType.Instances(world)
	s.colorComponent = colorComponentType.Instances(world)
	s.movementComponent = movementComponentType.Instances(world)
}
func (s *spawnSystem) Run(world *ecs.World) {
	for range rand.Intn(1000) {

		newCreature := world.CreateEntity("Creature")

		t := transform{
			x: rand.Int31n(1000),
			y: rand.Int31n(1000),
		}

		s.transformComponent.Set(newCreature, t)

		maxHp := minMaxHp + rand.Int31n(maxMaxHp-minMaxHp)
		hp := int32(float32(maxHp) * float32(minHpPercentage+rand.Int31n(100-minHpPercentage)) / 100)

		h := health{
			hp:    hp,
			maxHp: maxHp,
		}

		s.healthComponent.Set(newCreature, h)

		c := color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}

		s.colorComponent.Set(newCreature, c)

		m := movement{
			goToX: t.x,
			goToY: t.y,
		}

		s.movementComponent.Set(newCreature, m)

		entityCount++
	}
}
func (s *spawnSystem) Destroy(world *ecs.World) {}
