/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package ecs

import (
	"math/bits"
	"testing"
)

func TestNewComponentBitTable(t *testing.T) {
	// Test with different max component sizes
	tests := []struct {
		name               string
		maxComponentsLen   int
		expectedBitsetSize int
	}{
		{"Small", 10, 1},
		{"Medium", 64, 1},
		{"Large", 192, 3},
		{"Below", 172, 3},
		{"Above", 200, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewComponentBitTable(tt.maxComponentsLen)

			if cap(table.bitsetsBook) != initialBookSize {
				t.Errorf("Expected %d preallocated chunks, got %d", initialBookSize, cap(table.bitsetsBook))
			}
		})
	}
}

func TestComponentBitTable_Set(t *testing.T) {
	table := NewComponentBitTable(100)

	// Set a bit for a new entity
	entity1 := Entity(1)
	table.Create(entity1)
	table.Set(entity1, ComponentId(5))

	// Verify bit was set
	bitsId, ok := table.lookup.Get(entity1)
	if !ok {
		t.Fatalf("Entity %d not found in lookup", entity1)
	}

	pageId, bitsetId := table.getPageIDAndBitsetIndex(bitsId)
	offset := int(ComponentId(5) >> uintShift)
	mask := uint(1 << (ComponentId(5) % bits.UintSize))

	if (table.bitsetsBook[pageId][bitsetId+offset] & mask) == 0 {
		t.Errorf("Expected bit to be set for entity %d, component %d", entity1, 5)
	}

	// Set multiple bits for same entity
	table.Set(entity1, ComponentId(10))
	table.Set(entity1, ComponentId(63))
	table.Set(entity1, ComponentId(64))

	// Set bits for a different entity
	entity2 := Entity(2)
	table.Create(entity2)
	table.Set(entity2, ComponentId(5))
}

func TestComponentBitTable_Unset(t *testing.T) {
	table := NewComponentBitTable(100)
	entity := Entity(1)

	table.Create(entity)
	// Set and then unset
	table.Set(entity, ComponentId(5))
	table.Set(entity, ComponentId(10))
	table.Unset(entity, ComponentId(5))

	// Verify bit was unset
	bitsId, _ := table.lookup.Get(entity)
	pageId, bitsetId := table.getPageIDAndBitsetIndex(bitsId)
	offset := int(ComponentId(5) >> uintShift)
	mask := uint(1 << (ComponentId(5) % bits.UintSize))

	if (table.bitsetsBook[pageId][bitsetId+offset] & mask) != 0 {
		t.Errorf("Expected bit to be unset for entity %d, component %d", entity, 5)
	}

	// Verify other bit is still set
	offset = int(ComponentId(10) >> uintShift)
	mask = uint(1 << (ComponentId(10) % bits.UintSize))

	if (table.bitsetsBook[pageId][bitsetId+offset] & mask) == 0 {
		t.Errorf("Expected bit to still be set for entity %d, component %d", entity, 10)
	}
}

func TestComponentBitTable_Test(t *testing.T) {
	table := NewComponentBitTable(100)

	// Test for non-existent entity
	if table.Test(Entity(999), ComponentId(5)) {
		t.Error("Test should return false for non-existent entity")
	}

	// Set up an entity with some components
	entity := Entity(42)
	table.Create(entity)
	table.Set(entity, ComponentId(5))
	table.Set(entity, ComponentId(64))

	// Test for set components
	if !table.Test(entity, ComponentId(5)) {
		t.Error("Test should return true for set component 5")
	}
	if !table.Test(entity, ComponentId(64)) {
		t.Error("Test should return true for set component 64")
	}

	// Test for unset component
	if table.Test(entity, ComponentId(10)) {
		t.Error("Test should return false for unset component 10")
	}

	// Test after unsetting a component
	table.Unset(entity, ComponentId(5))
	if table.Test(entity, ComponentId(5)) {
		t.Error("Test should return false after component is unset")
	}

	// Test components at boundaries
	entity2 := Entity(43)
	table.Create(entity2)
	table.Set(entity2, ComponentId(0))
	table.Set(entity2, ComponentId(64)) // First bit in second uint
	table.Set(entity2, ComponentId(65)) // Second bit in second uint

	if !table.Test(entity2, ComponentId(0)) {
		t.Error("Test should return true for set component at uint boundary (0)")
	}
	if !table.Test(entity2, ComponentId(64)) {
		t.Error("Test should return true for set component at uint boundary (64)")
	}
	if !table.Test(entity2, ComponentId(65)) {
		t.Error("Test should return true for set component at uint boundary (65)")
	}
}

