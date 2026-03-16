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

package stdcomponents

import (
	"github.com/negrel/assert"
	"gomp/pkg/ecs"
	"gomp/vectors"
	"math"
	"math/bits"
	"slices"
)

const mortonPrecision = (1 << 16) - 1

type BvhNode struct {
	ChildIndex int32 // if < 0 then points to BvhLeaf
}

type BvhLeaf struct {
	Id ecs.Entity
}

type BvhComponent struct {
	Entity ecs.Entity
	Aabb   AABB
	Code   uint64
}

type BvhTree struct {
	Nodes      ecs.PagedArray[BvhNode]
	AabbNodes  ecs.PagedArray[AABB]
	Leaves     ecs.PagedArray[BvhLeaf]
	AabbLeaves ecs.PagedArray[AABB]
	Codes      ecs.PagedArray[uint64]
	Components ecs.PagedArray[BvhComponent]

	sorted []BvhComponent
}

func (t *BvhTree) Init() {
	t.Nodes = ecs.NewPagedArray[BvhNode]()
	t.AabbNodes = ecs.NewPagedArray[AABB]()
	t.Leaves = ecs.NewPagedArray[BvhLeaf]()
	t.AabbLeaves = ecs.NewPagedArray[AABB]()
	t.Codes = ecs.NewPagedArray[uint64]()
	t.Components = ecs.NewPagedArray[BvhComponent]()
}

func (t *BvhTree) AddComponent(entity ecs.Entity, aabb AABB) {
	code := t.morton2D(&aabb)
	t.Components.Append(BvhComponent{
		Entity: entity,
		Aabb:   aabb,
		Code:   code,
	})
}

func (t *BvhTree) Query(aabb AABB, result []ecs.Entity) []ecs.Entity {
	if t.Nodes.Len() == 0 { // Handle empty tree
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

		// Early exit if no AABB overlap
		if !t.aabbOverlap(*a, aabb) {
			continue
		}

		node := t.Nodes.Get(nodeIndex)
		if node.ChildIndex <= 0 {
			// Is a leaf
			index := -int(node.ChildIndex)
			leafAabb := t.AabbLeaves.Get(index)
			if t.aabbOverlap(*leafAabb, aabb) {
				result = append(result, t.Leaves.Get(index).Id)
			}
			continue
		}

		// Push child indices (right and left) onto the stack.
		stack[stackLen] = node.ChildIndex + 1
		stack[stackLen+1] = node.ChildIndex
		stackLen += 2
	}

	return result
}

// go:inline aabbOverlap checks if two AABB intersect
func (t *BvhTree) aabbOverlap(a, b AABB) bool {
	return a.Max.X >= b.Min.X && a.Min.X <= b.Max.X &&
		a.Max.Y >= b.Min.Y && a.Min.Y <= b.Max.Y
}

// Expands a 16-bit integer into 32 bits by inserting 1 zero after each bit
func (t *BvhTree) expandBits2D(v uint32) uint32 {
	v = (v | (v << 8)) & 0x00FF00FF
	v = (v | (v << 4)) & 0x0F0F0F0F
	v = (v | (v << 2)) & 0x33333333
	v = (v | (v << 1)) & 0x55555555
	return v
}

func (t *BvhTree) morton2D(aabb *AABB) uint64 {
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

func (t *BvhTree) Build() {
	// Reset tree
	t.Nodes.Reset()
	t.AabbNodes.Reset()
	t.Leaves.Reset()
	t.AabbLeaves.Reset()
	t.Codes.Reset()

	if cap(t.sorted) < t.Components.Len() {
		t.sorted = make([]BvhComponent, 0, t.Components.Len())
	}

	t.sorted = t.Components.Raw(t.sorted)

	slices.SortFunc(t.sorted, func(a, b BvhComponent) int {
		return int(a.Code - b.Code)
	})

	// Add leaves
	for i := range t.sorted {
		component := t.sorted[i]
		t.Leaves.Append(BvhLeaf{Id: component.Entity})
		t.AabbLeaves.Append(component.Aabb)
		t.Codes.Append(component.Code)
	}
	t.Components.Reset()
	t.sorted = t.sorted[:0]

	if t.Leaves.Len() == 0 {
		return
	}

	// Add root node
	t.Nodes.Append(BvhNode{-1})
	t.AabbNodes.Append(AABB{})

	type buildTask struct {
		parentIndex     int
		start           int
		end             int
		childrenCreated bool
	}

	stack := [64]buildTask{
		{parentIndex: 0, start: 0, end: t.Leaves.Len() - 1, childrenCreated: false},
	}
	stackLen := 1

	for stackLen > 0 {
		stackLen--
		// Pop the last task
		task := stack[stackLen]

		if !task.childrenCreated {
			if task.start == task.end {
				// Leaf node
				t.Nodes.Get(task.parentIndex).ChildIndex = -int32(task.start)
				t.AabbNodes.Set(task.parentIndex, t.AabbLeaves.GetValue(task.start))
				continue
			}

			split := t.findSplit(task.start, task.end)

			// Create left and right nodes
			leftIndex := t.Nodes.Len()
			t.Nodes.Append(BvhNode{-1})
			t.Nodes.Append(BvhNode{-1})
			t.AabbNodes.Append(AABB{})
			t.AabbNodes.Append(AABB{})

			// Set parent's childIndex to leftIndex
			t.Nodes.Get(task.parentIndex).ChildIndex = int32(leftIndex)

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
			leftChildIndex := int(t.Nodes.Get(task.parentIndex).ChildIndex)
			rightChildIndex := leftChildIndex + 1

			leftAABB := t.AabbNodes.Get(leftChildIndex)
			rightAABB := t.AabbNodes.Get(rightChildIndex)

			merged := t.mergeAABB(leftAABB, rightAABB)
			t.AabbNodes.Set(task.parentIndex, merged)
		}
	}
	t.Components.Reset()
}

func (t *BvhTree) findSplit(start, end int) int {
	// Identical Morton sortedMortonCodes => split the range in the middle.
	first := t.Codes.GetValue(start)
	last := t.Codes.GetValue(end)

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
			splitCode := t.Codes.GetValue(newSplit)
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

func (t *BvhTree) mergeAABB(a, b *AABB) AABB {
	return AABB{
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

type BvhTreeComponentManager = ecs.ComponentManager[BvhTree]

func NewBvhTreeComponentManager() BvhTreeComponentManager {
	return ecs.NewComponentManager[BvhTree](BvhTreeComponentId)
}
