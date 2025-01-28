/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

type Scene struct {
	Name string

	Systems        []*System
	Entities       []Entity
	SceneComponent ComponentFactory[SceneData]
}

type SceneData struct {
	Name string
}

type sceneFactoryEntities struct {
	scene *Scene
}

func (f sceneFactoryEntities) AddEntities(ent ...[]Entity) Scene {
	for i := 0; i <= len(ent); i++ {
		f.scene.Entities = ent[0]
	}

	return *f.scene
}

func CreateScene(name string) sceneFactoryEntities {
	var SceneComponent = CreateComponent[SceneData]()
	scene := new(Scene)
	scene.SceneComponent = SceneComponent

	scene.Name = name

	factory := sceneFactoryEntities{
		scene: scene,
	}

	return factory
}
