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

package stdentities

import (
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
)

type CreateCollisionGridManagers struct {
	EntityManager *ecs.EntityManager
	Grid          *stdcomponents.CollisionGridComponentManager
}

func CreateCollisionGrid(
	props *CreateCollisionGridManagers,
	layer stdcomponents.CollisionLayer,
	cellSize float32,
) ecs.Entity {
	e := props.EntityManager.Create()

	props.Grid.Create(e, stdcomponents.NewCollisionGrid(layer, cellSize))

	return e
}
