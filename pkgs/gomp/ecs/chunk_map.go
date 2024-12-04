/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type ChunkMap[T any] struct {
	buffer           []ChunkMapElement[T]
	pageSize         int
	bufferInitialCap int
}

func NewChunkMap[T any](bufferCapacity int, pageSize int) (arr *ChunkMap[T]) {
	arr = new(ChunkMap[T])
	arr.bufferInitialCap = bufferCapacity
	arr.pageSize = pageSize
	arr.buffer = make([]ChunkMapElement[T], bufferCapacity)

	return arr
}

func (cm *ChunkMap[T]) Get(index int) (value T, ok bool) {
	pageId := index / cm.pageSize
	bufferLen := len(cm.buffer)
	if pageId >= bufferLen {
		return value, false
	}

	page := cm.buffer[pageId]

	index %= cm.pageSize
	if index >= len(page.data) {
		return value, false
	}
	data := page.data[index]

	if !data.exists {
		return value, false
	}

	return data.value, true
}

func (cm *ChunkMap[T]) Set(index int, value T) {
	var page *ChunkMapElement[T]

	pageId := index / cm.pageSize
	bufferLastIndex := len(cm.buffer) - 1
	if pageId > bufferLastIndex {
		delta := pageId - bufferLastIndex
		if delta < cm.bufferInitialCap {
			cm.buffer = append(cm.buffer, make([]ChunkMapElement[T], cm.bufferInitialCap)...)
			cm.bufferInitialCap += cm.bufferInitialCap
		} else {
			cm.buffer = append(cm.buffer, make([]ChunkMapElement[T], delta)...)
		}
	}
	page = &cm.buffer[pageId]

	index %= cm.pageSize
	if index >= len(page.data) {
		page.data = make([]ChunkMapElementData[T], cm.pageSize)
	}

	data := &page.data[index]
	data.exists = true
	data.value = value
}

func (cm *ChunkMap[T]) Delete(index int) {
	var page *ChunkMapElement[T]

	pageId := index / cm.pageSize
	if pageId >= len(cm.buffer)-1 {
		return
	}

	page = &cm.buffer[pageId]

	index %= cm.pageSize
	if index >= len(page.data) {
		return
	}

	data := &page.data[index]
	data.exists = false
}

type ChunkMapElement[T any] struct {
	data []ChunkMapElementData[T]
}

type ChunkMapElementData[T any] struct {
	exists bool
	value  T
}
