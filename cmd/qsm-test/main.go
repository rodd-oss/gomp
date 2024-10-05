package main

import (
	"fmt"
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

	playerStateMachine.SetMutationRuleNto1([]State{idle, walk}, death, &qsm.MutationRule{
		Before: func() {
			fmt.Println("before idle or walk->death")
		},
		While: func(dt time.Duration) {
			fmt.Println("while idle or walk->death")
		},
		After: func() {
			fmt.Println("after idle or walk->death")
		},
		Cancel: func() {
			fmt.Println("cancelling idle or walk->death")
		},
		Duration: time.Second,
	})
	playerStateMachine.SetMutationRule(idle, walk, &qsm.MutationRule{
		Before: func() {
			fmt.Println("before idle->walk")
		},
		While: func(dt time.Duration) {
			fmt.Println("while idle->walk")
		},
		After: func() {
			fmt.Println("after idle->walk")
		},
		Duration: time.Second / 5,
	})

	playerStateMachine.SetMutationRule(walk, idle, &qsm.MutationRule{
		Before: func() {
			fmt.Println("before walk->idle")
		},
		While: func(dt time.Duration) {
			fmt.Println("while walk->idle")
		},
		After: func() {
			fmt.Println("after walk->idle")
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
		fmt.Println(err)
	}
	time.Sleep(2 * time.Second)

	err = playerStateMachine.Mutate(idle)
	if err != nil {
		fmt.Println(err)
	}

	err = playerStateMachine.Mutate(death)
	if err != nil {
		fmt.Println(err)
	}

	err = playerStateMachine.Mutate(walk)
	if err != nil {
		fmt.Println(err)
	}
	err = playerStateMachine.Mutate(death)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(10 * time.Second)
}
