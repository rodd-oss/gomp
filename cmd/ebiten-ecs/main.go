/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"flag"
	"fmt"
	"gomp_game/pkgs/gomp/ecs"
	"image/color"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type pixel struct {
	x      int32
	y      int32
	hp     int32
	color  color.RGBA
	breath bool
}

const (
	width  = 1000
	height = 1000
	scale  = float64(4)
)

var pixelComponentType = ecs.CreateComponent[pixel]()

type pixelSystem struct {
	pixelComponent ecs.WorldComponents[pixel]
}

func (s *pixelSystem) Init(world *ecs.World) {
	s.pixelComponent = pixelComponentType.Instances(world)

	for i := range height {
		for j := range width {
			newPixel := world.CreateEntity("Pixel")

			randomGreen := uint8(135 / (rand.Intn(10) + 1))
			randomBlue := uint8(135 / (rand.Intn(10) + 1))

			randomColor := color.RGBA{
				G: randomGreen,
				B: randomBlue,
				A: 255,
			}
			s.pixelComponent.Set(newPixel, pixel{
				x:     int32(j),
				y:     int32(i),
				hp:    100,
				color: randomColor,
			})
		}
	}
}

func (s *pixelSystem) Run(world *ecs.World) {
	s.pixelComponent.AllDataParallel(func(pixel *pixel) bool {
		color := &pixel.color

		if pixel.breath {
			if color.G < 135 {
				color.G++
			} else {
				pixel.hp++
			}
			if color.B < 135 {
				color.B++
			} else {
				pixel.hp++
			}
		} else {
			if color.G > 0 {
				color.G--
			} else {
				pixel.hp--
			}
			if color.B > 0 {
				color.B--
			} else {
				pixel.hp--
			}
		}

		if pixel.hp <= 0 {
			pixel.breath = true
		} else if pixel.hp >= 100 {
			pixel.breath = false
		}
		return true
	})
}

func (s *pixelSystem) Destroy(world *ecs.World) {}

type game struct {
	world           *ecs.World
	pixelComponents ecs.WorldComponents[pixel]

	imageBuffer  *ebiten.Image
	screenBuffer []byte
}

func (g *game) Update() error {
	g.world.RunSystems()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	g.pixelComponents.AllDataParallel(func(pixel *pixel) bool {
		index := (pixel.x + pixel.y*width) * 4
		*(*[4]byte)(unsafe.Pointer(&g.screenBuffer[index])) = *(*[4]byte)(unsafe.Pointer(&pixel.color))
		return true
	})

	g.imageBuffer.WritePixels(g.screenBuffer)
	op.GeoM.Scale(scale, scale)

	screen.DrawImage(g.imageBuffer, op)

	var debugInfo = make([]string, 0)

	debugInfo = append(debugInfo, fmt.Sprintf("TPS %0.2f", ebiten.ActualTPS()))
	debugInfo = append(debugInfo, fmt.Sprintf("FPS %0.2f", ebiten.ActualFPS()))

	ebitenutil.DebugPrint(screen, strings.Join(debugInfo, "\n"))
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("1 mil pixel ecs")

	world := ecs.New("1 mil pixel")

	world.RegisterComponentTypes(
		&pixelComponentType,
	)

	world.RegisterSystems().Parallel(
		new(pixelSystem),
	)

	newGame := game{
		world:           &world,
		pixelComponents: pixelComponentType.Instances(&world),
		imageBuffer:     ebiten.NewImage(width, height),
	}

	newGame.screenBuffer = make([]byte, 4*width*height)

	if err := ebiten.RunGame(&newGame); err != nil {
		panic(err)
	}
}
