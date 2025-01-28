/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

type System struct {
	Controller SystemController
	ID         uint16
}

type SystemController interface {
	Init(game *Game)
	Update(dt float64)
}

func (s *System) Init(game *Game) {
	s.Controller.Init(game)
}

func (s *System) Update(dt float64) {
	s.Controller.Update(dt)
}
