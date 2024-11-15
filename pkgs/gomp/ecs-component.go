package gomp

import "github.com/yohamta/donburi"

type IComponent = donburi.IComponentType

type Component struct {
	ComponentType IComponent
	Set           func(*donburi.Entry)
}
