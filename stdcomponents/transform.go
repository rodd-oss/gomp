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

type Transform2d struct {
	Position vectors.Vec2
	Rotation vectors.Radians
	Scale    vectors.Vec2
}

type TransformComponentManager = ecs.ComponentManager[Transform2d]

func NewTransformComponentManager() TransformComponentManager {
	return ecs.NewComponentManager[Transform2d](TransformComponentId)
}
