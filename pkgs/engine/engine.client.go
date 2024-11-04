//go:build !renderless

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import (
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

func (e *Engine) RunClient() (err error) {
	client := &Client{}
	client.engine = e

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Engine")
	tps := int(1000 / e.tickRate.Milliseconds())
	log.Println("TPS:", tps)
	ebiten.SetTPS(tps)

	if err = ebiten.RunGame(client); err != nil {
		return err
	}

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image, dt float64) {
	e.mx.Lock()
	defer e.mx.Unlock()

	if e.DebugDraw {
		log.Println("=========ENGINE DRAW START==========")
		defer log.Println("=========ENGINE DRAW FINISH=========")
	}

	sceneLen := len(e.LoadedScenes)
	if sceneLen == 0 {
		if e.Debug {
			log.Println("NO ACTIVE SCENES")
		}

		return
	}

	e.wg.Add(sceneLen)
	for i := range e.LoadedScenes {
		if !e.LoadedScenes[i].ShouldRender {
			e.wg.Done()
			continue
		}

		go drawSceneAsync(e.LoadedScenes[i], screen, dt, e.wg)
	}

	e.wg.Wait()
}

func drawSceneAsync(scene *Scene, screen *ebiten.Image, dt float64, wg *sync.WaitGroup) {
	defer wg.Done()
	scene.Draw(screen, dt)
}
