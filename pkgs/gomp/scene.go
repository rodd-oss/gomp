/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"gomp_game/pkgs/gomp/ecs"

	"github.com/yohamta/donburi"
)

type Scene struct {
	Name string

	Systems        []*ecs.System
	Entities       []ecs.Entity
	SceneComponent *donburi.ComponentType[SceneData]
}

type SceneData struct {
	Name string
}

var SceneComponent = CreateComponent[SceneData]

type sceneFactoryEntities struct {
	scene *Scene
}

func (f sceneFactoryEntities) AddEntities(ent ...ecs.Entity) sceneFactorySystems {
	f.scene.Entities = ent
	return sceneFactorySystems(f)
}

type sceneFactorySystems struct {
	scene *Scene
}

func (f sceneFactorySystems) AddSystems(sys ...ecs.System) Scene {
	lenSys := len(sys)
	for i := 0; i < lenSys; i++ {
		f.scene.Systems = append(f.scene.Systems, &sys[i])
	}
	return *f.scene
}

func CreateScene(name string) sceneFactoryEntities {
	scene := new(Scene)
	scene.SceneComponent = SceneComponent(SceneData{
		Name: name,
	})

	scene.Name = name

	factory := sceneFactoryEntities{
		scene: scene,
	}

	return factory
}
