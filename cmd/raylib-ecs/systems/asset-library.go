/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/pkgs/gomp/ecs"
)

type assetLibController struct {
	assets []ecs.AnyAssetLibrary
}

func (s *assetLibController) Init(world *ecs.World) {}
func (s *assetLibController) Update(world *ecs.World) {
	for _, asset := range s.assets {
		asset.LoadAll()
	}
}
func (s *assetLibController) FixedUpdate(world *ecs.World) {}
func (s *assetLibController) Destroy(world *ecs.World) {
	for _, asset := range s.assets {
		asset.UnloadAll()
	}
}
