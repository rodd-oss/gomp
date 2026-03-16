package main

import (
	"fmt"
	"github.com/jupiterrider/purego-sdl3/sdl"
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

	sdl.SetHint(sdl.HintRenderDriver, "gpu,software")
	defer sdl.Quit()
	if !sdl.Init(sdl.InitVideo) {
		panic(sdl.GetError())
	}

	var w *sdl.Window
	var r *sdl.Renderer

	if !sdl.CreateWindowAndRenderer("Hello, World!", 1280, 720, sdl.WindowResizable, &w, &r) {
		panic(sdl.GetError())
	}
	defer sdl.DestroyRenderer(r)
	defer sdl.DestroyWindow(w)

	fmt.Println(sdl.GetRendererName(r), sdl.GetRenderDriver(0))

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
			switch event.Type() {
			case sdl.EventQuit:
				break Outer
			case sdl.EventKeyDown:
				if event.Key().Scancode == sdl.ScancodeEscape {
					break Outer
				}
			}
		}
		sdl.SetRenderDrawColor(r, 0, 0, 0, 255)
		sdl.RenderClear(r)

		sdl.SetRenderDrawColor(r, 255, 255, 255, 255)
		for i := 0; i < len(rects); i += batchSize {
			sdl.RenderFillRects(r, rects[i:i+batchSize]...)
		}

		sdl.RenderDebugText(r, 10, 10, fmt.Sprintf("FPS %f", fps))
		sdl.RenderDebugText(r, 10, 20, fmt.Sprintf("Avg FPS %f", avgFps))
		sdl.RenderDebugText(r, 10, 30, fmt.Sprintf("dt %s", dt.String()))

		sdl.RenderPresent(r)
	}
}

func must(err error) {
	if err != nil {
		println(err.Error())
		panic(err)
	}
}
