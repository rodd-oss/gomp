package ecs

import "testing"

const testEntitiesLen Entity = 100_000
const maxComponentsLen = 1024

func BenchmarkComponentBitTable_SetAndTest(b *testing.B) {
	// using a fixed maximum components length
	// create and extend the table
	table := NewComponentBitTable(maxComponentsLen)
	for i := Entity(0); i < testEntitiesLen; i++ {
		comp := ComponentId(i % maxComponentsLen)
		table.Create(i)
		table.Set(i, comp)
		if !table.Test(i, comp) {
			b.Fatalf("BitTable: expected entity %d to have component %d set", i, comp)
		}
	}
	for i := Entity(0); i < testEntitiesLen; i++ {
		table.Delete(i)
	}

	b.ReportAllocs()
	for b.Loop() {
		for i := Entity(0); i < testEntitiesLen; i++ {
			comp := ComponentId(i % maxComponentsLen)
			table.Create(i)
			table.Set(i, comp)
			if !table.Test(i, comp) {
				b.Fatalf("BitTable: expected entity %d to have component %d set", i, comp)
			}
		}
		b.StopTimer()
		for i := Entity(0); i < testEntitiesLen; i++ {
			table.Delete(i)
		}
		b.StartTimer()
	}
}

func BenchmarkComponentBitTable_SetTestDelete(b *testing.B) {
	table := NewComponentBitTable(maxComponentsLen)
	//for i := Entity(1); i < testEntitiesLen; i++ {
	//	table.Create(Entity(i))
	//}
	// using a fixed maximum components length
	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		for i := Entity(0); i < testEntitiesLen; i++ {
			table.Create(i)
		}
		b.StartTimer()

		for i := Entity(0); i < testEntitiesLen; i++ {
			comp := ComponentId(i % maxComponentsLen)
			table.Set(i, comp)
			if !table.Test(i, comp) {
				b.Fatalf("BitTable: expected entity %d to have component %d set", i, comp)
			}
			table.Delete(i)
		}
	}
}

func BenchmarkComponentBitTable_Delete(b *testing.B) {
	// using a fixed maximum components length
	table := NewComponentBitTable(maxComponentsLen)
	// Setup - create and set components for all entities
	for i := Entity(0); i < testEntitiesLen; i++ {
		entity := Entity(i)
		table.Create(entity)
		table.Set(entity, ComponentId(i%maxComponentsLen))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		// Delete all entities
		for i := Entity(0); i < testEntitiesLen; i++ {
			entity := Entity(i)
			table.Delete(entity)
		}

		// Recreate for next iteration
		for i := Entity(0); i < testEntitiesLen; i++ {
			entity := Entity(i)
			table.Create(entity)
			table.Set(entity, ComponentId(i%maxComponentsLen))
		}
	}
}

func BenchmarkComponentBitTable_AllSet(b *testing.B) {
	table := NewComponentBitTable(maxComponentsLen)
	// Prepare entities with multiple components each
	for i := Entity(0); i < testEntitiesLen; i++ {
		entity := Entity(i)
		table.Create(entity)
		// Give each entity 100 components
		for j := 0; j < 100; j++ {
			table.Set(entity, ComponentId(j%maxComponentsLen))
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		for i := Entity(0); i < testEntitiesLen; i++ {
			entity := Entity(i)
			count := 0
			table.AllSet(entity, func(id ComponentId) bool {
				count++
				return true
			})
			if count != 100 {
				b.Fatalf("Expected 3 components, got %d", count)
			}
		}
	}
}
