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

type networkSendController struct{}

func (s *networkSendController) Init(world *ecs.World)        {}
func (s *networkSendController) Update(world *ecs.World)      {}
func (s *networkSendController) FixedUpdate(world *ecs.World) {}
func (s *networkSendController) Destroy(world *ecs.World)     {}
