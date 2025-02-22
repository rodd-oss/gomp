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

package gomp

import (
	"github.com/negrel/assert"
	"reflect"
	"time"
)

type AnyGame interface {
	Init()
	Update(dt time.Duration)
	FixedUpdate(dt time.Duration)
	Render(interpolation float32)
	Destroy()
	ShouldDestroy() bool
}

func NewGame(scenes ...AnyScene) Game {
	sceneSet := make(map[SceneId]AnyScene, len(scenes))

	for i := range len(scenes) {
		id := scenes[i].Id()

		_, exists := sceneSet[id]
		assert.False(exists, "Scene with id %d already exists. Duplicating ids?", id)

		sceneSet[id] = scenes[i]
	}

	game := Game{
		Scenes: sceneSet,
	}

	return game
}

type Game struct {
	Scenes         map[SceneId]AnyScene
	CurrentSceneId SceneId

	shouldDestroy bool
}

func (g *Game) Init() {
	for _, scene := range g.Scenes {
		g.injectToScene(scene)
	}

	scene, ok := g.Scenes[g.CurrentSceneId]
	assert.True(ok, "Scene not found")
	scene.Init()
}

func (g *Game) Update(dt time.Duration) {
	// Scenes
	scene, ok := g.Scenes[g.CurrentSceneId]
	assert.True(ok, "Scene not found")
	g.CurrentSceneId = scene.Update(dt)
}

func (g *Game) FixedUpdate(dt time.Duration) {
	// Scenes
	scene, ok := g.Scenes[g.CurrentSceneId]
	assert.True(ok, "Scene not found")
	scene.FixedUpdate(dt)
}

func (g *Game) Render(interpolation float32) {
	scene, ok := g.Scenes[g.CurrentSceneId]
	assert.True(ok, "Scene not found")
	scene.Render(interpolation)
}

func (g *Game) Destroy() {
	scene, ok := g.Scenes[g.CurrentSceneId]
	assert.True(ok, "Scene not found")
	scene.Destroy()
}

func (g *Game) ShouldDestroy() bool {
	return g.shouldDestroy
}

func (g *Game) SetShouldDestroy(value bool) {
	g.shouldDestroy = value
}

func (g *Game) injectToScene(scene AnyScene) {
	reflectedScene := reflect.ValueOf(scene).Elem()
	sceneLen := reflectedScene.NumField()

	reflectedGame := reflect.ValueOf(g)
	gameType := reflect.TypeOf(g)

	for i := 0; i < sceneLen; i++ {
		field := reflectedScene.Field(i)
		fieldType := field.Type()

		if fieldType == gameType {
			reflectedScene.Field(i).Set(reflectedGame)
			continue
		}
	}
}
