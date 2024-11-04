//go:build !renderless

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import "github.com/hajimehoshi/ebiten/v2"

type RenderSystemController interface {
	Init(scene *Scene)
	Draw(screen *ebiten.Image, dt float64)
}

type RenderSystem struct {
	controller RenderSystemController
}

func (s *RenderSystem) Init(scene *Scene) {
	s.controller.Init(scene)
}

func (s *RenderSystem) Draw(screen *ebiten.Image, dt float64) {
	s.controller.Draw(screen, dt)
}

func CreateRenderSystem(controller RenderSystemController) RenderSystem {
	return RenderSystem{
		controller: controller,
	}
}
