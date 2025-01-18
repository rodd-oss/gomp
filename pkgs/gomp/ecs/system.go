/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "sync"

type AnySystemPtr[W any] interface {
	Init(*W)
	Run(*W)
	Destroy(*W)
}

type SystemFunctionMethod int

const (
	systemFunctionInit SystemFunctionMethod = iota
	systemFunctionUpdate
	SystemFunctionFixedUpdate
	systemFunctionDestroy
)

type AnySystemControllerPtr interface {
	Init(*World)
	Update(*World)
	FixedUpdate(*World)
	Destroy(*World)
}

type AnySystemServicePtr interface {
	register(*World) *SystemServiceInstance
	registerDependencyFor(*World, AnySystemServicePtr)
	getInstance(*World) *SystemServiceInstance
}

type SystemServiceInstance struct {
	dependsOnNum    int
	controller      AnySystemControllerPtr
	dependencyFor   []*SystemServiceInstance
	dependencyWg    *sync.WaitGroup
	world           *World
	initChan        chan struct{}
	updateChan      chan struct{}
	fixedUpdateChan chan struct{}
	destroyChan     chan struct{}
}

func (s *SystemServiceInstance) PrepareWg() {
	s.dependencyWg.Add(s.dependsOnNum)
}

func (s *SystemServiceInstance) WaitForDeps() {
	s.dependencyWg.Wait()
}

func (s *SystemServiceInstance) SendDone() {
	for i := range s.dependencyFor {
		s.dependencyFor[i].dependencyWg.Done()
	}
}

func (s *SystemServiceInstance) Add(delta int) {
	s.dependencyWg.Add(delta)
}

func (s *SystemServiceInstance) Run() {
	controller := s.controller

	controller.Init(s.world)
	defer controller.Destroy(s.world)

	shoudDestroy := false
	for !shoudDestroy {
		select {
		case <-s.initChan:
			s.WaitForDeps()
			controller.Init(s.world)
			s.SendDone()
		case <-s.updateChan:
			s.WaitForDeps()
			controller.Update(s.world)
			s.SendDone()
		case <-s.fixedUpdateChan:
			s.WaitForDeps()
			controller.FixedUpdate(s.world)
			s.SendDone()
		case <-s.destroyChan:
			s.WaitForDeps()
			shoudDestroy = true
			s.SendDone()
		}
		s.world.wg.Done()
	}
}

func (s *SystemServiceInstance) asyncInit() {
	s.initChan <- struct{}{}
}

func (s *SystemServiceInstance) asyncUpdate() {
	s.updateChan <- struct{}{}
}

func (s *SystemServiceInstance) asyncFixedUpdate() {
	s.fixedUpdateChan <- struct{}{}
}

func (s *SystemServiceInstance) asyncDestroy() {
	s.destroyChan <- struct{}{}
}

type SystemService[T AnySystemControllerPtr] struct {
	initValue T
	dependsOn []AnySystemServicePtr
	instances map[*World]*SystemServiceInstance
}

func (s *SystemService[T]) register(world *World) *SystemServiceInstance {
	s.instances[world] = &SystemServiceInstance{
		controller:      s.initValue,
		dependsOnNum:    len(s.dependsOn),
		dependencyWg:    new(sync.WaitGroup),
		world:           world,
		initChan:        make(chan struct{}),
		updateChan:      make(chan struct{}),
		fixedUpdateChan: make(chan struct{}),
		destroyChan:     make(chan struct{}),
	}

	for i := range s.dependsOn {
		s.dependsOn[i].registerDependencyFor(world, s)
	}

	return s.instances[world]
}

func (s *SystemService[T]) registerDependencyFor(world *World, dep AnySystemServicePtr) {
	instance := s.instances[world]
	instance.dependencyFor = append(instance.dependencyFor, dep.getInstance(world))
}

func (s *SystemService[T]) getInstance(world *World) *SystemServiceInstance {
	return s.instances[world]
}

// TODO: dependsOn
func CreateSystem[T AnySystemControllerPtr](controller T, dependsOn ...AnySystemServicePtr) SystemService[T] {
	return SystemService[T]{
		initValue: controller,
		dependsOn: dependsOn,
		instances: make(map[*World]*SystemServiceInstance),
	}
}
