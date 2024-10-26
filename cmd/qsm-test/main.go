/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"log"
	"time"
	"tomb_mates/internal/qsm"
)

type State int

const (
	idle State = iota
	walk
	death
)

func main() {
	playerStateMachine := qsm.Init[State]()

	playerStateMachine.SetMutationRuleNtoN([]State{idle, walk}, []State{idle, walk}, &qsm.MutationRule{
		Before: func() {
			log.Println("before idle or walk->death")
		},
		While: func(dt time.Duration) {
			log.Println("while idle or walk->death")
		},
		After: func() {
			log.Println("after idle or walk->death")
		},
		Cancel: func() {
			log.Println("cancelling idle or walk->death")
		},
		Duration: time.Second,
	})

	playerStateMachine.SetMutationRule(idle, walk, &qsm.MutationRule{
		Before: func() {
			log.Println("before idle->walk")
		},
		While: func(dt time.Duration) {
			log.Println("while idle->walk")
		},
		After: func() {
			log.Println("after idle->walk")
		},
		Duration: time.Second / 5,
	})

	playerStateMachine.SetMutationRule(walk, idle, &qsm.MutationRule{
		Before: func() {
			log.Println("before walk->idle")
		},
		While: func(dt time.Duration) {
			log.Println("while walk->idle")
		},
		After: func() {
			log.Println("after walk->idle")
		},
		Duration: time.Second / 2,
	})

	playerStateMachine.Start(idle)

	go func() {
		dt := time.Second / 10
		ticker := time.NewTicker(dt)
		for {
			select {
			case <-ticker.C:
				playerStateMachine.Update(dt)
			}
		}
	}()

	var err error

	err = playerStateMachine.Mutate(walk)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(2 * time.Second)

	err = playerStateMachine.Mutate(idle)
	if err != nil {
		log.Println(err)
	}

	err = playerStateMachine.Mutate(death)
	if err != nil {
		log.Println(err)
	}

	err = playerStateMachine.Mutate(walk)
	if err != nil {
		log.Println(err)
	}
	err = playerStateMachine.Mutate(death)
	if err != nil {
		log.Println(err)
	}

	time.Sleep(10 * time.Second)
}
