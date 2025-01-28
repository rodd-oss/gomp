//go:build !graphics
// +build !graphics

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type ebitenGame struct {
	game         *Game
	dt           time.Duration
	lastUpdateAt time.Time
	inputSystem  *input.System
}

func (e *ebitenGame) Update() error {
	e.dt = time.Since(e.lastUpdateAt)

	e.inputSystem.UpdateWithDelta(e.dt.Seconds())
	e.game.Update(e.dt.Seconds())

	e.lastUpdateAt = time.Now()
	return nil
}

func (e *ebitenGame) Draw(screen *ebiten.Image) {
	// e.game.Draw(screen)

	op := &ebiten.DrawImageOptions{}

	query := donburi.NewQuery(filter.Contains(BodyComponent.Query, RenderComponent.Query))

	query.Each(e.game.World, func(e *donburi.Entry) {
		render := RenderComponent.Query.Get(e)

		if render == nil {
			log.Fatalln("RenderComponent is nil")
		}

		body := BodyComponent.Query.Get(e)

		if body == nil {
			log.Fatalln("BodyComponent is nil")
		}

		op.GeoM.Reset()
		op.GeoM.Translate(body.Position().X, body.Position().Y)

		screen.DrawImage(render.Sprite, op)
	})
}

func (e *ebitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (e *ebitenGame) RegisterInputHandlers(...input.Handler) {
	// e.inputSystem
}

// TODO: create input system that could register the gomp.CreateKeyboardHandler() and gomp.CreateMouseHandler()
