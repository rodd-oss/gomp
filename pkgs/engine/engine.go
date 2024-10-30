package engine

import (
	"context"
	"log"
	"sync"
	"time"
)

type Engine struct {
	mx sync.Mutex
	wg *sync.WaitGroup

	Network      *Network
	LoadedScenes map[string]*Scene
	tickRate     time.Duration
}

func NewEngine(tickRate time.Duration) *Engine {
	e := new(Engine)

	e.tickRate = tickRate
	e.wg = new(sync.WaitGroup)
	e.LoadedScenes = make(map[string]*Scene)

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
	e.mx.Lock()
	defer e.mx.Unlock()

	log.Println("=========ENGINE UPDATE START==========")
	defer log.Println("=========ENGINE UPDATE FINISH=========")

	e.Network.Update()

	sceneLen := len(e.LoadedScenes)
	if sceneLen == 0 {
		log.Println("NO ACTIVE SCENES")
		return
	}

	e.wg.Add(sceneLen)
	for i := range e.LoadedScenes {
		go updateSceneAsync(e.LoadedScenes[i], dt, e.wg)
	}

	e.wg.Wait()
}

func (e *Engine) LoadScene(name string, scene Scene) {
	e.mx.Lock()
	defer e.mx.Unlock()

	scene.Engine = e
	scene.Name = name

	scene.Contoller.OnLoad(e.LoadedScenes[name])

	e.LoadedScenes[name] = &scene
}

func (e *Engine) UnloadScene(name string) {
	e.mx.Lock()
	defer e.mx.Unlock()

	// check if scene exists
	if _, ok := e.LoadedScenes[name]; !ok {
		return
	}

	e.LoadedScenes[name].Contoller.OnUnload(e.LoadedScenes[name])

	delete(e.LoadedScenes, name)
}

func (e *Engine) UnloadAllScenes() {
	for i := range e.LoadedScenes {
		e.UnloadScene(i)
	}
}

func updateSceneAsync(scene *Scene, dt float64, wg *sync.WaitGroup) {
	defer wg.Done()
	scene.Update(dt)
}