func TestComponentBitTable_AllSet(t *testing.T) {
	table := NewComponentBitTable(200)
	entity := Entity(1)
	table.Create(entity)

	// Set several components
	expectedComponents := []ComponentId{5, 10, 64, 128, 199}
	for _, id := range expectedComponents {
		table.Set(entity, id)
	}

	// Use AllSet to collect components
	var foundComponents []ComponentId
	table.AllSet(entity, func(id ComponentId) bool {
		foundComponents = append(foundComponents, id)
		return true
	})

	// Verify all components were found
	if len(foundComponents) != len(expectedComponents) {
		t.Errorf("Expected %d components, found %d", len(expectedComponents), len(foundComponents))
	}

	// Check each component
	componentMap := make(map[ComponentId]bool)
	for _, id := range foundComponents {
		componentMap[id] = true
	}

	for _, id := range expectedComponents {
		if !componentMap[id] {
			t.Errorf("Component %d not found in AllSet results", id)
		}
	}

	// Test early termination
	count := 0
	table.AllSet(entity, func(id ComponentId) bool {
		count++
		return count < 3 // Stop after finding 2 components
	})

	if count != 3 {
		t.Errorf("Early termination didn't work as expected. Count: %d", count)
	}
}

func TestComponentBitTable_extend(t *testing.T) {
	// Create a table with a small chunk size for testing
	table := NewComponentBitTable(bits.UintSize)

	// Create entities to fill the first chunk
	for i := 0; i < pageSize; i++ {
		e := Entity(i)
		table.Create(e)
		table.Set(e, ComponentId(i%bits.UintSize))
	}
	if len(table.bitsetsBook) != 1 {
		t.Errorf("Expected table to extend, got %d chunks", len(table.bitsetsBook))
	}

	// Create entity to force extension
	table.Create(pageSize)
	if len(table.bitsetsBook) <= 1 {
		t.Errorf("Expected table to be extended, got %d chunks", len(table.bitsetsBook))
	}

	// Create more entities
	for i := pageSize + 1; i < (pageSize*2)+1; i++ {
		e := Entity(i)
		table.Create(e)
		table.Set(e, ComponentId(1))
	}
	if len(table.bitsetsBook) <= 2 {
		t.Errorf("Expected table to extend beyond 2 chunks, got %d chunks", len(table.bitsetsBook))
	}
}

func TestComponentBitTable_EdgeCases(t *testing.T) {
	table := NewComponentBitTable(100)

	// Test AllSet on non-existent entity
	called := false
	table.AllSet(Entity(999), func(id ComponentId) bool {
		called = true
		return true
	})

	if called {
		t.Errorf("AllSet callback should not be called for non-existent entity")
	}

	// Test setting across bit boundaries
	entity := Entity(42)
	table.Create(entity)
	table.Set(entity, ComponentId(0))  // First bit in first uint
	table.Set(entity, ComponentId(63)) // Last bit in first uint
	table.Set(entity, ComponentId(64)) // First bit in second uint

	// Verify all bits
	var found []ComponentId
	table.AllSet(entity, func(id ComponentId) bool {
		found = append(found, id)
		return true
	})

	if len(found) != 3 || !contains(found, ComponentId(0)) || !contains(found, ComponentId(63)) || !contains(found, ComponentId(64)) {
		t.Errorf("Expected components 0, 63 and 64, got %v", found)
	}
}

func TestComponentBitTable_Create(t *testing.T) {
	table := NewComponentBitTable(100)
	entity := Entity(42)

	// Create entity
	table.Create(entity)

	// Verify entity is in lookup
	_, ok := table.lookup.Get(entity)
	if !ok {
		t.Fatalf("Entity %d not found in lookup after Create", entity)
	}

	// Check that entity ID is stored in entities book
	bitsId, _ := table.lookup.Get(entity)
	pageId, entityId := table.getPageIDAndEntityIndex(bitsId)

	if table.entitiesBook[pageId][entityId] != entity {
		t.Errorf("Expected entity ID %d stored in entities book, got %d",
			entity, table.entitiesBook[pageId][entityId])
	}
}

