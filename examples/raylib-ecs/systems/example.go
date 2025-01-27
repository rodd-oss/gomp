/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp/pkgs/ecs"
)

type exampleController struct{}

func (s *exampleController) Init(world *ecs.World)        {}
func (s *exampleController) Update(world *ecs.World)      {}
func (s *exampleController) FixedUpdate(world *ecs.World) {}
func (s *exampleController) Destroy(world *ecs.World)     {}
