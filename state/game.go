package state

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	player    *Player
	chaser    *Chaser
	score     int32
	playfield *Playfield
	target    *Location
	rng       *rand.Rand
	wall      []*Box
	area      []*Box
	opening   []*Box
)

// Location represents a coordinate in the game.
type Location struct {
	X int32
	Y int32
	Z int32
}

// Player represents the player's avatar in the game.
type Player struct {
	Location  Location
	Direction float64
	Speed     float64
	Bounds    *Box
	HalfW     int32
	HalfH     int32
}

// Chaser represents a mobile in the game.
type Chaser struct {
	Location  Location
	Direction float64
	Speed     float64
	Bounds    *Box
}

// Playfield represents the playable area in the game.
type Playfield struct {
	Width  int32
	Height int32
}

// InitState initializes the game structures to a known, playable state.
func InitState(p *Playfield) {
	seed := time.Now().UnixNano()
	//var seed int64 = 1421978894553050386
	// this seed has a maze gen bug
	seed = 1422165503556202364
	fmt.Println("seed:", seed)
	rng = rand.New(rand.NewSource(seed))
	playfield = p
	initPlayer()
	initChaser()
	createWalls()
}

func random(min int32, max int32) int32 {
	fmt.Println("random", min, max)
	return int32(rng.Intn(int(max-min))) + min
}

func initPlayer() {
	player = &Player{Location{0, 0, 10}, -0.75, 1.0, nil, 0, 0}
}

func initChaser() {
	chaser = &Chaser{Location{10, 10, 9}, -0.5, 0.75, nil}
}

func createWalls() {
	// the following has been removed since it's not as appropriate for what I want to create
	// wall = make([]*Box, 0, 100)
	// opening = make([]*Box, 0, 100)
	// makeRecursiveMaze(&Box{0, 0, playfield.Width, playfield.Height, 0})

	// for testing, I'll use Prim's method since I'm intending to have a cell- or grid-based game
	var cellSize int32 = 40
	grid := makeRandomizedPrimsMaze(&Box{X: 0, Y: 0, W: playfield.Width / cellSize, H: playfield.Height / cellSize})
	for row := 0; row < len(grid); row++ {
		cells := grid[row]
		for col := 0; col < len(cells); col++ {
			if grid[row][col] {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}

		fmt.Println("")
	}
}

// GetPlayer is a convenience method to return the Player struct.
func GetPlayer() *Player {
	return player
}

// GetChaser is a convenience method to return the Chaser struct.
func GetChaser() *Chaser {
	return chaser
}

// GetScore is a convenience method to return the current score.
func GetScore() int32 {
	return score
}

// UpdateState is called once per frame or tick to update the game state.
func UpdateState() {
	player.Location = translateLocation(player.Location, player.Direction, player.Speed)
	collisionBox := player.GetCollisionBox()

	// TODO fix player bounds using collision box
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

	// check for wall collisions
	for i := 0; i < len(wall); i++ {
		w := wall[i]
		// TODO check sprite bounds instead of the center of the sprite
		// floats so we gotta check within a range
		/*
			if ((w.X-0.5 <= player.Location.X && player.Location.X <= w.X+0.5) ||
			(w.X+w.W-0.5 <= player.Location.X && player.Location.X <= w.X+w.W+0.5)) &&
			((w.Y-0.5 <= player.Location.Y && player.Location.Y <= w.Y+0.5) ||
			(w.Y+w.H-0.5 <= player.Location.Y && player.Location.Y <= w.Y+w.H+0.5)) {
		*/
		if collisionBox.Intersects(w) {
			player.Speed = 0
		}
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

// SetClickLocation sets the location of a player's mouse click.
func SetClickLocation(click Location, behavior uint8) {
	target = &click

	// calculate the target angle
	y := player.Location.Y - click.Y
	x := click.X - player.Location.X

	if behavior == BehaviorAvoid {
		x *= -1
		y *= -1
	}

	theta := math.Atan2(float64(y), float64(x))

	player.SetDirection(theta)

	// also reset the speed (in case the sprite stopped)
	player.SetSpeed(1.0)
}

// UpdateChaser updates the state of the chaser to follow the player.
func UpdateChaser() {
	// TODO refactor this - copied chase logic from SetClickLocation
	// set the chaser to follow the player
	y := chaser.Location.Y - player.Location.Y
	x := player.Location.X - chaser.Location.X
	theta := math.Atan2(float64(y), float64(x))
	chaser.SetDirection(theta)
	chaser.SetSpeed(0.75)
}

// currently 2D, Z ignored
func translateLocation(origin Location, theta float64, speed float64) Location {
	speedFudge := 2.0
	x := int32(float64(origin.X) + math.Cos(theta)*speed*speedFudge)
	y := int32(float64(origin.Y) + math.Sin(theta)*-speed*speedFudge)
	return Location{x, y, origin.Z}
}

// AdjustSpeed adds the amount to the speed at which the player moves.
func (p *Player) AdjustSpeed(amount float64) {
	if p.Speed > -1.0 && p.Speed < 1.0 {
		p.Speed += amount
	}
}

// SetSpeed sets the speed of the player.
func (p *Player) SetSpeed(amount float64) {
	p.Speed = amount
}

// SetSpeed sets the speed of the chaser.
func (c *Chaser) SetSpeed(amount float64) {
	c.Speed = amount
}

// SetDirection sets the absolute direction in radians of the player.
func (p *Player) SetDirection(theta float64) {
	p.Direction = theta
}

// SetDirection sets the absolute direction in radians of the chaser.
func (c *Chaser) SetDirection(theta float64) {
	c.Direction = theta
}

// AdjustDirection adds the amount of radians to the player's direction.
func (p *Player) AdjustDirection(amount float64) {
	p.Direction += amount
	if p.Direction > TwoPi {
		p.Direction = p.Direction - TwoPi
	} else if p.Direction < NegTwoPi {
		p.Direction = p.Direction - NegTwoPi
	}
}

// GetCollisionBox returns the bounds of the player sprite.
func (p *Player) GetCollisionBox() *Box {
	b := &Box{p.Location.X - p.HalfW, p.Location.Y - p.HalfH,
		p.Bounds.W, p.Bounds.H, 0}
	return b
}

// GetWalls is a convenience method to get the impassable areas of the playfield.
func GetWalls() []*Box {
	return wall
}

// GetAreas is a convenience method to get the subdivided areas of the playfield, suitable for R-trees.  This method is likely to go away.
func GetAreas() []*Box {
	return area
}

// GetOpenings is a convenience method to get the openings in the playfield walls.  This method is likely to go away.
func GetOpenings() []*Box {
	return opening
}
