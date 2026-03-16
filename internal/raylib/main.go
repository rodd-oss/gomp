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
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	batchSize    = 1 << 14
	rectCounter  = 1 << 18
	framerate    = 600
	frameCount   = 100 // Number of frames to average
)

func main() {
	f, err := os.Create("cpu.out")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	fmt.Println("CPU Profile Started")
	defer fmt.Println("CPU Profile Stopped")

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	rl.InitWindow(screenWidth, screenHeight, "raylib example")
	defer rl.CloseWindow()

	rl.SetTargetFPS(framerate)

	rects := make([]rl.Rectangle, rectCounter)
	for i := range rects {
		rects[i] = rl.Rectangle{X: 300, Y: 300, Width: 10, Height: 10}
	}

	var dt time.Duration = time.Second
	var fps, avgFps float64
	updatedAt := time.Now()

	frameTimes := make([]float64, 0, frameCount)

	for !rl.WindowShouldClose() {
		fps = float64(time.Second) / float64(dt)
		dt = time.Since(updatedAt)
		updatedAt = time.Now()

		frameTimes = append(frameTimes, fps)
		if len(frameTimes) > frameCount {
			frameTimes = frameTimes[1:]
		}

		var totalFps float64
		for _, frameTime := range frameTimes {
			totalFps += frameTime
		}
		avgFps = totalFps / float64(len(frameTimes))

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.DrawText(fmt.Sprintf("FPS %f", fps), 10, 10, 20, rl.White)
		rl.DrawText(fmt.Sprintf("Avg FPS %f", avgFps), 10, 30, 20, rl.White)
		rl.DrawText(fmt.Sprintf("dt %s", dt.String()), 10, 50, 20, rl.White)

		for _, rect := range rects {
			rl.DrawRectangleRec(rect, rl.White)
		}

		rl.EndDrawing()
	}
}
