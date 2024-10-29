package engine

import "sync"

type Engine struct {
	wg *sync.WaitGroup

	Network *Network
	Scenes  []*Scene
}

func NewEngine() *Engine {
	e := new(Engine)

	e.wg = new(sync.WaitGroup)

	return e
}

func (e *Engine) Update(dt float64) {
	e.Network.Update()

	e.wg.Add(len(e.Scenes))

	for i := range e.Scenes {
		go updateSceneAsync(e.Scenes[i], dt, e.wg)
	}

	e.wg.Wait()
}

func updateSceneAsync(scene *Scene, dt float64, wg *sync.WaitGroup) {
	defer wg.Done()
	scene.Update(dt)
}
