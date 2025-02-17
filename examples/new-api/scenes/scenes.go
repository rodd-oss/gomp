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

package scenes

import (
	"gomp"
	"gomp/examples/new-api/components"
	"gomp/pkg/ecs"
)

func NewSceneSet(world *ecs.World, components *components.GameComponents) map[gomp.SceneId]gomp.AnyScene {
	sceneSet := make(map[gomp.SceneId]gomp.AnyScene)

	sceneSet[MainSceneId] = NewMainScene(world, components)

	return sceneSet
}
