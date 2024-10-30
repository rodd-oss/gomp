package engine

import (
	"context"
	"log"
	"sync"
	"time"
)

type Engine struct {
	wg *sync.WaitGroup

	Network  *Network
	Scenes   []*Scene
	tickRate time.Duration
}

func NewEngine(tickRate time.Duration) *Engine {
	e := new(Engine)

	e.tickRate = tickRate
	e.wg = new(sync.WaitGroup)

	return e
}

func (e *Engine) Run(ctx context.Context) {
	ticker := time.NewTicker(e.tickRate)
	defer ticker.Stop()

	dt := e.tickRate.Seconds()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.Update(dt)
		}
	}
}

func (e *Engine) Update(dt float64) {
	log.Println("Engine update")
	e.Network.Update()

	e.wg.Add(len(e.Scenes))

	for i := range e.Scenes {
		go updateSceneAsync(e.Scenes[i], dt, e.wg)
	}

	e.wg.Wait()
}

func (e *Engine) LoadScene(scene *Scene) {
	e.Scenes = append(e.Scenes, scene)
}

func updateSceneAsync(scene *Scene, dt float64, wg *sync.WaitGroup) {
	defer wg.Done()
	scene.Update(dt)
}
