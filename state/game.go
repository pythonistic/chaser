package state

import (
	"math"
)

const TWO_PI = math.Pi * 2
const NEG_TWO_PI = math.Pi * -2
const BEHAVIOR_AVOID = 1
const BEHAVIOR_ATTRACT = 2

var (
	player    *Player
	chaser    *Chaser
	score     uint32
	playfield *Playfield
	target    *Location
)

type Location struct {
	X float64
	Y float64
	Z float64
}

type Player struct {
	Location  Location
	Direction float64
	Speed     float64
}

type Chaser struct {
	Location  Location
	Direction float64
	Speed     float64
}

type Playfield struct {
	Width  float64
	Height float64
}

func InitState(p *Playfield) {
	playfield = p
	initPlayer()
	initChaser()
}

func initPlayer() {
	player = &Player{Location{0.0, 0.0, 0.0}, -0.75, 1.0}
}

func initChaser() {
	chaser = &Chaser{Location{10.0, 10.0, 10.0}, -0.5, 0.9}
}

func GetPlayer() *Player {
	return player
}

func GetChaser() *Chaser {
	return chaser
}

func GetScore() uint32 {
	return score
}

func UpdateState() {
	player.Location = translateLocation(player.Location, player.Direction, player.Speed)

	if player.Location.X >= playfield.Width {
		player.Location.X = playfield.Width - 1
		player.Speed = 0
	} else if player.Location.X < 0 {
		player.Location.X = 0
		player.Speed = 0
	}

	if player.Location.Y >= playfield.Width {
		player.Location.Y = playfield.Height - 1
		player.Speed = 0
	} else if player.Location.Y < 0 {
		player.Location.Y = 0
		player.Speed = 0
	}

	chaser.Location = translateLocation(chaser.Location, chaser.Direction, chaser.Speed)

	if chaser.Location.X >= playfield.Width {
		chaser.Location.X = playfield.Width - 1
		chaser.Speed = 0
	} else if chaser.Location.X < 0 {
		chaser.Location.X = 0
		chaser.Speed = 0
	}

	if chaser.Location.Y >= playfield.Width {
		chaser.Location.Y = playfield.Height - 1
		chaser.Speed = 0
	} else if chaser.Location.Y < 0 {
		chaser.Location.Y = 0
		chaser.Speed = 0
	}
}

func SetClickLocation(click Location, behavior uint8) {
	target = &click

	// calculate the target angle
	y := player.Location.Y - click.Y
	x := click.X - player.Location.X

	if behavior == BEHAVIOR_AVOID {
		x *= -1
		y *= -1
	}

	theta := math.Atan2(y, x)

	player.SetDirection(theta)
}

// currently 2D, Z ignored
func translateLocation(origin Location, theta float64, speed float64) Location {
	x := origin.X + math.Cos(theta)*speed
	y := origin.Y + math.Sin(theta)*-speed
	return Location{x, y, origin.Z}
}

func (p *Player) AdjustSpeed(amount float64) {
	if p.Speed > -1.0 && p.Speed < 1.0 {
		p.Speed += amount
	}
}

func (p *Player) SetDirection(theta float64) {
	p.Direction = theta
}

func (p *Player) AdjustDirection(amount float64) {
	p.Direction += amount
	if p.Direction > TWO_PI {
		p.Direction = p.Direction - TWO_PI
	} else if p.Direction < NEG_TWO_PI {
		p.Direction = p.Direction - NEG_TWO_PI
	}
}
