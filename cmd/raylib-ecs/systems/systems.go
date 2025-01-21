/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import "gomp_game/pkgs/gomp/ecs"

var SpawnService = ecs.CreateSystemService(&spawnController{})
var HpService = ecs.CreateSystemService(&hpController{}, &SpawnService)
var ColorService = ecs.CreateSystemService(&colorController{}, &HpService)
var SpriteService = ecs.CreateSystemService(&spriteController{}, &ColorService)
var RenderService = ecs.CreateSystemService(&renderController{width: 800, height: 600})
var ExampleService = ecs.CreateSystemService(&exampleController{})
