package qsm

import (
	"fmt"
	"sync"
	"time"
)

type State comparable

type QSM[T State] struct {
	Mx sync.Mutex

	currentState    *T
	currentMutation *Mutation[T]
	mutations       map[T]map[T]*MutationRule
}

func Init[T State]() *QSM[T] {
	qsm := &QSM[T]{
		currentState: nil,
		currentMutation: &Mutation[T]{
			targetState: nil,
			progress:    0.0,
			rule:        &MutationRule{},
		},
		mutations: make(map[T]map[T]*MutationRule),
	}

	return qsm
}

func (fsm *QSM[T]) Start(initialState T) {
	fsm.currentState = &initialState
}

func (fsm *QSM[T]) SetMutationRuleNto1(from []T, to T, rule *MutationRule) *QSM[T] {
	if rule.Duration < 0 {
		panic("Time travel is not available in this patch. Please come back a few updates later.")
	}

	for _, f := range from {
		if fsm.mutations[f] == nil {
			fsm.mutations[f] = make(map[T]*MutationRule)
		}

		fsm.mutations[f][to] = rule
	}

	return fsm
}

func (fsm *QSM[T]) SetMutationRule1toN(from T, to []T, rule *MutationRule) *QSM[T] {
	if rule.Duration < 0 {
		panic("Time travel is not available in this patch. Please come back a few updates later.")
	}

	for _, t := range to {
		if fsm.mutations[from] == nil {
			fsm.mutations[from] = make(map[T]*MutationRule)
		}
		fsm.mutations[from][t] = rule
	}

	return fsm
}

func (fsm *QSM[T]) SetMutationRuleNtoN(from []T, to []T, rule *MutationRule) *QSM[T] {
	if rule.Duration < 0 {
		panic("Time travel is not available in this patch. Please come back a few updates later.")
	}

	for _, f := range from {
		for _, t := range to {
			if f == t {
				continue
			}

			if fsm.mutations[f] == nil {
				fsm.mutations[f] = make(map[T]*MutationRule)
			}

			fsm.mutations[f][t] = rule
		}
	}

	return fsm
}

func (fsm *QSM[T]) SetMutationRule(from T, to T, rule *MutationRule) *QSM[T] {
	if rule.Duration < 0 {
		panic("Time travel is not available in this patch. Please come back a few updates later.")
	}

	if fsm.mutations[from] == nil {
		fsm.mutations[from] = make(map[T]*MutationRule)
	}

	fsm.mutations[from][to] = rule

	return fsm
}

func (fsm *QSM[T]) RemoveMutationRule(from T, to T) {
	delete(fsm.mutations[from], to)
}

func (fsm *QSM[T]) CancelMutation() {
	fsm.currentMutation.rule.runCancel()
	fsm.currentMutation.reset()
}

func (fsm *QSM[T]) Update(dt time.Duration) {
	fsm.Mx.Lock()
	defer fsm.Mx.Unlock()

	if fsm.currentMutation == nil {
		return
	}

	rule := fsm.currentMutation.rule
	if rule == nil {
		return
	}

	rule.runWhile(dt)

	// fmt.Printf("from %v to %v -- %0.1f \n", *fsm.currentState, *fsm.currentMutation.targetState, fsm.currentMutation.progress)

	if fsm.currentMutation.progress < 1.0 {
		fsm.currentMutation.progress += float32(dt) / float32(rule.Duration)
	} else {
		fsm.currentState = fsm.currentMutation.targetState
		rule.runAfter()
		fsm.currentMutation.reset()
	}
}

func (fsm *QSM[T]) Mutate(toState T) error {
	if fsm.currentState == nil {
		panic("No current state. Call Start() first.")
	}

	if fsm.currentMutation != nil {
		fsm.CancelMutation()
	}

	fsm.currentMutation.rule = fsm.mutations[*fsm.currentState][toState]
	if fsm.currentMutation.rule == nil {
		return fmt.Errorf("No mutation rule for %v -> %v", *fsm.currentState, toState)
	}

	fsm.currentMutation.rule.runBefore()
	fsm.currentMutation.targetState = &toState

	return nil
}
