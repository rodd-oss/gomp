package engine

import (
	"context"
	"log"
	"sync"
	"time"
)

type Engine struct {
	wg *sync.WaitGroup

	Network          *Network
	registeredScenes map[string]*Scene
	tickRate         time.Duration
}

func NewEngine(tickRate time.Duration) *Engine {
	e := new(Engine)

	e.tickRate = tickRate
	e.wg = new(sync.WaitGroup)
	e.registeredScenes = make(map[string]*Scene)

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
	log.Println("===================")
	log.Println("ENGINE UPDATE START")
	e.Network.Update()

	e.wg.Add(len(e.registeredScenes))

	for i := range e.registeredScenes {
		go updateSceneAsync(e.registeredScenes[i], dt, e.wg)
	}

	e.wg.Wait()
	log.Println("ENGINE UPDATE FINISH")
	log.Println("===================")

}

func (e *Engine) RegisterScene(name string, scene Scene) {
	scene.Name = name
	e.registeredScenes[name] = &scene
}

func updateSceneAsync(scene *Scene, dt float64, wg *sync.WaitGroup) {
	defer wg.Done()
	scene.Update(dt)
}
