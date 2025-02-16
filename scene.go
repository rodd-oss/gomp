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

package gomp

import "time"

type Scene struct {
	IsInitialized bool
	Enabled       bool
	Systems       AnyEngineSystemSet
}

func (s *Scene) Init() {
	s.Systems.Init()
	s.IsInitialized = true
}
func (s *Scene) Update(dt time.Duration) {
	s.Systems.Update(dt)
}
func (s *Scene) FixedUpdate(dt time.Duration) {
	s.Systems.FixedUpdate(dt)
}
func (s *Scene) Destroy() {
	s.Systems.Destroy()
	s.IsInitialized = false
}
