package ecs

import (
	"sync"
	"testing"
)

func TestLookupMap(t *testing.T) {
	lm := &LookupMap[int]{}

	// Test setting and getting a value
	key := uint32(0x12345678)
	value := 42
	lm.Set(key, value)

	// Test getting the value
	ret, ok := lm.Get(key)
	if !ok {
		t.Errorf("Expected to find value for key %x, but got not found", key)
	}
	if ret != value {
		t.Errorf("Expected value %d for key %x, but got %d", value, key, ret)
	}

	// Test getting a non-existent key
	nonExistentKey := uint32(0x87654321)
	ret, ok = lm.Get(nonExistentKey)
	if ok {
		t.Errorf("Expected not to find value for key %x, but got found", nonExistentKey)
	}
	if ret != 0 {
		t.Errorf("Expected zero value for non-existent key %x, but got %d", nonExistentKey, ret)
	}

	// Test overwriting a value
	newValue := 100
	lm.Set(key, newValue)
	ret, ok = lm.Get(key)
	if !ok {
		t.Errorf("Expected to find value for key %x after overwrite, but got not found", key)
	}
	if ret != newValue {
		t.Errorf("Expected value %d for key %x after overwrite, but got %d", newValue, key, ret)
	}
}

func TestLookupMapConcurrency(t *testing.T) {
	lm := &LookupMap[int]{}
	key := uint32(0x12345678)
	value := 42

	// Test concurrent writes
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lm.Set(key, value+i)
		}(i)
	}
	wg.Wait()

	// Test concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ret, ok := lm.Get(key)
			if !ok {
				t.Errorf("Expected to find value for key %x, but got not found", key)
			}
			if ret < value || ret >= value+100 {
				t.Errorf("Expected value between %d and %d for key %x, but got %d", value, value+99, key, ret)
			}
		}()
	}
	wg.Wait()
}
