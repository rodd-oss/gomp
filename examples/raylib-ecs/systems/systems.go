/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp/examples/raylib-ecs/assets"
	ecs2 "gomp/pkgs/ecs"
)

var InitService = ecs2.CreateSystemService(&initController{windowWidth: 800, windowHeight: 600})
var DebugService = ecs2.CreateSystemService(&debugController{})

var ExampleService = ecs2.CreateSystemService(&exampleController{})

var PlayerService = ecs2.CreateSystemService(&playerController{})
var SpawnService = ecs2.CreateSystemService(&spawnController{})
var HpService = ecs2.CreateSystemService(&hpController{}, &PlayerService)
var ColorService = ecs2.CreateSystemService(&colorController{}, &HpService)

var AnimationSpriteMatrixService = ecs2.CreateSystemService(&animationSpriteMatrixController{})
var AnimationPlayerService = ecs2.CreateSystemService(&animationPlayerController{}, &AnimationSpriteMatrixService)

// Texture Render Main Systems
var TRSpriteService = ecs2.CreateSystemService(&trSpriteController{})
var TRSpriteSheetService = ecs2.CreateSystemService(&trSpriteSheetController{})
var TRSpriteMatrixService = ecs2.CreateSystemService(&trSpriteMatrixController{})

// Texture Render Sub Systems
var TRAnimationService = ecs2.CreateSystemService(&trAnimationController{})
var TRMirroredService = ecs2.CreateSystemService(&trMirroredController{}, &TRAnimationService)
var TRPositionService = ecs2.CreateSystemService(&trPositionController{})
var TRRotationService = ecs2.CreateSystemService(&trRotationController{})
var TRScaleService = ecs2.CreateSystemService(&trScaleController{})
var TRTintService = ecs2.CreateSystemService(&trTintController{})

// Render System
var AssetLibService = ecs2.CreateSystemService(&assetLibController{
	assets: []ecs2.AnyAssetLibrary{assets.Textures},
})
var RenderService = ecs2.CreateSystemService(&renderController{})
