/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package entities

import (
	"gomp_game/pkgs/example/resources"
	"gomp_game/pkgs/gomp"

	"github.com/jakecoffman/cp/v2"
)

type HealthData struct {
	Health    int
	MaxHealth int
}

var HealthComponent = gomp.CreateComponent[HealthData]()
var ManaComponent = gomp.CreateComponent[uint16]()

var Player = gomp.CreateEntity(
	gomp.SpriteComponent.New(resources.Sprites["player.png"]),
	gomp.BodyComponent.New(*cp.NewKinematicBody()),

	HealthComponent.New(HealthData{Health: 100, MaxHealth: 100}),
	ManaComponent.New(125),
)
