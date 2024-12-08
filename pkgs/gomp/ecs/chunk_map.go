/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "math/bits"

type ChunkMap[T any] struct {
	buffer               []ChunkMapElement[T]
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
	for i := range arr.buffer {
		arr.buffer[i].startingIndex = ((1<<i - 1) << chunkCapacityPower)
	}

	return arr
}

func (cm *ChunkMap[T]) Get(index int) (value T, ok bool) {
	pageId := cm.getPageIdByIndex(index)
	if pageId >= len(cm.buffer) {
		return value, false
	}
	page := &cm.buffer[pageId]

	index -= page.startingIndex
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

	pageId := cm.getPageIdByIndex(index)
	bufferLastIndex := len(cm.buffer) - 1
	if pageId > bufferLastIndex {
		delta := pageId - bufferLastIndex
		if delta < 1<<cm.bufferCapPower {
			newBuf := make([]ChunkMapElement[T], 1<<cm.bufferCapPower)
			for i := range newBuf {
				newBuf[i].startingIndex = ((1<<(bufferLastIndex+i+1) - 1) << cm.initialChunkCapPower)
			}
			cm.buffer = append(cm.buffer, newBuf...)
			cm.bufferCapPower++
		} else {
			newBuf := make([]ChunkMapElement[T], delta)
			cm.buffer = append(cm.buffer, newBuf...)
			for i := range newBuf {
				newBuf[i].startingIndex = ((1<<(bufferLastIndex+i+1) - 1) << cm.initialChunkCapPower)
			}
		}
	}
	page = &cm.buffer[pageId]

	index -= page.startingIndex
	if index >= len(page.data) {
		page.data = make([]ChunkMapElementData[T], 1<<(cm.chunkCapPower+pageId))
	}

	data := &page.data[index]
	data.exists = true
	data.value = value
}

func (cm *ChunkMap[T]) Delete(index int) {
	var page *ChunkMapElement[T]

	pageId := cm.getPageIdByIndex(index)
	if pageId >= len(cm.buffer) {
		return
	}
	page = &cm.buffer[pageId]

	index -= page.startingIndex
	if index >= len(page.data) {
		return
	}

	data := &page.data[index]
	data.exists = false
}

func (cm *ChunkMap[T]) getPageIdByIndex(index int) int {
	return bits.Len64(uint64(index>>cm.initialChunkCapPower+1)) - 1
}

type ChunkMapElement[T any] struct {
	data          []ChunkMapElementData[T]
	startingIndex int
}

type ChunkMapElementData[T any] struct {
	exists bool
	value  T
}
