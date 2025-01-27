//go:build !server

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package components

import (
	input "github.com/quasilyte/ebitengine-input"
	ecs "github.com/yohamta/donburi"
	"gomp/internal/tomb-mates-demo-v2/protos"
)

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
)

type LocalControllerData struct {
	input       *input.Handler
	playerInput struct {
		X, Y float64
	}
}

var LocalController = ecs.NewComponentType[*LocalControllerData](&LocalControllerData{
	input: nil,
	playerInput: struct{ X, Y float64 }{
		X: 0,
		Y: 0,
	},
})

func (lc *LocalControllerData) Update() {
	var event *protos.Event
	pInput := lc.playerInput

	if lc.input.ActionIsPressed(ActionMoveLeft) {
		pInput.X = -1
	} else if lc.input.ActionIsPressed(ActionMoveRight) {
		pInput.X = 1
	} else {
		pInput.X = 0
	}

	if lc.input.ActionIsPressed(ActionMoveUp) {
		pInput.Y = 1
	} else if lc.input.ActionIsPressed(ActionMoveDown) {
		pInput.Y = -1
	} else {
		pInput.Y = 0
	}

	if pInput != lc.playerInput {
		event = &protos.Event{
			Type: protos.EventType_move,
			Data: &protos.Event_Move{
				Move: &protos.EventMove{
					Direction: &protos.Vector2{
						X: pInput.X,
						Y: pInput.Y,
					},
				},
			},
			// PlayerId: *g.NetworkManager.MyID,
		}

		lc.playerInput = pInput
	}

	if event != nil {
		// message, err := proto.Marshal(event)
		// if err != nil {
		// 	return err
		// }

		// err = c.Write(context.Background(), websocket.MessageBinary, message)
		// if err != nil {
		// 	return err
		// }
	}
}
