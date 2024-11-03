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

	Debug bool
}

func NewEngine(tickRate time.Duration) *Engine {
	e := new(Engine)

	e.tickRate = tickRate
	e.wg = new(sync.WaitGroup)
	e.LoadedScenes = make(map[string]*Scene)
	e.Debug = false

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

	if e.Debug {
		log.Println("=========ENGINE UPDATE START==========")
		defer log.Println("=========ENGINE UPDATE FINISH=========")
	}

	e.Network.Update()

	sceneLen := len(e.LoadedScenes)
	if sceneLen == 0 {
		if e.Debug {
			log.Println("NO ACTIVE SCENES")
		}

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

	if e.Debug {
		log.Println("Loading scene:", name)
		defer log.Println("Scene loaded:", name)
	}

	scene.Engine = e
	scene.Name = name

	systems, entities := scene.Contoller.Load(&scene)

	for i := range entities {
		scene.World.Create(entities[i]...)
	}

	for i := range systems {
		systems[i].Init(&scene)
	}

	scene.Systems = systems
	e.LoadedScenes[name] = &scene
}

func (e *Engine) UnloadScene(name string) {
	e.mx.Lock()
	defer e.mx.Unlock()

	if e.Debug {
		log.Println("Unloading scene: ", name)
		defer log.Println("Scene unloaded: ", name)
	}

	// check if scene exists
	if _, ok := e.LoadedScenes[name]; !ok {
		return
	}

	e.LoadedScenes[name].Contoller.Unload(e.LoadedScenes[name])

	delete(e.LoadedScenes, name)
}

func (e *Engine) UnloadAllScenes() {
	for i := range e.LoadedScenes {
		e.UnloadScene(i)
	}
}

func (e *Engine) SetDebug(mode bool) {
	e.Debug = mode
}

func updateSceneAsync(scene *Scene, dt float64, wg *sync.WaitGroup) {
	defer wg.Done()
	scene.Update(dt)
}
