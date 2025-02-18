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

package instances

import (
	"gomp/examples/new-api/assets"
	"gomp/examples/new-api/systems"
	"gomp/pkg/ecs"
	"gomp/stdsystems"
	"reflect"
)

func NewSystemList(world *ecs.World, components *ComponentList) SystemList {
	newSystemList := SystemList{
		Player: systems.NewPlayerSystem(),
		Debug:  stdsystems.NewDebugSystem(),

		Velocity: stdsystems.NewVelocitySystem(),

		Network:        stdsystems.NewNetworkSystem(),
		NetworkReceive: stdsystems.NewNetworkReceiveSystem(),
		NetworkSend:    stdsystems.NewNetworkSendSystem(),

		AnimationSpriteMatrix: stdsystems.NewAnimationSpriteMatrixSystem(),
		AnimationPlayer:       stdsystems.NewAnimationPlayerSystem(),

		TextureRenderSprite:      stdsystems.NewTextureRenderSpriteSystem(),
		TextureRenderSpriteSheet: stdsystems.NewTextureRenderSpriteSheetSystem(),
		TextureRenderMatrix:      stdsystems.NewTextureRenderMatrixSystem(),

		TextureRenderAnimation: stdsystems.NewTextureRenderAnimationSystem(),
		TextureRenderFlip:      stdsystems.NewTextureRenderFlipSystem(),
		TextureRenderPosition:  stdsystems.NewTextureRenderPositionSystem(),
		TextureRenderRotation:  stdsystems.NewTextureRenderRotationSystem(),
		TextureRenderScale:     stdsystems.NewTextureRenderScaleSystem(),
		TextureRenderTint:      stdsystems.NewTextureRenderTintSystem(),

		AssetLib: stdsystems.NewAssetLibSystem([]ecs.AnyAssetLibrary{assets.Textures}),
		Render:   stdsystems.NewRenderSystem(),
	}

	InjectECSToSystems(&newSystemList, world, components)

	return newSystemList
}

type SystemList struct {
	Player systems.PlayerSystem
	Debug  stdsystems.DebugSystem

	Velocity stdsystems.VelocitySystem

	// Network
	Network        stdsystems.NetworkSystem
	NetworkReceive stdsystems.NetworkReceiveSystem
	NetworkSend    stdsystems.NetworkSendSystem
	// Animation
	AnimationSpriteMatrix stdsystems.AnimationSpriteMatrixSystem
	AnimationPlayer       stdsystems.AnimationPlayerSystem
	// Prerender init
	TextureRenderSprite      stdsystems.TextureRenderSpriteSystem
	TextureRenderSpriteSheet stdsystems.TextureRenderSpriteSheetSystem
	TextureRenderMatrix      stdsystems.TextureRenderMatrixSystem
	// Prerender fill
	TextureRenderAnimation stdsystems.TextureRenderAnimationSystem
	TextureRenderFlip      stdsystems.TextureRenderFlipSystem
	TextureRenderPosition  stdsystems.TextureRenderPositionSystem
	TextureRenderRotation  stdsystems.TextureRenderRotationSystem
	TextureRenderScale     stdsystems.TextureRenderScaleSystem
	TextureRenderTint      stdsystems.TextureRenderTintSystem
	// Render
	AssetLib stdsystems.AssetLibSystem
	Render   stdsystems.RenderSystem
}

type AnySystemList interface{}
type AnyComponentList interface{}

func InjectECSToSystems(systemList AnySystemList, world *ecs.World, componentList AnyComponentList) {
	reflectedSystemList := reflect.ValueOf(systemList).Elem()
	systemsLen := reflectedSystemList.NumField()

	reflectedComponentList := reflect.ValueOf(componentList).Elem()
	componentsLen := reflectedComponentList.NumField()

	worldType := reflect.TypeOf(world)

	for i := range systemsLen {
		system := reflectedSystemList.Field(i)
		systemLen := system.NumField()

		for j := range systemLen {
			systemField := system.Field(j)
			systemFieldType := systemField.Type()

			if systemFieldType == worldType {
				system.Field(j).Set(reflect.ValueOf(world))
				continue
			}

			shouldEscape := false
			for k := range componentsLen {
				component := reflectedComponentList.Field(k).Elem()
				componentType := component.Type()

				if systemFieldType.Kind() == reflect.Ptr {
					if systemFieldType.Elem() == componentType {
						system.Field(j).Set(reflectedComponentList.Field(k))
						shouldEscape = true
					}
				} else {
					if systemFieldType == componentType {
						system.Field(j).Set(reflectedComponentList.Field(k))
						shouldEscape = true
					}
				}

				if shouldEscape {
					break
				}
			}
		}
	}
}
