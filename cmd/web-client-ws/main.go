package main

import (
	"context"
	"syscall/js"
	"tomb_mates/internal/client"
	"tomb_mates/internal/game"

	e "github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

func main() {
	url := js.Global().Get("document").Get("location").Get("origin").String() + "/ws"
	url = "ws" + url[4:]

	inputs := client.NewInputs(input.AnyDevice)
	transport := client.NewWsTransport(context.Background(), url)
	config := client.NewConfig(640, 480, "MMO 360 no scope", e.WindowResizingModeEnabled)

	gameClient := game.NewClient(inputs, transport, config)

	gameClient.Run()
}
