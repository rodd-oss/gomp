/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"context"
	"syscall/js"
	"tomb_mates/internal/client"

	e "github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

func main() {
	ctx := context.Background()

	url := js.Global().Get("document").Get("location").Get("origin").String() + "/ws"
	url = "ws" + url[4:]

	inputs := client.NewInputs(input.AnyDevice)
	config := client.NewConfig(640, 480, "MMO 360 no scope", e.WindowResizingModeEnabled, true)
	transport := client.NewWsTransport(ctx, url)

	gameClient := client.New(ctx, inputs, transport, config)

	err := gameClient.Run()
	if err != nil {
		panic(err)
	}
}
