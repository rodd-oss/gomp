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

package systems

import (
	"gomp/pkg/ecs"
)

type networkReceiveController struct{}

func (s *networkReceiveController) Init(world *ecs.EntityManager)        {}
func (s *networkReceiveController) Update(world *ecs.EntityManager)      {}
func (s *networkReceiveController) FixedUpdate(world *ecs.EntityManager) {}
func (s *networkReceiveController) Destroy(world *ecs.EntityManager)     {}
