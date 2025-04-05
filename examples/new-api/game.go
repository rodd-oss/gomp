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

package main

import (
	"errors"
	"fmt"
	"gomp"
	"gomp/examples/new-api/instances"
	"gomp/pkg/ecs"
	"gomp/stdsystems"
	"time"
)

func NewGame(initialScene gomp.AnyScene) Game {
	game := Game{
		worlds:            make([]instances.World, 0),
		scenes:            make([]gomp.AnyScene, 0),
		lookup:            make(map[gomp.SceneId]int),
		renderSystem:      stdsystems.NewRenderSystem(),
		minimapSystem:     stdsystems.NewMinimapRenderSystem(),
		finalRenderSystem: stdsystems.NewRenderTexture2DRenderSystem(),
	}

	game.LoadScene(initialScene)
	game.SetActiveScene(initialScene.Id())

	return game
}

type Game struct {
	// Scenes
	worlds         []instances.World
	scenes         []gomp.AnyScene
	lookup         map[gomp.SceneId]int
	currentSceneId gomp.SceneId

	// Systems
	renderSystem  stdsystems.RenderSystem
	minimapSystem stdsystems.MinimapRenderSystem

	// Utils
	shouldClose       bool
	finalRenderSystem stdsystems.RenderTexture2DRenderSystem
}

func (g *Game) Init() {
	err := g.injectWorldToRender()
	if err != nil {
		panic(fmt.Sprint("injectWorldToRender failed: ", err.Error()))
	}

	for i := range g.scenes {
		world := &g.worlds[i]
		world.Init()

		world.Systems.ColliderSystem.Init()
		world.Systems.Velocity.Init()
		world.Systems.CollisionDetectionBVH.Init()
		world.Systems.CollisionResolution.Init()
		world.Systems.AnimationSpriteMatrix.Init()
		world.Systems.AnimationPlayer.Init()
		world.Systems.SpriteMatrix.Init()
		world.Systems.Sprite.Init()
		world.Systems.YSort.Init()
		world.Systems.Debug.Init()
		world.Systems.AssetLib.Init()
		world.Systems.Audio.Init()
		world.Systems.SpatialAudio.Init()

		g.scenes[i].Init(world)
	}
	g.renderSystem.Init()
	g.minimapSystem.Init()
	g.finalRenderSystem.Init()
}

func (g *Game) Update(dt time.Duration) {
	for i := range g.scenes {
		g.scenes[i].Update(dt)
	}
}

func (g *Game) FixedUpdate(dt time.Duration) {
	for i := range g.worlds {
		world := &g.worlds[i]

		world.Systems.ColliderSystem.Run(dt)
		world.Systems.Velocity.Run(dt)
		world.Systems.CollisionDetectionBVH.Run(dt)
		world.Systems.CollisionResolution.Run(dt)

		g.scenes[i].FixedUpdate(dt)
	}
}

func (g *Game) Render(dt time.Duration) {

	id := g.lookup[g.currentSceneId]
	world := &g.worlds[id]

	world.Systems.AnimationSpriteMatrix.Run()
	world.Systems.AnimationPlayer.Run()
	world.Systems.SpriteMatrix.Run()
	world.Systems.Sprite.Run()
	world.Systems.Debug.Run()
	world.Systems.AssetLib.Run()
	world.Systems.YSort.Run()
	world.Systems.Audio.Run(dt)
	world.Systems.SpatialAudio.Run(dt)

	scene := g.scenes[g.lookup[g.currentSceneId]]
	scene.Render(dt)

	g.shouldClose = g.renderSystem.Run(dt)
	g.minimapSystem.Run(dt)
	g.finalRenderSystem.Run(dt)
}

