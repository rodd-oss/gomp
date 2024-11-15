/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/pkgs/example/components"
	"gomp_game/pkgs/gomp"
	"math"

	input "github.com/quasilyte/ebitengine-input"
	"github.com/yohamta/donburi"
)

var IntentSystem = gomp.CreateSystem(new(intentContoller))

const (
	heroActionMoveUp input.Action = iota
	heroActionMoveDown
	heroActionMoveLeft
	heroActionMoveRight
	heroActionJump
	heroActionFire
)

type intentContoller struct {
	world            donburi.World
	heroInputHandler *input.Handler
}

func (c *intentContoller) Init(game *gomp.Game) {
	c.world = game.World

	c.heroInputHandler = gomp.InputSystem.NewHandler(0, input.Keymap{
		heroActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
		heroActionMoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
		heroActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		heroActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		heroActionJump:      {input.KeyGamepadA, input.KeySpace},
	})
}

func (c *intentContoller) Update(dt float64) {
	components.HeroIntentComponent.Query.Each(c.world, func(e *donburi.Entry) {
		intent := components.HeroIntentComponent.Query.Get(e)

		if c.heroInputHandler.ActionIsPressed(heroActionMoveUp) {
			intent.Move.Y = 1
		} else if c.heroInputHandler.ActionIsPressed(heroActionMoveDown) {
			intent.Move.Y = -1
		} else {
			intent.Move.Y = 0
		}

		if c.heroInputHandler.ActionIsPressed(heroActionMoveLeft) {
			intent.Move.X = -1
		} else if c.heroInputHandler.ActionIsPressed(heroActionMoveRight) {
			intent.Move.X = 1
		} else {
			intent.Move.X = 0
		}

		// Normalize diagonal movement
		sumVector := math.Sqrt(intent.Move.X*intent.Move.X + intent.Move.Y*intent.Move.Y)
		if sumVector != 0 {
			intent.Move.X = intent.Move.X / sumVector
			intent.Move.Y = intent.Move.Y / sumVector
		}

		intent.Jump = c.heroInputHandler.ActionIsJustPressed(heroActionJump)

		intent.Fire = c.heroInputHandler.ActionIsPressed(heroActionFire)
	})
}
