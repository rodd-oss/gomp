package components

import (
	"gomp_game/pkgs/gomp"
	"gomp_game/pkgs/gomp/ecs"

	"github.com/yohamta/donburi"
)

type Scene struct {
	Name string

	Systems        []ecs.System
	Entities       []ecs.Entity
	SceneComponent *donburi.ComponentType[SceneData]
}

type SceneData struct {
	ID uint8
}

var SceneComponent = gomp.CreateComponent[SceneData]
