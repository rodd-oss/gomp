package client

import input "github.com/quasilyte/ebitengine-input"

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
)

type Inputs struct {
	System   input.System
	keyMaps  map[uint8]*input.Keymap
	Handlers map[uint8]*input.Handler
}

func NewInputs(availableDevices input.DeviceKind) *Inputs {
	inputs := &Inputs{
		keyMaps:  make(map[uint8]*input.Keymap),
		Handlers: make(map[uint8]*input.Handler),
	}

	inputs.System.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})

	return inputs
}

func (i *Inputs) Register(deviceId uint8, keyMap *input.Keymap) {
	handler := i.System.NewHandler(deviceId, *keyMap)

	i.keyMaps[deviceId] = keyMap
	i.Handlers[deviceId] = handler
}
