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

package stdsystems

import (
	"cmp"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"gomp/vectors"
	"math/bits"
	"slices"
	"time"
)

func NewBvhTreeSystem() BvhTreeSystem {
	return BvhTreeSystem{}
}

type BvhTreeSystem struct {
	EntityManager *ecs.EntityManager
}

func (s *BvhTreeSystem) Init()                {}
func (s *BvhTreeSystem) Run(dt time.Duration) {}
func (s *BvhTreeSystem) Destroy()             {}

func (s *BvhTreeSystem) build(t *stdcomponents.BvhTree) {
	// Reset tree
	t.Nodes.Reset()
	t.AabbNodes.Reset()
	t.Leaves.Reset()
	t.AabbLeaves.Reset()
	t.Codes.Reset()

	var sorted []stdcomponents.BvhComponent
	sorted = t.Components.Raw(sorted)

	slices.SortFunc(sorted, func(a, b stdcomponents.BvhComponent) int {
		return cmp.Compare(a.Code, b.Code)
	})

	// Add leaves
	for i := range sorted {
		component := sorted[i]
		t.Leaves.Append(stdcomponents.BvhLeaf{Id: component.Entity})
		t.AabbLeaves.Append(component.Aabb)
		t.Codes.Append(component.Code)
	}
	t.Components.Reset()

	if t.Leaves.Len() == 0 {
		return
	}

	// Add root node
	t.Nodes.Append(stdcomponents.BvhNode{-1})
	t.AabbNodes.Append(stdcomponents.AABB{})

	type buildTask struct {
		parentIndex     int
		start           int
		end             int
		childrenCreated bool
	}

	stack := []buildTask{
		{parentIndex: 0, start: 0, end: t.Leaves.Len() - 1, childrenCreated: false},
	}

	for len(stack) > 0 {
		// Pop the last task
		task := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if !task.childrenCreated {
			if task.start == task.end {
				// Leaf node
				t.Nodes.Get(task.parentIndex).ChildIndex = -int32(task.start)
				t.AabbNodes.Set(task.parentIndex, t.AabbLeaves.GetValue(task.start))
				continue
			}

			split := s.findSplit(t, task.start, task.end)

			// Create left and right nodes
			leftIndex := t.Nodes.Len()
			t.Nodes.Append(stdcomponents.BvhNode{-1})
			t.Nodes.Append(stdcomponents.BvhNode{-1})
			t.AabbNodes.Append(stdcomponents.AABB{})
			t.AabbNodes.Append(stdcomponents.AABB{})

			// Set parent's childIndex to leftIndex
			t.Nodes.Get(task.parentIndex).ChildIndex = int32(leftIndex)

			// Push parent task back with childrenCreated=true
			stack = append(stack, buildTask{
				parentIndex:     task.parentIndex,
				start:           task.start,
				end:             task.end,
				childrenCreated: true,
			})

			// Push right child task (split+1 to end)
			stack = append(stack, buildTask{
				parentIndex:     leftIndex + 1,
				start:           split + 1,
				end:             task.end,
				childrenCreated: false,
			})

			// Push left child task (start to split)
			stack = append(stack, buildTask{
				parentIndex:     leftIndex,
				start:           task.start,
				end:             split,
				childrenCreated: false,
			})
		} else {
			// Merge children's AABBs into parent
			leftChildIndex := int(t.Nodes.Get(task.parentIndex).ChildIndex)
			rightChildIndex := leftChildIndex + 1

			leftAABB := t.AabbNodes.Get(leftChildIndex)
			rightAABB := t.AabbNodes.Get(rightChildIndex)

			merged := s.mergeAABB(leftAABB, rightAABB)
			t.AabbNodes.Set(task.parentIndex, merged)
		}
	}
	t.Components.Reset()
}

func (s *BvhTreeSystem) findSplit(t *stdcomponents.BvhTree, start, end int) int {
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

func (s *BvhTreeSystem) mergeAABB(a, b *stdcomponents.AABB) stdcomponents.AABB {
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
