/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"time"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	rectCounter  = 1 << 18
	frameCount   = 100
)

type Rect struct {
	X float32
	Y float32
	W float32
	H float32
}

type Game struct {
	rects      []Rect
	frameTimes []float64
	updatedAt  time.Time
	dt         time.Duration
	fps        float64
	avgFps     float64
}

func (g *Game) Update() error {
	g.dt = time.Since(g.updatedAt)
	g.updatedAt = time.Now()
	g.fps = 1.0 / g.dt.Seconds()

	g.frameTimes = append(g.frameTimes, g.fps)
	if len(g.frameTimes) > frameCount {
		g.frameTimes = g.frameTimes[1:]
	}

	var totalFps float64
	for _, frameTime := range g.frameTimes {
		totalFps += frameTime
	}
	g.avgFps = totalFps / float64(len(g.frameTimes))

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("quit")
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for _, rect := range g.rects {
		vector.DrawFilledRect(screen, rect.X, rect.Y, rect.W, rect.H, color.White, false)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nAvg FPS: %0.2f\nDT: %s", g.fps, g.avgFps, g.dt))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		rects:      make([]Rect, rectCounter),
		frameTimes: make([]float64, 0, frameCount),
		updatedAt:  time.Now(),
	}

	for i := range game.rects {
		game.rects[i] = Rect{X: 300, Y: 300, W: 10, H: 10}
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(game); err != nil && err.Error() != "quit" {
		log.Fatal(err)
	}
}
