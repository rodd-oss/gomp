/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"gomp"
)

func NewAssetLibSystem(assets []gomp.AnyAssetLibrary) AssetLibSystem {
	return AssetLibSystem{
		assets: assets,
	}
}

type AssetLibSystem struct {
	assets []gomp.AnyAssetLibrary
}

func (s *AssetLibSystem) Init() {}
func (s *AssetLibSystem) Run() {
	for _, asset := range s.assets {
		asset.LoadAll()
	}
}
func (s *AssetLibSystem) Destroy() {
	for _, asset := range s.assets {
		asset.UnloadAll()
	}
}