func TestComponentBitTable_Delete(t *testing.T) {
	table := NewComponentBitTable(100)
	entity := Entity(42)
	table.Create(entity)

	// Set multiple components
	table.Set(entity, ComponentId(5))
	table.Set(entity, ComponentId(10))

	// Ensure components are set
	if !table.Test(entity, ComponentId(5)) || !table.Test(entity, ComponentId(10)) {
		t.Fatalf("Expected components to be set for entity %d", entity)
	}

	// Delete the entity
	table.Delete(entity)

	// Verify entity is no longer in lookup
	_, ok := table.lookup.Get(entity)
	if ok {
		t.Errorf("Entity %d should be removed from lookup after deletion", entity)
	}

	// Test should return false for deleted entity
	if table.Test(entity, ComponentId(5)) || table.Test(entity, ComponentId(10)) {
		t.Errorf("Test should return false for deleted entity %d", entity)
	}
}

func TestComponentBitTable_DeleteWithSwap(t *testing.T) {
	table := NewComponentBitTable(100)

	// Create two entities
	entity1 := Entity(1)
	entity2 := Entity(2)
	table.Create(entity1)
	table.Create(entity2)

	// Set different components for each
	table.Set(entity1, ComponentId(5))
	table.Set(entity1, ComponentId(10))
	table.Set(entity2, ComponentId(15))
	table.Set(entity2, ComponentId(20))

	// Get entity2's bits ID before deletion of entity1
	entity2BitsId, _ := table.lookup.Get(entity2)

	// Delete the first entity - should swap with entity2
	table.Delete(entity1)

	// Verify entity1 is gone
	_, ok := table.lookup.Get(entity1)
	if ok {
		t.Errorf("Entity %d should be removed from lookup", entity1)
	}

	// Verify entity2's data is still accessible
	if !table.Test(entity2, ComponentId(15)) || !table.Test(entity2, ComponentId(20)) {
		t.Errorf("Entity %d should still have its components after swap", entity2)
	}

	// Entity2's lookup entry should now point to entity1's old position
	newEntity2BitsId, ok := table.lookup.Get(entity2)
	if !ok {
		t.Fatalf("Entity %d not found after swap", entity2)
	}

	// entity2 should have been moved to entity1's position
	if newEntity2BitsId == entity2BitsId {
		t.Errorf("Entity %d position should have changed after swap", entity2)
	}

	// Ensure entity is stored in the entity book
	pageId, entityId := table.getPageIDAndEntityIndex(newEntity2BitsId)
	if table.entitiesBook[pageId][entityId] != entity2 {
		t.Errorf("Entity ID %d not correctly stored after swap, got %d",
			entity2, table.entitiesBook[pageId][entityId])
	}
}

func TestComponentBitTable_MultipleOperations(t *testing.T) {
	table := NewComponentBitTable(100)

	// Create, set, delete several entities in sequence
	for i := 1; i <= 5; i++ {
		entity := Entity(i)
		table.Create(entity)
		table.Set(entity, ComponentId(i))
		table.Set(entity, ComponentId(i+10))
	}

	// Delete entity 2 and 4
	table.Delete(Entity(2))
	table.Delete(Entity(4))

	// Verify entities 1, 3, 5 still exist with correct components
	for _, id := range []Entity{1, 3, 5} {
		if !table.Test(id, ComponentId(int(id))) || !table.Test(id, ComponentId(int(id)+10)) {
			t.Errorf("Entity %d should still have its components", id)
		}
	}

	// Verify entities 2 and 4 are gone
	for _, id := range []Entity{2, 4} {
		if _, ok := table.lookup.Get(id); ok {
			t.Errorf("Entity %d should have been deleted", id)
		}
	}

	// Create new entities
	for i := 6; i <= 7; i++ {
		entity := Entity(i)
		table.Create(entity)
		table.Set(entity, ComponentId(i))
	}

	// Verify new entities have correct components
	for i := 6; i <= 7; i++ {
		entity := Entity(i)
		if !table.Test(entity, ComponentId(i)) {
			t.Errorf("Entity %d should have component %d set", entity, i)
		}
	}

	// The length should be 5 (entities 1, 3, 5, 6, 7)
	if table.length != 5 {
		t.Errorf("Expected length 5, got %d", table.length)
	}
}

// Helper function to check if slice contains a value
func contains(s []ComponentId, id ComponentId) bool {
	for _, v := range s {
		if v == id {
			return true
		}
	}
	return false
}
