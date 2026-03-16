package bench

import (
	"math"
	"math/rand"
	"testing"
)

const (
	numCircles = 10000
	numAABBs   = 10000
)

// Circle represents a circle with center (x, y) and radius r
type Circle struct {
	X, Y, R float64
}

// AABB represents an axis-aligned bounding box
type AABB struct {
	XMin, YMin, XMax, YMax float64
}

// circleCircleIntersect1 checks intersection using distance-based method
func circleCircleIntersect1(c1, c2 Circle) bool {
	dx := c1.X - c2.X
	dy := c1.Y - c2.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < c1.R+c2.R
}

// circleCircleIntersect2 checks intersection using squared distance
func circleCircleIntersect2(c1, c2 Circle) bool {
	dx := c1.X - c2.X
	dy := c1.Y - c2.Y
	rSum := c1.R + c2.R
	return dx*dx+dy*dy < rSum*rSum
}

// circleAABBIntersect checks if circle intersects with AABB
func circleAABBIntersect(c Circle, a AABB) bool {
	// Find closest point on AABB to circle center
	closestX := math.Max(a.XMin, math.Min(c.X, a.XMax))
	closestY := math.Max(a.YMin, math.Min(c.Y, a.YMax))

	// Calculate distance between closest point and circle center
	dx := c.X - closestX
	dy := c.Y - closestY
	distance := math.Sqrt(dx*dx + dy*dy)

	return distance < c.R
}

// aabbIntersect checks if two AABBs intersect
func aabbIntersect(a, b AABB) bool {
	// Check if one box is to the left of the other
	return a.XMax >= b.XMin && a.XMin <= b.XMax && a.YMax >= b.YMin && a.YMin <= b.YMax
}

func aabbIntersect2(a, b AABB) bool {
	// Check if one box is to the left of the other
	if a.XMax < b.XMin || b.XMax < a.XMin {
		return false
	}

	return !(a.YMax < b.YMin || b.YMax < a.YMin)
}

// generateRandomCircle generates a random circle
func generateRandomCircle() Circle {
	return Circle{
		X: rand.Float64() * 100,
		Y: rand.Float64() * 100,
		R: rand.Float64()*10 + 5,
	}
}

// generateRandomAABB generates a random AABB
func generateRandomAABB() AABB {
	x1 := rand.Float64() * 100
	y1 := rand.Float64() * 100
	w := rand.Float64()*20 + 10
	h := rand.Float64()*20 + 10
	return AABB{
		XMin: x1,
		YMin: y1,
		XMax: x1 + w,
		YMax: y1 + h,
	}
}

func BenchmarkCircleCircleIntersect1(b *testing.B) {
	rand.Seed(42)
	circles := make([]Circle, numCircles)
	for i := range circles {
		circles[i] = generateRandomCircle()
	}

	for b.Loop() {
		lastIndex := len(circles)
		for j := 0; j < lastIndex; j++ {
			for k := j + 1; k < lastIndex; k++ {
				circleCircleIntersect1(circles[j], circles[k])
			}
		}
	}
}

func BenchmarkCircleCircleIntersect2(b *testing.B) {
	rand.Seed(42)
	circles := make([]Circle, numCircles)
	for i := range circles {
		circles[i] = generateRandomCircle()
	}

	for b.Loop() {
		lastIndex := len(circles)
		for j := 0; j < lastIndex; j++ {
			for k := j + 1; k < lastIndex; k++ {
				circleCircleIntersect2(circles[j], circles[k])
			}
		}
	}
}

//
//func BenchmarkCircleAABBIntersect(b *testing.B) {
//	rand.Seed(42)
//	circles := make([]Circle, numCircles)
//	aabbs := make([]AABB, numAABBs)
//	for i := range circles {
//		circles[i] = generateRandomCircle()
//		aabbs[i] = generateRandomAABB()
//	}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		for j := 0; j < len(aabbs); j++ {
//			circleAABBIntersect(circles[j], aabbs[j])
//		}
//	}
//}

func BenchmarkAABBIntersect(b *testing.B) {
	rand.Seed(42)
	aabbs := make([]AABB, numAABBs)
	for i := range aabbs {
		aabbs[i] = generateRandomAABB()
	}

	for b.Loop() {
		lastIndex := len(aabbs)
		for j := 0; j < lastIndex; j++ {
			for k := j + 1; k < lastIndex; k++ {
				aabbIntersect(aabbs[j], aabbs[k])
			}
		}
	}
}

func BenchmarkAABBIntersect2(b *testing.B) {
	rand.Seed(42)
	aabbs := make([]AABB, numAABBs)
	for i := range aabbs {
		aabbs[i] = generateRandomAABB()
	}

	for b.Loop() {
		lastIndex := len(aabbs)
		for j := 0; j < lastIndex; j++ {
			for k := j + 1; k < lastIndex; k++ {
				aabbIntersect2(aabbs[j], aabbs[k])
			}
		}
	}
}
