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
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const batchSize = 1 << 14
const rectCounter = 1 << 18
const framerate = 600
const frameCount = 100 // Number of frames to average

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

	defer binsdl.Load().Unload() // sdl.LoadLibrary(sdl.Path())

	sdl.SetHint(sdl.HINT_RENDER_DRIVER, "gpu,software") // <- gpu driver creates memory leak
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err.Error())
	}
	defer sdl.Quit()

	//var w *sdl.Window
	//var r *sdl.Renderer
	w, r, err := sdl.CreateWindowAndRenderer("Hello, World!", 1280, 720, sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err.Error())
	}
	defer r.Destroy()
	defer w.Destroy()

	var dt time.Duration = time.Second
	var fps, avgFps float64
	updatedAt := time.Now()

	var renderTicker *time.Ticker
	if framerate > 0 {
		renderTicker = time.NewTicker(time.Second / time.Duration(framerate))
		defer renderTicker.Stop()
	}

	rects := make([]sdl.FRect, rectCounter)
	for i := range rects {
		rects[i] = sdl.FRect{X: 300, Y: 300, W: 10, H: 10}
	}

	fmt.Println(sdl.GetRenderDriver(0))
	frameTimes := make([]float64, 0, frameCount)

Outer:
	for {
		if renderTicker != nil {
			<-renderTicker.C
		}
		fps = time.Second.Seconds() / dt.Seconds()
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

		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type {
			case sdl.EVENT_QUIT:
				break Outer
			case sdl.EVENT_KEY_DOWN:
				if event.KeyboardEvent().Scancode == sdl.SCANCODE_ESCAPE {
					break Outer
				}
			}
		}
		r.SetDrawColor(0, 0, 0, 255)
		r.Clear()

		r.SetDrawColor(255, 255, 255, 255)
		for i := 0; i < len(rects); i += batchSize {
			r.RenderFillRects(rects[i : i+batchSize])
		}

		r.DebugText(10, 10, fmt.Sprintf("FPS %f", fps))
		r.DebugText(10, 20, fmt.Sprintf("Avg FPS %f", avgFps))
		r.DebugText(10, 30, fmt.Sprintf("dt %s", dt.String()))

		r.Present()
	}
}

func must(err error) {
	if err != nil {
		println(err.Error())
		panic(err)
	}
}
