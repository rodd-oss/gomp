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
var DebugService = ecs.CreateSystemService(&debugController{})

var ExampleService = ecs.CreateSystemService(&exampleController{})

var PlayerService = ecs.CreateSystemService(&playerController{})
var SpawnService = ecs.CreateSystemService(&spawnController{})
var HpService = ecs.CreateSystemService(&hpController{}, &PlayerService)
var ColorService = ecs.CreateSystemService(&colorController{}, &HpService)

var AnimationSpriteMatrixService = ecs.CreateSystemService(&animationSpriteMatrixController{})
var AnimationPlayerService = ecs.CreateSystemService(&animationPlayerController{}, &AnimationSpriteMatrixService)

// Texture Render Main Systems
var trSpriteSystemServices = []ecs.AnySystemServicePtr{&TRSpriteService, &TRSpriteSheetService, &TRSpriteMatrixService}

var TRSpriteService = ecs.CreateSystemService(&trSpriteController{})
var TRSpriteSheetService = ecs.CreateSystemService(&trSpriteSheetController{})
var TRSpriteMatrixService = ecs.CreateSystemService(&trSpriteMatrixController{})

// Texture Render Sub Systems
var TRAnimationService = ecs.CreateSystemService(&trAnimationController{})
var TRMirroredService = ecs.CreateSystemService(&trMirroredController{}, &TRAnimationService)
var TRPositionService = ecs.CreateSystemService(&trPositionController{})
var TRRotationService = ecs.CreateSystemService(&trRotationController{})
var TRScaleService = ecs.CreateSystemService(&trScaleController{})
var TRTintService = ecs.CreateSystemService(&trTintController{})

// Render System
var AssetLibService = ecs.CreateSystemService(&assetLibController{
	assets: []ecs.AnyAssetLibrary{assets.Textures},
})
var RenderService = ecs.CreateSystemService(&renderController{})
