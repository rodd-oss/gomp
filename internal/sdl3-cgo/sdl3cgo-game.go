package main

import "C"
import (
	"fmt"
	"github.com/jfreymuth/go-sdl3/sdl"
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
	must(sdl.Init(sdl.InitVideo))
	defer sdl.Quit()

	w, r, e := sdl.CreateWindowAndRenderer("sld c-go", 1280, 720, sdl.WindowResizable)
	must(e)
	defer w.Destroy()
	defer r.Destroy()

	n, err := r.Name()
	fmt.Println(n, sdl.GetRenderDriver(0))

	var renderTicker *time.Ticker
	if framerate > 0 {
		renderTicker = time.NewTicker(time.Second / time.Duration(framerate))
		defer renderTicker.Stop()
	}

	var dt time.Duration = time.Second
	var fps, avgFps float64
	updatedAt := time.Now()

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
				if event.Keyboard().Scancode == sdl.ScancodeEscape {
					break Outer
				}
			}
		}

		must(r.SetDrawColor(0, 0, 0, 255))
		must(r.Clear())

		must(r.SetDrawColor(255, 255, 255, 255))

		for i := 0; i < len(rects); i += batchSize {
			must(r.FillRects(rects[i : i+batchSize]))
		}

		must(r.DebugText(10, 10, fmt.Sprintf("FPS %f", fps)))
		must(r.DebugText(10, 20, fmt.Sprintf("Avg FPS %f", avgFps)))
		must(r.DebugText(10, 30, fmt.Sprintf("dt %s", dt.String())))

		must(r.Present())
	}
}

func must(err error) {
	if err != nil {
		println(err.Error())
		panic(err)
	}
}
