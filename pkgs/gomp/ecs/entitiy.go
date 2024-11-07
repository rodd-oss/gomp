package ecs

import "github.com/yohamta/donburi"

type Entity func(world donburi.World, extra ...Component)
