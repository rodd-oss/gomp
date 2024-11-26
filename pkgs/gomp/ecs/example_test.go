/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import (
	"log"
	"testing"
	"time"
)

type Transform struct {
	X, Y, Z float32
}

var transformComponent = CreateComponent[Transform]()

type Rotation struct {
	RX, RY, RZ int
}

var _ = CreateComponent[Rotation]()

type Scale struct {
	Value float32
}

var scaleComponent = CreateComponent[Scale]()

func TestExample(t *testing.T) {
	var world = New("Main")

	world.RegisterComponents(
		&scaleComponent,
		&transformComponent,
	)

	tra := Transform{0, 1, 2}

	start := time.Now()
	for i := 0; i < 10000000; i++ {
		v := world.CreateEntity("Player")
		transformComponent.Set(v, tra)
	}
	middle := time.Since(start)
	log.Println(middle)

	arr := transformComponent.Instances[&world].dense
	l := len(arr)
	for i := 0; i < l; i++ {
		arr[i].X += 1
		arr[i].Y = 0
		arr[i].Z += 2
	}
	log.Println(time.Since(start) - middle)
}
