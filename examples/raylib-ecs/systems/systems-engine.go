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

package systems

import (
	"gomp"
	"gomp/examples/raylib-ecs/assets"
	"gomp/pkg/ecs"
)

// Engine Initial Systems

var InitService = ecs.CreateSystemService(&initController{})
var DebugService = ecs.CreateSystemService(&debugController{})

// Engine Network systems

var NetworkService = ecs.CreateSystemService(&networkController{})
var NetworkSendService = ecs.CreateSystemService(&networkSendController{})
var NetworkReceiveService = ecs.CreateSystemService(&networkReceiveController{})

// Engine Texture PreRender systems

var AnimationSpriteMatrixService = ecs.CreateSystemService(&animationSpriteMatrixController{})
var AnimationPlayerService = ecs.CreateSystemService(&animationPlayerController{}, &AnimationSpriteMatrixService)
var TRSpriteService = ecs.CreateSystemService(&trSpriteController{}, &AnimationPlayerService)
var TRSpriteSheetService = ecs.CreateSystemService(&trSpriteSheetController{}, &AnimationPlayerService)
var TRSpriteMatrixService = ecs.CreateSystemService(&trSpriteMatrixController{}, &AnimationPlayerService)
var TRAnimationService = ecs.CreateSystemService(&trAnimationController{}, &TRSpriteService, &TRSpriteSheetService, &TRSpriteMatrixService)
var TRMirroredService = ecs.CreateSystemService(&trMirroredController{}, &TRSpriteService, &TRSpriteSheetService, &TRSpriteMatrixService, &TRAnimationService)
var TRPositionService = ecs.CreateSystemService(&trPositionController{}, &TRSpriteService, &TRSpriteSheetService, &TRSpriteMatrixService)
var TRRotationService = ecs.CreateSystemService(&trRotationController{}, &TRSpriteService, &TRSpriteSheetService, &TRSpriteMatrixService)
var TRScaleService = ecs.CreateSystemService(&trScaleController{}, &TRSpriteService, &TRSpriteSheetService, &TRSpriteMatrixService)
var TRTintService = ecs.CreateSystemService(&trTintController{}, &TRSpriteService, &TRSpriteSheetService, &TRSpriteMatrixService)

// Engine Render Systems

var AssetLibService = ecs.CreateSystemService(&assetLibController{
	assets: []gomp.AnyAssetLibrary{assets.Textures},
})
var RenderService = ecs.CreateSystemService(&renderController{
	windowWidth:  1024,
	windowHeight: 768,
})
