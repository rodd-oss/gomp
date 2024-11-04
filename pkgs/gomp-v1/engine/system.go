/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

type SystemController interface {
	Init(scene *Scene)
	Update(dt float64)
}

type System struct {
	controller SystemController
}

func (s *System) Init(scene *Scene) {
	s.controller.Init(scene)
}

func (s *System) Update(dt float64) {
	s.controller.Update(dt)
}

func CreateSystem(controller SystemController) System {
	return System{
		controller: controller,
	}
}
