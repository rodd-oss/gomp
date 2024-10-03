package components

import (
	e "github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/anim"
	ecs "github.com/yohamta/donburi"
)

type RenderData struct {
	Image      *e.Image
	Controller *anim.AnimationPlayer
}

var Render = ecs.NewComponentType[RenderData]()
