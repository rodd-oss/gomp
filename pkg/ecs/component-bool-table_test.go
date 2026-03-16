package ecs

import "testing"

func TestComponentByteTable_SetAndTest(t *testing.T) {
	// ...existing code...
	table := NewComponentBoolTable(10)
	entity := Entity(1)
	table.Set(entity, ComponentId(3))
	if !table.Test(entity, ComponentId(3)) {
		t.Errorf("Expected component 3 to be set for entity %d", entity)
	}
}

func TestComponentByteTable_Unset(t *testing.T) {
	// ...existing code...
	table := NewComponentBoolTable(10)
	entity := Entity(2)
	table.Set(entity, ComponentId(5))
	if !table.Test(entity, ComponentId(5)) {
		t.Errorf("Expected component 5 to be set for entity %d", entity)
	}
	table.Unset(entity, ComponentId(5))
	if table.Test(entity, ComponentId(5)) {
		t.Errorf("Expected component 5 to be unset for entity %d", entity)
	}
}

func TestComponentByteTable_AllSet(t *testing.T) {
	// ...existing code...
	table := NewComponentBoolTable(10)
	entity := Entity(3)
	components := []ComponentId{2, 4, 7}
	for _, id := range components {
		table.Set(entity, id)
	}

	var got []ComponentId
	table.AllSet(entity, func(id ComponentId) bool {
		got = append(got, id)
		return true
	})

	if len(got) != len(components) {
		t.Errorf("Expected %d components, got %d", len(components), len(got))
	}
}

func TestComponentByteTable_MultipleEntities(t *testing.T) {
	// ...existing code...
	table := NewComponentBoolTable(10)
	for i := 1; i <= 5; i++ {
		e := Entity(i)
		table.Set(e, ComponentId(i))
	}
	for i := 1; i <= 5; i++ {
		e := Entity(i)
		if !table.Test(e, ComponentId(i)) {
			t.Errorf("Entity %d should have component %d set", e, i)
		}
	}
}
