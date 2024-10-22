package components

import (
	"tomb_mates/internal/protos"

	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
	"google.golang.org/protobuf/proto"
)

type NetworkAreaData struct {
	Area *protos.Area
}

var NetworkArea = ecs.NewComponentType[NetworkAreaData]()

type NetworkManagerData struct {
	MyID *uint32

	UnhandledEvents []*protos.Event
	NetworkEntities map[uint32]*NetworkEntityData

	// Sync game state
	IncomingPatch       *protos.GameStatePatch // Clientside
	OutgoingPatch       *protos.GameStatePatch // Serverside
	NetworkIdToEntityId map[uint32]ecs.Entity
	EntityIdToNetworkId map[ecs.Entity]uint32

	LastNetworkEntityId uint32

	Broadcast chan []byte
	World     ecs.World
	Space     *cp.Space
}

func (nm *NetworkManagerData) Register(entityId ecs.Entity, id uint32) *protos.NetworkEntity {
	entity := nm.World.Entry(entityId)

	transform := Transform.GetValue(entity)
	physics := Physics.GetValue(entity)

	ned := &NetworkEntityData{
		Id:        id,
		LastPatch: nil,
		Transform: &transform,
		Body:      physics.Body,
	}

	NetworkEntity.SetValue(entity, ned)

	if nm.OutgoingPatch == nil {
		nm.OutgoingPatch = &protos.GameStatePatch{}
	}

	// Prevent from being in outgoing patch
	if nm.OutgoingPatch.DeletedEntities != nil {
		delete(nm.OutgoingPatch.DeletedEntities, id)
	}

	if nm.OutgoingPatch.CreatedEntities == nil {
		nm.OutgoingPatch.CreatedEntities = make(map[uint32]*protos.NetworkEntity)
	}

	velocity := physics.Body.Velocity()

	nm.OutgoingPatch.CreatedEntities[id] = &protos.NetworkEntity{
		Id: ned.Id,
		Transform: &protos.Transform{
			Position: &protos.Vector2{
				X: transform.LocalPosition.X,
				Y: transform.LocalPosition.Y,
			},
			Rotation: transform.LocalRotation,
			Scale: &protos.Vector2{
				X: transform.LocalScale.X,
				Y: transform.LocalScale.Y,
			},
		},
		Physics: &protos.Physics{
			Velocity: &protos.Vector2{
				X: velocity.X,
				Y: velocity.Y,
			},
		},
	}

	// Register new entities
	nm.NetworkIdToEntityId[id] = entityId
	nm.EntityIdToNetworkId[entityId] = id
	nm.NetworkEntities[id] = ned

	return nm.OutgoingPatch.CreatedEntities[id]
}

func (nm *NetworkManagerData) Deregister(entityId ecs.Entity) {
	id := nm.EntityIdToNetworkId[entityId]

	if nm.OutgoingPatch == nil {
		nm.OutgoingPatch = &protos.GameStatePatch{}
	}

	delete(nm.OutgoingPatch.CreatedEntities, id)

	// Prevent from being in outgoing patch

	if nm.OutgoingPatch.DeletedEntities == nil {
		nm.OutgoingPatch.DeletedEntities = make(map[uint32]*protos.Empty)
	}

	nm.OutgoingPatch.DeletedEntities[id] = &protos.Empty{}

	delete(nm.NetworkEntities, id)
	delete(nm.NetworkIdToEntityId, id)
	delete(nm.EntityIdToNetworkId, entityId)
}

