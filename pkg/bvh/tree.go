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

package bvh

import (
	"cmp"
	"github.com/negrel/assert"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"math"
	"math/bits"
	"slices"
)

type node struct {
	childIndex int32 // if < 0 then points to leaf
}

type leaf struct {
	id ecs.Entity
}

type component struct {
	entity ecs.Entity
	aabb   stdcomponents.AABB
	code   uint64
}

func NewTree(layer stdcomponents.CollisionLayer) Tree {
	return Tree{
		nodes:      ecs.NewPagedArray[node](),
		AabbNodes:  ecs.NewPagedArray[stdcomponents.AABB](),
		leaves:     ecs.NewPagedArray[leaf](),
		AabbLeaves: ecs.NewPagedArray[stdcomponents.AABB](),
		codes:      ecs.NewPagedArray[uint64](),
		components: ecs.NewPagedArray[component](),
		layer:      layer,
	}
}

type Tree struct {
	nodes      ecs.PagedArray[node]
	AabbNodes  ecs.PagedArray[stdcomponents.AABB]
	leaves     ecs.PagedArray[leaf]
	AabbLeaves ecs.PagedArray[stdcomponents.AABB]
	codes      ecs.PagedArray[uint64]
	components ecs.PagedArray[component]
	layer      stdcomponents.CollisionLayer

	componentsSlice []component
}

func (t *Tree) AddComponent(entity ecs.Entity, aabb stdcomponents.AABB) {
	code := t.morton2D(&aabb)
	t.components.Append(component{
		entity: entity,
		aabb:   aabb,
		code:   code,
	})
}

func (t *Tree) Build() {
	// Reset tree
	t.nodes.Reset()
	t.AabbNodes.Reset()
	t.leaves.Reset()
	t.AabbLeaves.Reset()
	t.codes.Reset()

	// Extract and sort components by morton code
	if cap(t.componentsSlice) < t.components.Len() {
		t.componentsSlice = make([]component, 0, max(cap(t.componentsSlice)*2, t.components.Len()))
	}

	t.componentsSlice = t.components.Raw(t.componentsSlice)

	slices.SortFunc(t.componentsSlice, func(a, b component) int {
		return cmp.Compare(a.code, b.code)
	})

	// Add leaves
	for i := range t.componentsSlice {
		component := &t.componentsSlice[i]
		t.leaves.Append(leaf{id: component.entity})
		t.AabbLeaves.Append(component.aabb)
		t.codes.Append(component.code)
	}
	t.components.Reset()

	if t.leaves.Len() == 0 {
		return
	}

	// Add root node
	t.nodes.Append(node{-1})
	t.AabbNodes.Append(stdcomponents.AABB{})

	type buildTask struct {
		parentIndex     int
		start           int
		end             int
		childrenCreated bool
	}

	stack := [64]buildTask{
		{parentIndex: 0, start: 0, end: t.leaves.Len() - 1, childrenCreated: false},
	}
	stackLen := 1

	for stackLen > 0 {
		stackLen--
		// Pop the last task
		task := stack[stackLen]

		if !task.childrenCreated {
			if task.start == task.end {
				// Leaf node
				t.nodes.Get(task.parentIndex).childIndex = -int32(task.start)
				t.AabbNodes.Set(task.parentIndex, t.AabbLeaves.GetValue(task.start))
				continue
			}

			split := t.findSplit(task.start, task.end)

			// Create left and right nodes
			leftIndex := t.nodes.Len()
			t.nodes.Append(node{-1})
			t.nodes.Append(node{-1})
			t.AabbNodes.Append(stdcomponents.AABB{})
			t.AabbNodes.Append(stdcomponents.AABB{})

			// Set parent's childIndex to leftIndex
			t.nodes.Get(task.parentIndex).childIndex = int32(leftIndex)

			// Push parent task back with childrenCreated=true
			stack[stackLen] = buildTask{
				parentIndex:     task.parentIndex,
				start:           task.start,
				end:             task.end,
				childrenCreated: true,
			}
			stackLen++

			// Push right child task (split+1 to end)
			stack[stackLen] = buildTask{
				parentIndex:     leftIndex + 1,
				start:           split + 1,
				end:             task.end,
				childrenCreated: false,
			}
			stackLen++

			// Push left child task (start to split)
			stack[stackLen] = buildTask{
				parentIndex:     leftIndex,
				start:           task.start,
				end:             split,
				childrenCreated: false,
			}
			stackLen++
		} else {
			// Merge children's AABBs into parent
			leftChildIndex := int(t.nodes.Get(task.parentIndex).childIndex)
			rightChildIndex := leftChildIndex + 1

			leftAABB := t.AabbNodes.Get(leftChildIndex)
			rightAABB := t.AabbNodes.Get(rightChildIndex)

			merged := t.mergeAABB(leftAABB, rightAABB)
			t.AabbNodes.Set(task.parentIndex, merged)
		}
	}
	t.components.Reset()
}

