/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package gomp

import (
	"math/rand"
	"testing"
)

type Float3 struct {
	X, Y, Z float64
}

// Array of Structures (AoS)
type Unit struct {
	Pos      Float3
	Velocity Float3
	HP       int
	Offset   [AdditionalDataPer64]int64
}

type UnitsAoS struct {
	Units []Unit
}

// Structure of Arrays (SoA)
type UnitsSoA struct {
	Positions  [][pageSize]Float3
	Velocities [][pageSize]Float3
	HPs        [][pageSize]int
	Offsets    [][pageSize][AdditionalDataPer64]int64
}

func UpdatePosition(pos *Float3, velocity *Float3) {
	pos.X += velocity.X
	pos.Y += velocity.Y
	pos.Z += velocity.Z
}

func TakeDamage(hp *int, damage int) {
	*hp = max(0, *hp-damage)
}

func UpdateOffset(offset *[AdditionalDataPer64]int64) {
	for i := range offset {
		offset[i]++
	}
}

func ShouldRender(pos *Float3) bool {
	return pos.X > 10.0 && pos.Y < 5.0 && pos.Z > 0.0
}

func GenerateUnitsAoS(count int) UnitsAoS {
	units := UnitsAoS{Units: make([]Unit, count)}
	for i := 0; i < count; i++ {
		units.Units[i] = Unit{
			Pos: Float3{
				X: rand.Float64()*40 - 20,
				Y: rand.Float64()*40 - 20,
				Z: rand.Float64()*40 - 20,
			},
			Velocity: Float3{
				X: rand.Float64()*2 - 1,
				Y: rand.Float64()*2 - 1,
				Z: rand.Float64()*2 - 1,
			},
			HP: rand.Intn(151) + 50,
		}
	}
	return units
}

func MakeUnitsSoA(aos UnitsAoS) UnitsSoA {
	soa := UnitsSoA{
		Positions:  make([][pageSize]Float3, len(aos.Units)/pageSize+1),
		Velocities: make([][pageSize]Float3, len(aos.Units)/pageSize+1),
		HPs:        make([][pageSize]int, len(aos.Units)/pageSize+1),
		Offsets:    make([][pageSize][AdditionalDataPer64]int64, len(aos.Units)/pageSize+1),
	}

	for i, unit := range aos.Units {
		a, b := i>>pageShift, i%pageSize
		soa.Positions[a][b] = unit.Pos
		soa.Velocities[a][b] = unit.Velocity
		soa.HPs[a][b] = unit.HP
		soa.Offsets[a][b] = unit.Offset
	}
	return soa
}

func PositionSystem(positions [][pageSize]Float3, velocities [][pageSize]Float3) {
	for i := 0; i < len(positions); i++ {
		for j := 0; j < len(positions[i]); j++ {
			UpdatePosition(&positions[i][j], &velocities[i][j])
		}
	}
}

func TakeDamageSystem(hps [][pageSize]int) {
	for i := 0; i < len(hps); i++ {
		for j := 0; j < len(hps[i]); j++ {
			TakeDamage(&hps[i][j], 33)
		}
	}
}

func OffsetSystem(offsets [][pageSize][AdditionalDataPer64]int64) {
	for i := 0; i < len(offsets); i++ {
		for j := 0; j < len(offsets[i]); j++ {
			UpdateOffset(&offsets[i][j])
		}
	}
}

func renderSystem(positions [][pageSize]Float3, unitsToRender []int) {
	for i := 0; i < len(positions); i++ {
		for j := 0; j < len(positions[i]); j++ {
			if ShouldRender(&positions[i][j]) {
				index := j + i*len(positions[i])
				unitsToRender = append(unitsToRender, index)
			}
		}
	}
}

// RuDimbo Team
func BenchmarkUpdateRuDimbo(b *testing.B) {
	b.ReportAllocs()

	units := GenerateUnitsAoS(NumOfEntities)
	unitsToRender := make([]int, 0, NumOfEntities)

	for b.Loop() {
		for i := range units.Units {
			unit := &units.Units[i]
			UpdatePosition(&unit.Pos, &unit.Velocity)
			if ShouldRender(&unit.Pos) {
				unitsToRender = append(unitsToRender, i)
			}
			TakeDamage(&unit.HP, 33)
			UpdateOffset(&unit.Offset)
		}
		unitsToRender = unitsToRender[:0]
	}
}

// Rodd Team
func BenchmarkUpdateRodd(b *testing.B) {
	b.ReportAllocs()
	soaUnits := MakeUnitsSoA(GenerateUnitsAoS(NumOfEntities))
	unitsToRender := make([]int, 0, NumOfEntities)
	for b.Loop() {
		positions := soaUnits.Positions
		velocities := soaUnits.Velocities
		hps := soaUnits.HPs
		offsets := soaUnits.Offsets
		PositionSystem(positions, velocities)
		TakeDamageSystem(hps)
		OffsetSystem(offsets)
		renderSystem(positions, unitsToRender)
		unitsToRender = unitsToRender[:0]
	}
}

const NumOfEntities = 10_000_000
const AdditionalDataPer64 = 30
const pageShift = 10
const pageSize = 1 << pageShift
