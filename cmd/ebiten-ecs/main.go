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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type pixel struct {
	x      int
	y      int
	breath bool
	hp     int
	color  color.RGBA
}

var pixelComponentType = ecs.CreateComponent[pixel]()

type pixelSystem struct {
	pixelComponent ecs.WorldComponents[pixel]
}

func (s *pixelSystem) Init(world *ecs.World) {
	s.pixelComponent = pixelComponentType.Instances(world)

	for i := range 1000 {
		for j := range 1000 {
			newPixel := world.CreateEntity("Pixel")

			randomGreen := uint8(135 / (rand.Intn(10) + 1))
			randomBlue := uint8(135 / (rand.Intn(10) + 1))

			randomColor := color.RGBA{
				G: randomGreen,
				B: randomBlue,
				A: 255,
			}
			s.pixelComponent.Set(newPixel, pixel{
				x:     j,
				y:     i,
				hp:    100,
				color: randomColor,
			})
		}
	}
}

func (s *pixelSystem) Run(world *ecs.World) {
	for _, pixel := range s.pixelComponent.All() {
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
	}
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
	for _, pixel := range g.pixelComponents.All() {
		index := (pixel.x + pixel.y*1080) * 4
		color := &pixel.color

		g.screenBuffer[index+0] = color.R
		g.screenBuffer[index+1] = color.G
		g.screenBuffer[index+2] = color.B
		g.screenBuffer[index+3] = color.A
	}

	g.imageBuffer.WritePixels(g.screenBuffer)

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
	ebiten.SetWindowSize(1080, 1080)
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
		imageBuffer:     ebiten.NewImage(1080, 1080),
	}

	newGame.screenBuffer = make([]byte, 4*1080*1080)

	if err := ebiten.RunGame(&newGame); err != nil {
		panic(err)
	}
}
