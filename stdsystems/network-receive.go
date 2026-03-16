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

package stdsystems

import "time"

func NewNetworkReceiveSystem() NetworkReceiveSystem {
	return NetworkReceiveSystem{}
}

type NetworkReceiveSystem struct{}

func (s *NetworkReceiveSystem) Init()                {}
func (s *NetworkReceiveSystem) Run(dt time.Duration) {}
func (s *NetworkReceiveSystem) Destroy()             {}
