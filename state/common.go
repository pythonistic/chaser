package state

import "fmt"

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

func (b *Box) Contains(p *Point) bool {
	return b.X <= p.X && p.X <= b.X+b.W &&
		b.Y <= p.Y && p.Y <= b.Y+b.H
}

func (b *Box) Intersects(t *Box) bool {
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
}
