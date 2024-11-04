/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"gomp_game/pkgs/gomp/ecs"
	"gomp_game/pkgs/gomp/scene"
)

type sceneFactoryEntities struct {
	scene *scene.Scene
}

func (f sceneFactoryEntities) AddEntities(ent ...ecs.Entity) sceneFactorySystems {
	f.scene.Entities = ent
	return sceneFactorySystems(f)
}

type sceneFactorySystems struct {
	scene *scene.Scene
}

func (f sceneFactorySystems) AddSystems(sys ...ecs.System) scene.Scene {
	f.scene.Systems = sys
	return *f.scene
}

func CreateScene(name string) sceneFactoryEntities {
	scene := new(scene.Scene)

	scene.Name = name

	factory := sceneFactoryEntities{
		scene: scene,
	}

	return factory
}
