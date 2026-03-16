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

package stdcomponents

import (
	"gomp/pkg/ecs"
	"gomp/vectors"
)

type AABB struct {
	Min vectors.Vec2
	Max vectors.Vec2
}

func (a AABB) Center() vectors.Vec2 {
	return a.Min.Add(a.Max).Scale(0.5)
}

func (a AABB) Rect() vectors.Rectangle {
	return vectors.Rectangle{
		X:      a.Min.X,
		Y:      a.Min.Y,
		Width:  a.Max.X - a.Min.X,
		Height: a.Max.Y - a.Min.Y,
	}
}

type AABBComponentManager = ecs.ComponentManager[AABB]

func NewAABBComponentManager() AABBComponentManager {
	return ecs.NewComponentManager[AABB](AABBComponentId)
}
