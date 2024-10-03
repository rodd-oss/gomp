package components

import (
	ecs "github.com/yohamta/donburi"
)

type RenderData struct {
	// Image      *e.Image
	// Controller *anim.AnimationPlayer
}

var Render = ecs.NewComponentType[RenderData]()
