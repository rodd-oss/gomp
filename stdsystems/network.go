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

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/network"
	"time"
)

type NetworkMode int

const (
	None NetworkMode = iota
	Server
	Client
)

func NewNetworkSystem() *NetworkSystem {
	return &NetworkSystem{}
}

type NetworkSystem struct {
}

func (s *NetworkSystem) Init() {
}
func (s *NetworkSystem) Run(dt time.Duration) {
	if rl.IsKeyPressed(rl.KeyP) {
		network.Quic.Host("127.0.0.1:27015")
	}

	if rl.IsKeyPressed(rl.KeyO) {
		network.Quic.Connect("127.0.0.1:27015")
	}
}
func (s *NetworkSystem) Destroy() {}
