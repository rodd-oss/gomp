/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type ChunkMap[T any] struct {
	buffer []ChunkMapElement[T]

	initialChunkCapPower int
	initialBufferCap     int
	chunkCapPower        int
	bufferCapPower       int
}

// const a = 1<<12 | 1<<13 | 1<<14

func NewChunkMap[T any](bufferCapacityPower int, chunkCapacityPower int) (arr *ChunkMap[T]) {
	arr = new(ChunkMap[T])

	arr.bufferCapPower = bufferCapacityPower
	arr.initialBufferCap = 1 << bufferCapacityPower
	arr.chunkCapPower = chunkCapacityPower
	arr.initialChunkCapPower = chunkCapacityPower

	arr.buffer = make([]ChunkMapElement[T], 1<<bufferCapacityPower)

	return arr
}

func (cm *ChunkMap[T]) Get(index int) (value T, ok bool) {
	pageId := FastIntLog2(index>>cm.initialChunkCapPower + 1)
	if pageId >= len(cm.buffer) {
		return value, false
	}
	page := cm.buffer[pageId]

	index -= ((1<<pageId - 1) << cm.initialChunkCapPower)
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

	pageId := FastIntLog2(index>>cm.initialChunkCapPower + 1)
	bufferLastIndex := len(cm.buffer) - 1
	if pageId > bufferLastIndex {
		delta := pageId - bufferLastIndex
		if delta < 1<<cm.bufferCapPower {
			cm.buffer = append(cm.buffer, make([]ChunkMapElement[T], 1<<cm.bufferCapPower)...)
			cm.bufferCapPower++
		} else {
			cm.buffer = append(cm.buffer, make([]ChunkMapElement[T], delta)...)
		}
	}
	page = &cm.buffer[pageId]

	index -= ((1<<pageId - 1) << cm.initialChunkCapPower)
	if index >= len(page.data) {
		page.data = make([]ChunkMapElementData[T], 1<<(cm.chunkCapPower+pageId))
	}

	data := &page.data[index]
	data.exists = true
	data.value = value
}

func (cm *ChunkMap[T]) Delete(index int) {
	var page *ChunkMapElement[T]

	pageId := FastIntLog2(index>>cm.initialChunkCapPower + 1)
	if pageId >= len(cm.buffer) {
		return
	}
	page = &cm.buffer[pageId]

	index -= ((1<<pageId - 1) << cm.initialChunkCapPower)
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
