//go:build !renderless

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func (s *Scene) Draw(screen *ebiten.Image, dt float64) {
	if s.Engine.DebugDraw {
		log.Println("Scene Draw:", s.Name)
		defer log.Println("Scene Draw:", s.Name)
	}

	for i := range s.RenderSystems {
		s.RenderSystems[i].Draw(screen, dt)
	}
}