func (g *Game) Destroy() {
	err := g.injectWorldToRender()
	if err != nil {
		panic("jfdk")
	}

	for i := range g.scenes {
		g.scenes[i].Destroy()

		world := &g.worlds[i]
		world.Systems.Velocity.Destroy()
		world.Systems.CollisionDetectionBVH.Destroy()
		world.Systems.CollisionResolution.Destroy()
		world.Systems.AnimationSpriteMatrix.Destroy()
		world.Systems.AnimationPlayer.Destroy()
		world.Systems.SpriteMatrix.Destroy()
		world.Systems.Sprite.Destroy()
		world.Systems.YSort.Destroy()
		world.Systems.Debug.Destroy()
		world.Systems.AssetLib.Destroy()
		world.Systems.Audio.Destroy()
		world.Systems.SpatialAudio.Destroy()

		world.Destroy()
	}

	g.finalRenderSystem.Destroy()
	g.minimapSystem.Destroy()
	g.renderSystem.Destroy()
}

func (g *Game) ShouldDestroy() bool {
	return g.shouldClose
}

func (g *Game) LoadScene(scene gomp.AnyScene) {
	g.scenes = append(g.scenes, scene)
	g.worlds = append(g.worlds, ecs.NewWorld(instances.NewComponentList(), instances.NewSystemList()))
	g.lookup[scene.Id()] = len(g.scenes) - 1
}

func (g *Game) SetActiveScene(id gomp.SceneId) {
	g.currentSceneId = id
}

func (g *Game) injectWorldToRender() error {
	id, exists := g.lookup[g.currentSceneId]
	if !exists {
		return errors.New("scene not found")
	}

	world := &g.worlds[id]

	g.renderSystem.InjectWorld(
		&stdsystems.RenderInjector{
			EntityManager:                      &world.Entities,
			RlTexturePros:                      &world.Components.RLTexturePro,
			Positions:                          &world.Components.Position,
			Rotations:                          &world.Components.Rotation,
			Scales:                             &world.Components.Scale,
			AnimationPlayers:                   &world.Components.AnimationPlayer,
			Tints:                              &world.Components.Tint,
			Flips:                              &world.Components.Flip,
			Renderables:                        &world.Components.Renderable,
			AnimationStates:                    &world.Components.AnimationState,
			Sprites:                            &world.Components.Sprite,
			SpriteMatrixes:                     &world.Components.SpriteMatrix,
			RenderOrders:                       &world.Components.RenderOrder,
			BoxColliders:                       &world.Components.ColliderBox,
			CircleColliders:                    &world.Components.ColliderCircle,
			AABBs:                              &world.Components.AABB,
			Collisions:                         &world.Components.Collision,
			ColliderSleepStateComponentManager: &world.Components.ColliderSleepState,
			BvhTrees:                           &world.Components.BvhTree,
			MainCameras:                        &world.Components.MainCamera,
			RenderTexture2D:                    &world.Components.RenderTexture2D,
		})
	g.minimapSystem.InjectWorld(
		&stdsystems.MinimapRenderInjector{
			EntityManager:                      &world.Entities,
			RlTexturePros:                      &world.Components.RLTexturePro,
			Positions:                          &world.Components.Position,
			Rotations:                          &world.Components.Rotation,
			Scales:                             &world.Components.Scale,
			AnimationPlayers:                   &world.Components.AnimationPlayer,
			Tints:                              &world.Components.Tint,
			Flips:                              &world.Components.Flip,
			Renderables:                        &world.Components.Renderable,
			AnimationStates:                    &world.Components.AnimationState,
			Sprites:                            &world.Components.Sprite,
			SpriteMatrixes:                     &world.Components.SpriteMatrix,
			RenderOrders:                       &world.Components.RenderOrder,
			BoxColliders:                       &world.Components.ColliderBox,
			CircleColliders:                    &world.Components.ColliderCircle,
			AABBs:                              &world.Components.AABB,
			Collisions:                         &world.Components.Collision,
			ColliderSleepStateComponentManager: &world.Components.ColliderSleepState,
			BvhTrees:                           &world.Components.BvhTree,
			MinimapCameras:                     &world.Components.MinimapCamera,
			RenderTexture2D:                    &world.Components.RenderTexture2D,
		})
	g.finalRenderSystem.InjectWorld(
		&stdsystems.RenderTexture2DRenderInjector{
			EntityManager:   &world.Entities,
			RenderTexture2D: &world.Components.RenderTexture2D,
		})

	return nil
}
