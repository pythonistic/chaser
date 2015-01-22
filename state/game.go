package state

import (
	"math"
	"math/rand"
	"time"
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
	rng       *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	wall      []*Wall
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

type Wall struct {
	X1 float64
	Y1 float64
	W  float64
	H  float64
}

func InitState(p *Playfield) {
	playfield = p
	initPlayer()
	initChaser()
	createWalls()
}

func initPlayer() {
	player = &Player{Location{0.0, 0.0, 10.0}, -0.75, 1.0}
}

func initChaser() {
	chaser = &Chaser{Location{10.0, 10.0, 9.0}, -0.5, 0.75}
}

func createWalls() {
	wall = make([]*Wall, 10, 10)
	boundWidth := playfield.Width - 20
	boundHeight := playfield.Height - 20
	//bounds := sdl.Rect{10, 10, boundWidth, boundHeight}
	for i := 0; i < 10; i++ {
		width := rng.Float64()*100 + 1
		height := rng.Float64()*100 + 1

		if width >= height {
			height = 10
		} else {
			width = 10
		}

		x1 := rng.Float64()*boundWidth + 10
		y1 := rng.Float64()*boundHeight + 10
		wall[i] = &Wall{x1, y1, width, height}
	}
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

	UpdateChaser()
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

func UpdateChaser() {
	// TODO refactor this - copied chase logic from SetClickLocation
	// set the chaser to follow the player
	y := chaser.Location.Y - player.Location.Y
	x := player.Location.X - chaser.Location.X
	theta := math.Atan2(y, x)
	chaser.SetDirection(theta)
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

func (c *Chaser) SetDirection(theta float64) {
	c.Direction = theta
}

func (p *Player) AdjustDirection(amount float64) {
	p.Direction += amount
	if p.Direction > TWO_PI {
		p.Direction = p.Direction - TWO_PI
	} else if p.Direction < NEG_TWO_PI {
		p.Direction = p.Direction - NEG_TWO_PI
	}
}

func GetWalls() []*Wall {
	return wall
}
