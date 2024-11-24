/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type Mask uint64
type EntityID int
type ComponentID uint64
type ComponentInstanceID uint64

var World world = world{
	Entities: make(map[EntityID]Entity),

	nextEntityID:    0,
	nextComponentID: 0,
}

type world struct {
	Entities map[EntityID]Entity

	nextEntityID    EntityID
	nextComponentID ComponentID
}

func (w *world) Create(components ...Component[any]) EntityID {
	entity := Entity{}
	entity.ID = w.nextEntityID
	w.Entities[w.nextEntityID] = entity
	w.nextEntityID++

	return entity.ID
}

func (w *world) Get(id EntityID) *Entity {
	ent, ok := w.Entities[id]
	if !ok {
		return nil
	}

	return &ent
}

func (w *world) Destroy(id EntityID) {
	delete(w.Entities, id)
}

type Entity struct {
	ID                 EntityID
	ComponentsMask     []Mask
	ComponentInstances []*ComponentInstanceID
}

func CreateComponent[T any](data T) Component[T] {
	component := Component[T]{}
	component.ID = World.nextComponentID

	World.nextComponentID++

	return component
}

type Component[T any] struct {
	ID        ComponentID
	Instances []ComponentInstance[T]
}

func (e *Component[T]) Add(id EntityID, data T) {}

func (c *Component[T]) New(initialValue T) ComponentInstance[T] {
	id := ComponentInstanceID(len(c.Instances))
	instance := ComponentInstance[T]{}
	instance.ID = id
	instance.Data = initialValue

	c.Instances = append(c.Instances, instance)

	return instance
}

func (c *Component[T]) Get(id EntityID) {
	// World.Entities[id].ComponentInstances
}

type ComponentInstance[T any] struct {
	ID       ComponentInstanceID
	EntityID EntityID
	Data     T
}

type Transform struct {
	x, y, z float32
}

func main() {
	var transformComponent = NewSparseSet[Transform, EntityID]()

	transformComponent.Add(10, Transform{1, 2, 3})
	t := transformComponent.Get(10)
	t.x = 0
	transformComponent.Delete(10)
}
