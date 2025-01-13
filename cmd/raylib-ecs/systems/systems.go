/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import "gomp_game/pkgs/gomp/ecs"

var Spawn = ecs.CreateSystem(&systemSpawn{})
var CalcHp = ecs.CreateSystem(&systemCalcHp{}, &Spawn)
var CalcColor = ecs.CreateSystem(&systemCalcColor{}, &CalcHp)
var Render = ecs.CreateSystem(&systemRender{width: 800, height: 600}, &CalcColor, &CalcHp)
