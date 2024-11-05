/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"context"
	"fmt"
	"gomp_game/pkgs/gomp/ecs"
	"log"
	"sync"
	"time"

	"github.com/yohamta/donburi"
)

type Game struct {
	mx sync.Mutex
	wg *sync.WaitGroup

	world        donburi.World
	systems      []ecs.System
	LoadedScenes map[string]*Scene

	tickRate time.Duration
	Debug    bool
}

func (g *Game) Init(tickRate time.Duration) {
	g.world = donburi.NewWorld()
	g.systems = []ecs.System{}
	g.tickRate = tickRate
	g.wg = new(sync.WaitGroup)
	g.LoadedScenes = make(map[string]*Scene)
	g.Debug = false
}

func (g *Game) Run(ctx context.Context) {
	ticker := time.NewTicker(g.tickRate)
	defer ticker.Stop()

	dt := g.tickRate.Seconds()

	g.Update(dt)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			g.Update(dt)
		}
	}
}

func (g *Game) Update(dt float64) {
	g.mx.Lock()
	defer g.mx.Unlock()

	if g.Debug {
		log.Println("=========Game UPDATE START==========")
		defer log.Println("=========Game UPDATE FINISH=========")
		log.Println("dt:", dt)
	}

	if dt > g.tickRate.Seconds()*2 {
		if g.Debug {
			log.Println("WARNING: Game tick rate is too high")
		}
		return
	}

	systemsLen := len(g.systems)
	if systemsLen == 0 {
		if g.Debug {
			log.Println("NO ACTIVE SYSTEMS")
		}

		return
	}

	for i := 0; i < systemsLen; i++ {
		g.systems[i].Update(dt)
	}

	g.updateAsync(dt)
}

func (g *Game) LoadScene(scene Scene) *Scene {
	g.mx.Lock()
	defer g.mx.Unlock()

	if g.Debug {
		log.Println("Loading scene:", scene.Name)
		defer log.Println("Scene loaded:", scene.Name)
	}

	//check if scene already exists
	suffix := 1

	for {
		prefixedName := fmt.Sprintf("%s_%d", scene.Name, suffix)

		if _, ok := g.LoadedScenes[prefixedName]; ok {
			suffix++
			continue
		}

		scene.Name = prefixedName
		break
	}

	g.LoadedScenes[scene.Name] = &scene
	return &scene
}

func (g *Game) UnloadScene(scene *Scene) {
	g.mx.Lock()
	defer g.mx.Unlock()

	if scene == nil {
		panic("Trying to unload nil scene")
	}

	name := scene.Name

	if g.Debug {
		log.Println("Unloading scene: ", name)
		defer log.Println("Scene unloaded: ", name)
	}

	// check if scene exists
	if _, ok := g.LoadedScenes[name]; !ok {
		return
	}

	delete(g.LoadedScenes, name)
}

func (g *Game) UnloadAllScenes() {
	for i := range g.LoadedScenes {
		g.UnloadScene(g.LoadedScenes[i])
	}
}

func (g *Game) updateAsync(dt float64) {
	g.wg.Add(len(g.systems))
	for i := range g.LoadedScenes {
		go updateSceneAsync(g.LoadedScenes[i], dt, g.wg)
	}
	g.wg.Wait()
}

func updateSceneAsync(scene *Scene, dt float64, wg *sync.WaitGroup) {
	defer wg.Done()
}
