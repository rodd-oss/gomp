/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"fmt"
	"gomp_game/pkgs/gomp/ecs"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type cameraSystem struct {
	transformComponent   ecs.WorldComponents[transform]
	healthComponent      ecs.WorldComponents[health]
	colorComponent       ecs.WorldComponents[color.RGBA]
	movementComponent    ecs.WorldComponents[movement]
	cameraComponent      ecs.WorldComponents[camera]
	destroyComponentType ecs.WorldComponents[empty]

	buffer    *image.RGBA
	debugInfo []string
	p         *message.Printer
}

func (s *cameraSystem) Init(world *ecs.World) {
	s.transformComponent = transformComponentType.Instances(world)
	s.healthComponent = healthComponentType.Instances(world)
	s.colorComponent = colorComponentType.Instances(world)
	s.movementComponent = movementComponentType.Instances(world)
	s.cameraComponent = cameraComponentType.Instances(world)
	s.destroyComponentType = destroyComponentType.Instances(world)

	s.p = message.NewPrinter(language.Russian)

	newcamera := world.CreateEntity("camera")
	s.cameraComponent.Set(newcamera, camera{
		mainLayer: cameraLayer{
			image: ebiten.NewImage(width, height),
			zoom:  1,
		},
		debugLayer: cameraLayer{
			image: ebiten.NewImage(width, height),
			zoom:  2,
		},
	})

	s.buffer = image.NewRGBA(image.Rect(0, 0, width, height))
	s.debugInfo = make([]string, 0)
}

func (s *cameraSystem) Run(world *ecs.World) {
	_, dy := ebiten.Wheel()

	draw.Draw(s.buffer, s.buffer.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	s.colorComponent.AllParallel(func(entity ecs.EntityID, color *color.RGBA) bool {
		if color == nil {
			return true
		}

		transform := s.transformComponent.GetPtr(entity)
		if transform == nil {
			return true
		}
		s.buffer.SetRGBA(int(transform.x), int(transform.y), *color)
		return true
	})

	var mainCamera *camera
	s.cameraComponent.AllData(func(c *camera) bool {
		mainCamera = c
		return false
	})

	mainCamera.mainLayer.image.Clear()
	mainCamera.debugLayer.image.Clear()

	mainCamera.mainLayer.zoom += float64(dy)
	mainCamera.mainLayer.image.WritePixels(s.buffer.Pix)

	if mainCamera.mainLayer.zoom < 0.1 {
		mainCamera.mainLayer.zoom = 0.1
	} else if mainCamera.mainLayer.zoom > 100 {
		mainCamera.mainLayer.zoom = 100
	}

	s.debugInfo = append(s.debugInfo, fmt.Sprintf("TPS %0.2f", ebiten.ActualTPS()))
	s.debugInfo = append(s.debugInfo, fmt.Sprintf("FPS %0.2f", ebiten.ActualFPS()))
	s.debugInfo = append(s.debugInfo, fmt.Sprintf("Zoom %0.2f", mainCamera.mainLayer.zoom))
	s.debugInfo = append(s.debugInfo, s.p.Sprintf("Entity count %d", entityCount))
	s.debugInfo = append(s.debugInfo, s.p.Sprintf("Transforms count %d", s.transformComponent.Len()))
	s.debugInfo = append(s.debugInfo, s.p.Sprintf("Healths count %d", s.healthComponent.Len()))
	s.debugInfo = append(s.debugInfo, s.p.Sprintf("Colors count %d", s.colorComponent.Len()))
	s.debugInfo = append(s.debugInfo, s.p.Sprintf("Movements count %d", s.movementComponent.Len()))
	s.debugInfo = append(s.debugInfo, s.p.Sprintf("Cameras count %d", s.cameraComponent.Len()))
	s.debugInfo = append(s.debugInfo, s.p.Sprintf("Destroys count %d", s.destroyComponentType.Len()))
	s.debugInfo = append(s.debugInfo, s.p.Sprintf("Pprof %d", pprofEnabled))

	ebitenutil.DebugPrint(mainCamera.debugLayer.image, strings.Join(s.debugInfo, "\n"))
	s.debugInfo = s.debugInfo[:0]
}
func (s *cameraSystem) Destroy(world *ecs.World) {}
