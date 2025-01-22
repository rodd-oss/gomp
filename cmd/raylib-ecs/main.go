/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file deveopment:
-===-===-===-===-===-===-===-===-===-===

<- Монтажер сука Donated 50 RUB

Thank you for your support!
*/

package main

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/cmd/raylib-ecs/systems"
	"gomp_game/pkgs/gomp/ecs"
)

func main() {
	world := ecs.CreateWorld("main-raylib")
	defer world.Destroy()

	// patch := world.GetPatch()
	// world.ApplyPatch(patch)

	// components.TransformService.OnPatch(func(newstate components.Transform, oldstate components.Transform) components.Transform {
	// 	return new
	// })

	world.RegisterComponentServices(
		&components.PositionService,
		&components.RotationService,
		&components.ScaleService,
		&components.HealthService,
		&components.SpriteService,
		&components.SpriteSheetService,
		&components.AnimationService,
		&components.TextureRenderService,
	)

	world.RegisterSystems().
		Sequential(
			&systems.SpawnService,
		).
		Parallel(
			&systems.HpService,
			&systems.ColorService,
		).
		Parallel(
			&systems.TRSpriteService,
			&systems.TRSpriteSheetService,
			&systems.TRAnimationService,
			&systems.TRPositionService,
			&systems.TRRotationService,
			&systems.TRScaleService,
		).
		Sequential(
			&systems.RenderService,
		)

	world.Run(60)
}