func (nm *NetworkManagerData) Update(dt float64, isClient bool) {
	// Patch existing entities
	if nm.IncomingPatch != nil {
		for id, patch := range nm.IncomingPatch.Entities {
			if patch == nil {
				continue
			}

			entId := nm.NetworkIdToEntityId[id]
			ent := nm.World.Entry(entId)
			if ent == nil {
				continue
			}

			ne := NetworkEntity.GetValue(ent)
			if ne == nil {
				continue
			}

			ne.ApplyPatch(patch)
		}

		nm.IncomingPatch = nil
	}

	// // Create new entities -> TODO: move to global manager?
	// for _, entity := range nm.IncomingPatch.CreatedEntities {
	// 	if entity == nil {
	// 		continue
	// 	}

	// 	newEntityId := nm.World.Create(Transform, Physics, NetworkEntity)
	// 	newEntity := nm.World.Entry(newEntityId)

	// 	// TODO: this should be done in a entity type switcher
	// 	// switch entity.Type {
	// 	// case protos.EntityType_unit:
	// 	// 	nm.CreatedUnits[entity.Id] = entity.Unit
	// 	// case protos.EntityType_area:
	// 	// 	nm.CreatedAreas[entity.Id] = entity.Area
	// 	// }
	// 	Transform.SetValue(newEntity, TransformData{
	// 		LocalPosition: math.Vec2{X: 1, Y: 2},
	// 		LocalRotation: 0,
	// 		LocalScale: math.Vec2{
	// 			X: 1,
	// 			Y: 1,
	// 		},
	// 	})

	// 	body := cp.NewKinematicBody()
	// 	body = nm.Space.AddBody(body)
	// 	body.SetPosition(cp.Vector{
	// 		X: 1,
	// 		Y: 2,
	// 	})

	// 	shape := nm.Space.AddShape(cp.NewCircle(body, 8, cp.Vector{}))
	// 	shape.SetElasticity(0)
	// 	shape.SetFriction(0)

	// 	Physics.SetValue(newEntity, PhysicsData{
	// 		Body: body,
	// 	})

	// 	// For all entities in the patch
	// 	newNetworkId := nm.LastNetworkEntityId + 1
	// 	NetworkEntity.SetValue(newEntity, NetworkEntityData{
	// 		Id: newNetworkId,
	// 	})
	// 	nm.LastNetworkEntityId = newNetworkId
	// }

	// // Delete entities -> TODO: move to global manager?
	// for id, entity := range nm.IncomingPatch.DeletedEntities {
	// 	if entity == nil {
	// 		continue
	// 	}

	// 	entityId := nm.NetworkIdToEntity[id]
	// 	entityToDelete := nm.World.Entry(entityId)
	// 	if entityToDelete == nil {
	// 		continue
	// 	}

	// 	physics := Physics.GetValue(entityToDelete)
	// 	nm.Space.RemoveBody(physics.Body)
	// 	nm.World.Remove(entityId)
	// }
}

func (n *NetworkManagerData) SendPatch() {
	if n.OutgoingPatch == nil {
		n.OutgoingPatch = &protos.GameStatePatch{}
	}

	NetworkEntity.Each(n.World, func(ent *ecs.Entry) {
		ne := NetworkEntity.GetValue(ent)
		if ne == nil {
			return
		}

		patch := ne.RequestPatch(ent)
		if patch == nil {
			return
		}

		if n.OutgoingPatch.Entities == nil {
			n.OutgoingPatch.Entities = make(map[uint32]*protos.PatchNetworkEntity)
		}

		if n.OutgoingPatch.Entities[ne.Id] == nil {
			n.OutgoingPatch.Entities[ne.Id] = &protos.PatchNetworkEntity{}
		}

		if patch.Physics != nil {
			n.OutgoingPatch.Entities[ne.Id].Physics = patch.Physics
		}

		if patch.Transform != nil {
			n.OutgoingPatch.Entities[ne.Id].Transform = patch.Transform
		}

		ne.ApplyPatch(patch)
	})

	if n.OutgoingPatch == nil {
		return
	}

	if n.OutgoingPatch.Entities == nil &&
		n.OutgoingPatch.CreatedEntities == nil &&
		n.OutgoingPatch.DeletedEntities == nil {
		return
	}

	statePatchEvent := &protos.Event{
		Type: protos.EventType_state_patch,
		Data: &protos.Event_StatePatch{
			StatePatch: n.OutgoingPatch,
		},
	}

	message, err := proto.Marshal(statePatchEvent)
	if err != nil {
		return
	}

	n.Broadcast <- message

	n.OutgoingPatch = &protos.GameStatePatch{}
}
