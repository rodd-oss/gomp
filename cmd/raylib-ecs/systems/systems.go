/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/cmd/raylib-ecs/assets"
	"gomp_game/pkgs/gomp/ecs"
)

var InitService = ecs.CreateSystemService(&initController{windowWidth: 800, windowHeight: 600})

var ExampleService = ecs.CreateSystemService(&exampleController{})

var PlayerService = ecs.CreateSystemService(&playerController{})
var SpawnService = ecs.CreateSystemService(&spawnController{})
var HpService = ecs.CreateSystemService(&hpController{}, &PlayerService)
var ColorService = ecs.CreateSystemService(&colorController{}, &HpService)

var AnimationService = ecs.CreateSystemService(&animationController{})

// Texture Render Systems
var TRSpriteService = ecs.CreateSystemService(&trSpriteController{})
var TRSpriteSheetService = ecs.CreateSystemService(&trSpriteSheetController{})
var TRAnimationService = ecs.CreateSystemService(&trAnimationController{}, &TRSpriteService, &TRSpriteSheetService)
var TRPositionService = ecs.CreateSystemService(&trPositionController{}, &TRSpriteService, &TRSpriteSheetService)
var TRRotationService = ecs.CreateSystemService(&trRotationController{}, &TRSpriteService, &TRSpriteSheetService)
var TRScaleService = ecs.CreateSystemService(&trScaleController{}, &TRSpriteService, &TRSpriteSheetService)

// Render System
var AssetLibService = ecs.CreateSystemService(&assetLibController{
	assets: []ecs.AnyAssetLibrary{assets.Textures},
})
var RenderService = ecs.CreateSystemService(&renderController{})
