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

package components

import (
	"gomp/pkg/ecs"
	"time"
)

type Weapon struct {
	Damage       int
	Cooldown     time.Duration
	CooldownLeft time.Duration
}

type WeaponComponentManager = ecs.ComponentManager[Weapon]

func NewWeaponComponentManager() WeaponComponentManager {
	return ecs.NewComponentManager[Weapon](WeaponComponentId)
}
