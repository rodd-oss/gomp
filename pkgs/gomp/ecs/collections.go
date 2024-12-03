/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

type Collection[T any] struct {
	buckets []Bucket[T]
	last    *Bucket[T]
	count   int

	initialBucketsCount int
	initialBucketSize   int
	defaultTValue       T
}

func (c *Collection[T]) init(buckets int, bucketSize int, defaultTvalue T) {
	c.initialBucketsCount = buckets
	c.initialBucketSize = bucketSize
	c.defaultTValue = defaultTvalue

	c.buckets = make([]Bucket[T], 1, buckets)
	c.buckets[0].init(0, bucketSize)
	c.last = &c.buckets[0]
}

func (c *Collection[T]) Get(id int) T {
	if id < c.initialBucketSize {
		if !c.buckets[0].Exists(id) {
			return c.defaultTValue
		}

		return c.buckets[0].Get(id)
	}

	bucketId := id / c.initialBucketSize
	if bucketId >= len(c.buckets) {
		return c.defaultTValue
	}

	bucket := &c.buckets[bucketId]
	id %= c.initialBucketSize

	if !bucket.Exists(id) {
		return c.defaultTValue
	}

	return bucket.Get(id)
}

func (c *Collection[T]) Exists(id int) bool {
	if id < c.initialBucketSize {
		return c.buckets[0].Exists(id)
	}

	bucketId := id / c.initialBucketSize
	bucket := &c.buckets[bucketId]
	id %= c.initialBucketSize

	return bucket.Exists(id)
}

func (c *Collection[T]) Last() (int, *T) {
	var bucketId int = c.last.size - 1
	return bucketId, &c.last.data[bucketId]
}

// [0,1,2][2][][] [0,1,2][][][9,10,11] [][sd][][9,10,11]
func (c *Collection[T]) Set(id int, val T) (value *T) {
	var bucketId int = 0

	if id >= c.initialBucketSize {
		bucketId = id / c.initialBucketSize
		id %= c.initialBucketSize
	}

	if bucketId >= len(c.buckets) {
		c.extend(bucketId - (len(c.buckets) - 1))
		c.buckets[bucketId].init(bucketId, c.initialBucketSize)
	}

	value = c.buckets[bucketId].Set(id, val, c.defaultTValue)
	c.count++
	return value
}

func (c *Collection[T]) Append(obj T) (int, *T) {
	if c.last == nil {
		c.init(c.initialBucketsCount, c.initialBucketSize, c.defaultTValue)
	}
	if c.last.CapLeft() < 1 {
		c.extend(1)
	}
	c.count++
	i, v := c.last.Append(obj)
	i = i + c.initialBucketSize*c.last.id
	return i, v
}

func (c *Collection[T]) SoftReduce() {
	for i := range c.buckets {
		c.last = &c.buckets[len(c.buckets)-1-i]
		if c.last.Len() != 0 {
			break
		}
	}

	c.last.SoftReduce()
	c.count--
}

func (c *Collection[T]) Swap(i int, j int) {
	iBucketId, iPos := c.indexToBucketIdAndPos(i)
	jBucketId, jPos := c.indexToBucketIdAndPos(j)

	if iBucketId == jBucketId {
		c.buckets[iBucketId].Swap(iPos, jPos)
		return
	}

	c.buckets[iBucketId].Set(iPos, c.buckets[iBucketId].Get(iPos), c.defaultTValue)
	c.buckets[jBucketId].Set(jPos, c.buckets[jBucketId].Get(jPos), c.defaultTValue)
}

func (c *Collection[T]) Clean() {
	for i := range c.buckets {
		c.buckets[i].Clean()
	}
}

func (c *Collection[T]) Len() (l int) {
	return c.count
}

func (c *Collection[T]) extend(buckets int) {
	capleft := cap(c.buckets) - len(c.buckets)
	if capleft >= buckets {
		c.buckets = append(c.buckets, make([]Bucket[T], buckets)...)
	} else {
		c.buckets = append(c.buckets, make([]Bucket[T], buckets, capleft+c.initialBucketsCount)...)
	}

	bucketId := len(c.buckets) - 1

	c.last = &c.buckets[bucketId]
	c.last.init(bucketId, c.initialBucketSize)
}

func (c *Collection[T]) indexToBucketIdAndPos(id int) (int, int) {
	var bucketId int = 0

	if id >= c.initialBucketSize {
		bucketId = id / c.initialBucketSize
		id %= c.initialBucketSize
	}

	return bucketId, id
}
