/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- Еблан Donated 228 RUB
<- VsehVertela Donated 500 RUB
<- Linkwayz Donated 500 RUB
<- thespacetime Donated 10 EUR
<- Linkwayz Donated 1500 RUB
<- mitwelve Donated 100 RUB
<- tema881 Donated 100 RUB

Thank you for your support!
*/

package worker

import (
	"context"
	"sync"
)

func NewPool(n int) Pool {
	return Pool{
		workers:  make([]Worker, n),
		taskChan: make(chan AnyTask, n*4),
	}
}

type Pool struct {
	workers   []Worker
	wg        sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc

	// Cache
	taskChan    chan AnyTask
	groupTaskWg *sync.WaitGroup
}

func (p *Pool) Start() {
	p.ctx, p.ctxCancel = context.WithCancel(context.Background())
	p.groupTaskWg = new(sync.WaitGroup)
	p.wg = sync.WaitGroup{}
	for i := range p.workers {
		p.workers[i] = NewWorker(p.ctx, WorkerId(i), p.taskChan, p.groupTaskWg)
		p.workers[i].Start(&p.wg)
	}
}

func (p *Pool) AddWorker() {
	p.workers = append(p.workers, NewWorker(p.ctx, WorkerId(len(p.workers)), p.taskChan, p.groupTaskWg))
	p.workers[len(p.workers)-1].Start(&p.wg)
}

func (p *Pool) RemoveWorker() {
	p.workers[len(p.workers)-1].Stop()
	p.workers = p.workers[:len(p.workers)-1]
}

func (p *Pool) Stop() {
	p.ctxCancel()
	p.wg.Wait()
}

func (p *Pool) GroupAdd(n int) {
	p.groupTaskWg.Add(n)
}

func (p *Pool) ProcessGroupTask(tasks AnyTask) {
	p.taskChan <- tasks
}

func (p *Pool) GroupWait() {
	p.groupTaskWg.Wait()
}

func (p *Pool) NumWorkers() int {
	return len(p.workers)
}
