package state

import (
	"fmt"
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

func (p *Point) Distance(o *Point) float64 {
	return math.Sqrt((p.X-o.X)*(p.X-o.X) +
		(p.Y-p.Y)*(p.Y-o.Y))
}

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

func (b *Box) Intersects(t *Box) bool {
	point1, point2, distance := b.ClosestPair(t)
	fmt.Println("Closest pair", point1, point2, distance)
/*	points := [2]Point{point1, point2}
  sort.Sort(ByX(points))
	points[0].X + distance / 2 */
	return distance < 0
	/*
		  points := [8]Point{Point{b.X, b.Y}, Point{b.X + b.W, b.Y},
		  Point{b.X + b.W, b.Y + b.H}, Point{b.X, b.Y + b.H},
		  Point{t.X, t.Y}, Point{t.X + t.W, t.Y},
		  Point{t.X + t.W, t.Y + t.H}, Point{t.X, t.Y + t.H}}
				sort.Sort(ByX(points))
				points1 := points[0:4]
				points2 := points[4:8]
				xMid := points1[3].Distance(points2[0])/2.0 + points1[3].X
				dLMin := points1[3].Distance(Point(xMid, 0.0))
				dRMin := points2[0].Distance(Point(xMid, 0.0))
				dLRMin := 9999999999.9
				for i := 0; i < len(points); i++ {

				}
	*/
	/*
		fmt.Println(b, t)
		does := (((b.X <= t.X && t.X <= b.X+b.W) ||
			(b.X <= t.X+t.W && t.X+t.W <= b.X+b.W)) &&
			((b.Y <= t.Y && t.Y <= b.Y+b.H) ||
				(b.Y <= t.Y+t.H && t.Y+t.H <= b.Y+b.H)) ||
			((t.X <= b.X && b.X <= t.X+t.W) ||
				(t.X <= b.X+b.W && b.X+b.W <= t.X+t.W)) &&
				((t.Y <= b.Y && b.Y <= t.Y+t.H) ||
					(t.Y <= b.Y+b.H && b.Y+b.H < b.Y+b.H)))
		fmt.Println("does", does)
		return does
	*/
}
