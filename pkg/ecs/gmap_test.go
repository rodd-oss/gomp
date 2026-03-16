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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenMap(t *testing.T) {
	t.Run("NewGenMap", func(t *testing.T) {
		m := NewGenMap[string, int](10)
		assert.Equal(t, m.Len(), 0)
	})

	t.Run("Set/Get/Has", func(t *testing.T) {
		m := NewGenMap[string, int](0)

		// Проверка отсутствующего ключа
		_, ok := m.Get("missing")
		assert.False(t, ok)
		assert.False(t, m.Has("missing"))

		// Добавление нового ключа
		m.Set("a", 1)
		val, ok := m.Get("a")
		assert.True(t, ok)
		assert.Equal(t, val, 1)
		assert.True(t, m.Has("a"))

		// Обновление существующего ключа
		m.Set("a", 42)
		val, _ = m.Get("a")
		assert.Equal(t, val, 42)
	})

	t.Run("Reset", func(t *testing.T) {
		m := NewGenMap[string, int](0)
		m.Set("a", 1)
		m.Set("b", 2)

		// Сброс поколения
		m.Reset()

		// Ключи не должны быть доступны
		assert.False(t, m.Has("a"))
		assert.False(t, m.Has("b"))
		assert.Equal(t, m.Len(), 0)

		// Добавление новых ключей после сброса
		m.Set("c", 3)
		assert.True(t, m.Has("c"))
		assert.Equal(t, m.Len(), 1)
	})

	t.Run("Delete", func(t *testing.T) {
		m := NewGenMap[string, int](0)
		m.Set("a", 1)

		// Удаление существующего ключа
		m.Delete("a")
		assert.False(t, m.Has("a"))

		// Повторное добавление после удаления
		m.Set("a", 2)
		assert.True(t, m.Has("a"))
		val, _ := m.Get("a")
		assert.Equal(t, val, 2)
	})

	t.Run("Len", func(t *testing.T) {
		m := NewGenMap[string, int](0)
		assert.Equal(t, m.Len(), 0)

		m.Set("a", 1)
		m.Set("b", 2)
		assert.Equal(t, m.Len(), 2)

		m.Reset()
		assert.Equal(t, m.Len(), 0)
	})

	t.Run("Each", func(t *testing.T) {
		m := NewGenMap[string, int](0)
		m.Set("a", 1)
		m.Set("b", 2)

		// Проверка итерации
		count := 0
		for k, v := range m.Each() {
			count++
			assert.True(t, k == "a" || k == "b")
			assert.True(t, v == 1 || v == 2)
		}
		assert.Equal(t, count, 2)

		// После сброса итерация пуста
		m.Reset()
		count = 0
		for range m.Each() {
			count++
		}
		assert.Equal(t, count, 0)
	})

	t.Run("Clear", func(t *testing.T) {
		m := NewGenMap[string, int](0)
		m.Set("a", 1)
		m.Reset()
		m.Set("b", 2)

		// Очистка старых поколений
		m.Clear()

		// Проверка, что старые записи удалены
		assert.Equal(t, len(m.data), 1)
		assert.True(t, m.Has("b"))
	})

	t.Run("MultipleGenerations", func(t *testing.T) {
		m := NewGenMap[string, int](0)

		// Generation 0
		m.Set("a", 1)
		m.Set("b", 2)

		// Generation 1
		m.Reset()
		m.Set("b", 20)
		m.Set("c", 30)

		// Проверка наличия ключей
		assert.False(t, m.Has("a")) // Из generation 0
		assert.True(t, m.Has("b"))  // Обновлено в generation 1
		assert.True(t, m.Has("c"))

		// Проверка значений
		val, _ := m.Get("b")
		assert.Equal(t, val, 20)
	})

	t.Run("ReuseKeyAfterReset", func(t *testing.T) {
		m := NewGenMap[string, int](0)
		m.Set("a", 1)
		m.Reset()
		m.Set("a", 2) // Тот же ключ, новое поколение

		assert.True(t, m.Has("a"))
		val, _ := m.Get("a")
		assert.Equal(t, val, 2)
	})
}
