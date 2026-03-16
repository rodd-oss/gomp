/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package worker

import (
	"context"
	"sync"
)

type WorkerId int

func NewWorker(ctx context.Context, id WorkerId, taskChan <-chan AnyTask, taskWg *sync.WaitGroup) Worker {
	return Worker{
		id:       id,
		ctx:      ctx,
		taskWg:   taskWg,
		taskChan: taskChan,
	}
}

type Worker struct {
	id       WorkerId
	ctx      context.Context
	taskChan <-chan AnyTask
	taskWg   *sync.WaitGroup
	wg       sync.WaitGroup
}

func (w *Worker) Start(poolWg *sync.WaitGroup) {
	poolWg.Add(1)
	go w.run(poolWg)
}

func (w *Worker) Stop() {
}

func (w *Worker) run(poolWg *sync.WaitGroup) {
	defer poolWg.Done()
	for {
		select {
		case <-w.ctx.Done():
			return
		case task := <-w.taskChan:
			err := task.Run(w.id)
			w.taskWg.Done()
			if err != nil {
				panic("not implemented")
			}
		}
	}
}

func (w *Worker) Id() WorkerId {
	return w.id
}
