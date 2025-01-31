/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp/pkg/ecs"
)

type initController struct {
}

func (s *initController) Init(world *ecs.World) {}

func (s *initController) Update(world *ecs.World)      {}
func (s *initController) FixedUpdate(world *ecs.World) {}
func (s *initController) Destroy(world *ecs.World)     {}
