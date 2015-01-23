package state

import (
	"math"
	)

type Point struct {
	X float64
	Y float64
}

type Box struct {
	X float64
	Y float64
	W float64
	H float64
}

// sort array of points by X
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

func (b *Box) Contains(p *Point) bool {
	return b.X <= p.X && p.X <= b.X+b.W &&
		b.Y <= p.Y && p.Y <= b.Y+b.H
}

// calculate the distance between two points
func (p *Point) Distance(o *Point) float64 {
	return math.Sqrt((p.X-o.X)*(p.X-o.X) +
		(p.Y-p.Y)*(p.Y-o.Y))
}

// find the closest pair of points with two boxes
func (b *Box) ClosestPair(t *Box) (*Point, *Point, float64) {
	points1 := [4]Point{Point{b.X, b.Y}, Point{b.X + b.W, b.Y},
		Point{b.X + b.W, b.Y + b.H}, Point{b.X, b.Y + b.H}}
	points2 := [4]Point{Point{t.X, t.Y}, Point{t.X + t.W, t.Y},
		Point{t.X + t.W, t.Y + t.H}, Point{t.X, t.Y + t.H}}
	var minimumDistance float64 = 99999999999.9
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

// return true if the boxes intersect
func (b *Box) Intersects(t *Box) bool {
	if b.Y + b.H < t.Y {
		return false
	}

  if b.Y > t.Y + t.H {
		return false
	}

	if b.X + b.W < t.X {
		return false
	}

	if b.X > t.X + t.W {
		return false
	}

	return true
}
