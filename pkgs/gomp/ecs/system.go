/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"github.com/yohamta/donburi"
)

type System struct {
	Controller SystemController
}

type SystemController interface {
	Init(world donburi.World)
	Update(dt float64)
}

func (s *System) Init(world donburi.World) {
	s.Controller.Init(world)
}

func (s *System) Update(dt float64) {
	s.Controller.Update(dt)
}
