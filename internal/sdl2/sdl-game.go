package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

	ok := sdl.SetHint(sdl.HINT_RENDER_DRIVER, "gpu,software")
	if !ok {
		panic(sdl.GetError())
	}
	must(ttf.Init())
	defer ttf.Quit()
	must(sdl.Init(sdl.INIT_EVERYTHING))
	defer sdl.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer(1280, 720, sdl.WINDOW_RESIZABLE)
	must(err)
	defer window.Destroy()
	defer renderer.Destroy()

	font, err := ttf.OpenFont("Roboto-SemiBold.ttf", 14)
	must(err)
	defer font.Close()

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

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break Outer
			case *sdl.KeyboardEvent:
				if event.(*sdl.KeyboardEvent).Keysym.Scancode == sdl.SCANCODE_ESCAPE {
					break Outer
				}
			}
		}

		must(renderer.SetDrawColor(0, 0, 0, 255))
		must(renderer.Clear())

		must(renderer.SetDrawColor(255, 255, 255, 255))

		for i := 0; i < len(rects); i += batchSize {
			must(renderer.FillRectsF(rects[i : i+batchSize]))
		}

		must(drawText(renderer, font, 10, 10, fmt.Sprintf("FPS %f", fps)))
		must(drawText(renderer, font, 10, 30, fmt.Sprintf("Avg FPS %f", avgFps)))
		must(drawText(renderer, font, 10, 50, fmt.Sprintf("dt %s", dt.String())))

		renderer.Present()
	}
}

func must(err error) {
	if err != nil {
		println(err.Error())
		panic(err)
	}
}

// Cache all char textures on the fly
// TODO: preload all chars initially (maybe +1FPS)
var letterCache [256]*sdl.Texture

func drawText(renderer *sdl.Renderer, font *ttf.Font, x, y int32, text string) error {
	var totalWidth int32
	var maxHeight int32

	for _, char := range text {
		if char < 256 {
			if texture := letterCache[char]; texture != nil {
				_, _, w, h, err := texture.Query()
				if err != nil {
					return err
				}
				dst := sdl.Rect{X: x + totalWidth, Y: y, W: w, H: h}
				if err := renderer.Copy(texture, nil, &dst); err != nil {
					return err
				}
				totalWidth += w
				if h > maxHeight {
					maxHeight = h
				}
				continue
			}
		}

		surface, err := font.RenderUTF8Blended(string(char), sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err != nil {
			return err
		}
		defer surface.Free()

		texture, err := renderer.CreateTextureFromSurface(surface)
		if err != nil {
			return err
		}
		if char < 256 {
			letterCache[char] = texture
		}

		_, _, w, h, err := texture.Query()
		if err != nil {
			return err
		}

		dst := sdl.Rect{X: x + totalWidth, Y: y, W: w, H: h}
		if err := renderer.Copy(texture, nil, &dst); err != nil {
			return err
		}
		totalWidth += w
		if h > maxHeight {
			maxHeight = h
		}
	}

	return nil
}
