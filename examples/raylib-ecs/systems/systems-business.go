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

package systems

import "gomp/pkg/ecs"

// Business logic systems

var PlayerService = ecs.CreateSystemService(&playerController{})
var SpawnService = ecs.CreateSystemService(&spawnController{})
var HpService = ecs.CreateSystemService(&hpController{}, &PlayerService)
var ColorService = ecs.CreateSystemService(&colorController{}, &HpService)
