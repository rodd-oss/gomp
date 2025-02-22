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

import "gomp/pkg/ecs"

type ViewPosition struct {
	LastPositionX, LastPositionY       float32
	CurrentPositionX, CurrentPositionY float32
}

type ViewPositionComponentManager = ecs.ComponentManager[ViewPosition]

func NewViewPositionComponentManager() ViewPositionComponentManager {
	return ecs.NewComponentManager[ViewPosition](ViewPositionId)
}
