/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- Монтажер сука Donated 50 RUB

Thank you for your support!
*/

package main

import (
	"gomp/examples/raylib-ecs/components"
	"gomp/examples/raylib-ecs/systems"
	"gomp/pkg/ecs"
)

func main() {
	world := ecs.CreateWorld("main-raylib")
	defer world.Destroy()

	// world.ApplyPatch(patch)

	// components.TransformService.OnPatch(func(newstate components.Transform, oldstate components.Transform) components.Transform {
	// 	return new
	// })

	world.RegisterComponentServices(
		&components.PositionService,
		&components.RotationService,
		&components.ScaleService,
		&components.MirroredService,
		&components.HealthService,
		&components.SpriteService,
		&components.SpriteSheetService,
		&components.SpriteMatrixService,
		&components.TintService,
		&components.AnimationPlayerService,
		&components.AnimationStateService,
		&components.TextureRenderService,
	)

	world.RegisterSystems().
		Sequential( // Initial systems: Main thread
			&systems.InitService,
		).
		Parallel( // Network receive systems
			&systems.NetworkService,
			&systems.NetworkReceiveService,
		).
		Parallel( // Business logic systems
			&systems.PlayerService,
			&systems.HpService,
			&systems.ColorService,
		).
		Parallel(
			// Network send systems
			&systems.NetworkSendService,
			// Prerender systems
			&systems.AnimationSpriteMatrixService,
			&systems.AnimationPlayerService,
			&systems.TRSpriteService,
			&systems.TRSpriteSheetService,
			&systems.TRSpriteMatrixService,
			&systems.TRAnimationService,
			&systems.TRMirroredService,
			&systems.TRPositionService,
			&systems.TRRotationService,
			&systems.TRScaleService,
			&systems.TRTintService,
		).
		Sequential( // Render systems: Main thread
			&systems.RenderService,
			&systems.DebugService,
			&systems.AssetLibService,
		)

	world.Run(1)
}