func (t *Tree) Layer() stdcomponents.CollisionLayer {
	return t.layer
}

func (t *Tree) Query(aabb *stdcomponents.AABB, result []ecs.Entity) []ecs.Entity {
	if t.nodes.Len() == 0 { // Handle empty tree
		return result
	}

	// Use stack-based traversal
	const stackSize = 32
	stack := [stackSize]int32{0}
	stackLen := 1

	for stackLen > 0 {
		stackLen--
		nodeIndex := int(stack[stackLen])
		a := t.AabbNodes.Get(nodeIndex)
		b := aabb

		// Early exit if no AABB overlap
		if !t.aabbOverlap(a, b) {
			continue
		}

		node := t.nodes.Get(nodeIndex)
		if node.childIndex <= 0 {
			// Is a leaf
			index := -int(node.childIndex)
			leafAabb := t.AabbLeaves.Get(index)
			if t.aabbOverlap(leafAabb, aabb) {
				result = append(result, t.leaves.Get(index).id)
			}
			continue
		}

		// Push child indices (right and left) onto the stack.
		stack[stackLen] = node.childIndex + 1
		stack[stackLen+1] = node.childIndex
		stackLen += 2
	}

	return result
}

// go:inline aabbOverlap checks if two AABB intersect
func (t *Tree) aabbOverlap(a, b *stdcomponents.AABB) bool {
	return a.Max.X >= b.Min.X && a.Min.X <= b.Max.X &&
		a.Max.Y >= b.Min.Y && a.Min.Y <= b.Max.Y
}

// findSplit finds the position where the highest bit changes
func (t *Tree) findSplit(start, end int) int {
	// Identical Morton sortedMortonCodes => split the range in the middle.
	first := t.codes.GetValue(start)
	last := t.codes.GetValue(end)

	if first == last {
		return (start + end) >> 1
	}

	// Calculate the number of highest bits that are the same
	// for all objects, using the count-leading-zeros intrinsic.
	commonPrefix := bits.LeadingZeros64(first ^ last)

	// Use binary search to find where the next bit differs.
	// Specifically, we are looking for the highest object that
	// shares more than commonPrefix bits with the first one.
	split := start
	step := end - start

	for {
		step = (step + 1) >> 1   // exponential decrease
		newSplit := split + step // proposed new position

		if newSplit < end {
			splitCode := t.codes.GetValue(newSplit)
			splitPrefix := bits.LeadingZeros64(first ^ splitCode)
			if splitPrefix > commonPrefix {
				split = newSplit
			}
		}

		if step <= 1 {
			break
		}
	}

	return split
}

// mergeAABB combines two AABB
func (t *Tree) mergeAABB(a, b *stdcomponents.AABB) stdcomponents.AABB {
	return stdcomponents.AABB{
		Min: vectors.Vec2{
			X: min(a.Min.X, b.Min.X),
			Y: min(a.Min.Y, b.Min.Y),
		},
		Max: vectors.Vec2{
			X: max(a.Max.X, b.Max.X),
			Y: max(a.Max.Y, b.Max.Y),
		},
	}
}

// Expands a 16-bit integer into 32 bits by inserting 1 zero after each bit
func (t *Tree) expandBits2D(v uint32) uint32 {
	v = (v | (v << 8)) & 0x00FF00FF
	v = (v | (v << 4)) & 0x0F0F0F0F
	v = (v | (v << 2)) & 0x33333333
	v = (v | (v << 1)) & 0x55555555
	return v
}

const mortonPrecision = (1 << 16) - 1

func (t *Tree) morton2D(aabb *stdcomponents.AABB) uint64 {
	center := aabb.Center()
	// Scale coordinates to 16-bit integers
	//assert.True(center.X >= 0 && center.Y >= 0, "morton2D: center out of range")

	xx := uint64(float64(center.X) * mortonPrecision)
	yy := uint64(float64(center.Y) * mortonPrecision)

	assert.True(xx < math.MaxUint64, "morton2D: x out of range")
	assert.True(yy < math.MaxUint64, "morton2D: y out of range")

	// Spread the bits of x into the even positions
	xx = (xx | (xx << 16)) & 0x0000FFFF0000FFFF
	xx = (xx | (xx << 8)) & 0x00FF00FF00FF00FF
	xx = (xx | (xx << 4)) & 0x0F0F0F0F0F0F0F0F
	xx = (xx | (xx << 2)) & 0x3333333333333333
	xx = (xx | (xx << 1)) & 0x5555555555555555

	// Spread the bits of y into the even positions and shift to odd positions
	yy = (yy | (yy << 16)) & 0x0000FFFF0000FFFF
	yy = (yy | (yy << 8)) & 0x00FF00FF00FF00FF
	yy = (yy | (yy << 4)) & 0x0F0F0F0F0F0F0F0F
	yy = (yy | (yy << 2)) & 0x3333333333333333
	yy = (yy | (yy << 1)) & 0x5555555555555555

	// Combine x (even bits) and y (odd bits)
	return xx | (yy << 1)
}
