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

package stdcomponents

import (
	"gomp/pkg/ecs"
)

type AnimationState int32

type AnimationStateComponentManager = ecs.ComponentManager[AnimationState]

func NewAnimationStateComponentManager() AnimationStateComponentManager {
	return ecs.NewComponentManager[AnimationState](AnimationStateComponentId)
}
