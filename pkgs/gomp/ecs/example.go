/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type Transform struct {
	X, Y, Z float32
}

type Rotation struct {
	RX, RY, RZ int
}

type Scale struct {
	Value float32
}

var _ = CreateComponent[Rotation]()
var transformComponent = CreateComponent[Transform]()
var scaleComponent = CreateComponent[Scale]()

type TransformSystem struct {
	n int
}

func (s *TransformSystem) Init()    {}
func (s *TransformSystem) Destroy() {}
func (s *TransformSystem) Run(world *ECS) {
	s.n++
	transformComponent.Each(world, func(entity *Entity, data *Transform) {
		data.X += 1
		data.Y -= 1
		data.Z += 2
	})
}

type ScaleSystem struct {
	n int
}

func (s *ScaleSystem) Init()    {}
func (s *ScaleSystem) Destroy() {}
func (s *ScaleSystem) Run(world *ECS) {
	s.n++
	scaleComponent.Each(world, func(entity *Entity, data *Scale) {
		tr := transformComponent.Get(entity)
		if tr == nil {
			return
		}

		data.Value += 0.1
	})
}
