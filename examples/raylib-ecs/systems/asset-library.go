/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	ecs2 "gomp/pkgs/ecs"
)

type assetLibController struct {
	assets []ecs2.AnyAssetLibrary
}

func (s *assetLibController) Init(world *ecs2.World) {}
func (s *assetLibController) Update(world *ecs2.World) {
	for _, asset := range s.assets {
		asset.LoadAll()
	}
}
func (s *assetLibController) FixedUpdate(world *ecs2.World) {}
func (s *assetLibController) Destroy(world *ecs2.World) {
	for _, asset := range s.assets {
		asset.UnloadAll()
	}
}
