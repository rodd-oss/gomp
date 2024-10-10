package components

import (
	"tomb_mates/internal/protos"

	input "github.com/quasilyte/ebitengine-input"
	ecs "github.com/yohamta/donburi"
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
