package state

import (
	"math"
)

// Point represents a point or 2D coordinate.
type Point struct {
	X int32
	Y int32
}

// Box is a rectangle with a Z coordinate, intended for stacking when rendering.
type Box struct {
	X int32
	Y int32
	W int32
	H int32
	Z int32
}

const (
	// TwoPi is 2 * Pi, used for calculating counter-clockwise rotation in radians.
	TwoPi = math.Pi * 2
	// NegTwoPi is -2 * Pi, used for calculating clockwise rotation in radians.
	NegTwoPi = math.Pi * -2
	// BehaviorAvoid indicates another object will try to avoid a target.
	BehaviorAvoid = 1
	// BehaviorAttract indicates another object will by attracted to a target.
	BehaviorAttract = 2
	// WallWidth indicates the default width of a wall when building levels or pathing.
	WallWidth = 10
	// CellSize indicates the default square size of a cell for grids/tiles.  Sprites may be smaller.
	CellSize = 40
)

// ByX is an interface for sorting Point instances by the X coordinate.
type ByX []Point

func (p ByX) Len() int {
	return len(p)
}

func (p ByX) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByX) Less(i, j int) bool {
	return p[i].X < p[j].X
}

// Contains returns True when the Point p is within Box b.
func (b *Box) Contains(p *Point) bool {
	return b.X <= p.X && p.X < b.X+b.W &&
		b.Y <= p.Y && p.Y < b.Y+b.H
}

// Distance calculates the distance between two Point instances.
func (p *Point) Distance(o *Point) float64 {
	return math.Sqrt(float64((p.X-o.X)*(p.X-o.X) +
		(p.Y-p.Y)*(p.Y-o.Y)))
}

// ClosestPair finds the closest pair of Point instances with two Box instances.
// This is implemented using a brute force method.
func (b *Box) ClosestPair(t *Box) (*Point, *Point, float64) {
	points1 := [4]Point{Point{b.X, b.Y}, Point{b.X + b.W, b.Y},
		Point{b.X + b.W, b.Y + b.H}, Point{b.X, b.Y + b.H}}
	points2 := [4]Point{Point{t.X, t.Y}, Point{t.X + t.W, t.Y},
		Point{t.X + t.W, t.Y + t.H}, Point{t.X, t.Y + t.H}}
	var minimumDistance = 99999999999.9
	var closest1, closest2 *Point

	for i := 0; i < len(points1); i++ {
		for j := 0; j < len(points2); j++ {
			dist := points1[i].Distance(&points2[j])
			if dist < minimumDistance {
				minimumDistance = dist
				closest1 = &points1[i]
				closest2 = &points2[j]
			}
		}
	}

	return closest1, closest2, minimumDistance
}

// Intersects returns true if the Box instances intersect.
func (b *Box) Intersects(t *Box) bool {
	if b.Y+b.H < t.Y {
		return false
	}

	if b.Y > t.Y+t.H {
		return false
	}

	if b.X+b.W < t.X {
		return false
	}

	if b.X > t.X+t.W {
		return false
	}

	return true
}
