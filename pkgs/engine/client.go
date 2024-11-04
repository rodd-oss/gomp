//go:build !renderless

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Client struct {
	engine *Engine
}

func (c *Client) Update() error {
	dt := 1 / ebiten.ActualTPS()

	c.engine.Update(dt)
	return nil
}

func (c *Client) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (c *Client) Draw(screen *ebiten.Image) {
	dt := 1 / ebiten.ActualFPS()

	c.engine.Draw(screen, dt)
}
