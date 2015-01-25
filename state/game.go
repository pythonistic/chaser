package state

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const TWO_PI = math.Pi * 2
const NEG_TWO_PI = math.Pi * -2
const BEHAVIOR_AVOID = 1
const BEHAVIOR_ATTRACT = 2
const WALL_WIDTH = 10

var (
	player    *Player
	chaser    *Chaser
	score     int32
	playfield *Playfield
	target    *Location
	rng       *rand.Rand
	wall      []*Box
)

type Location struct {
	X int32
	Y int32
	Z int32
}

type Player struct {
	Location  Location
	Direction float64
	Speed     float64
	Bounds    *Box
	HalfW     int32
	HalfH     int32
}

type Chaser struct {
	Location  Location
	Direction float64
	Speed     float64
	Bounds    *Box
}

type Playfield struct {
	Width  int32
	Height int32
}

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
	wall = make([]*Box, 0, 100)

	makeMaze(&Box{0, 0, playfield.Width, playfield.Height}, 0)
}

func makeMaze(box *Box, depth uint8) {
	depth++
fmt.Println("DEPTH ", depth)
/*
  if depth > 2 {
		return
	}
	*/

	if box.X > playfield.Width {
		panic("box X bigger than playfield")
	}
	if box.Y > playfield.Height {
		panic("box H bigger than playfield")
	}
	// subdivide

	// player gets loaded after state is created, so we need to hardcode
	// player bounds + a little fudge
	//minX := player.Bounds.W  (70, 95 minimum)
	//minY := player.Bounds.H
	var playerWidth int32 = 120
	var playerHeight int32 = 120
	var minimumWidth = playerWidth*2 + WALL_WIDTH
	var minumumHeight = playerHeight*2 + WALL_WIDTH

	// fail early if nothing to do
	if (box.W < minimumWidth &&
		box.H < minumumHeight) ||
		box.W+box.X > playfield.Width ||
		box.H+box.Y > playfield.Height {
		return
	}

	minX := box.X
	minY := box.Y
	maxX := box.W + minX
	maxY := box.H + minY
	var xPos, yPos int32 = -1, -1

	// make sure there's enough room to make a horizontal wall
	if box.W > minimumWidth {
		xPos = random(minX + playerHeight, maxX - playerHeight)
	}

	// make sure there's enough room to make a veritcal wall
	if box.H > minumumHeight {
		yPos = random(minY + playerWidth, maxY - playerWidth)
	}

	// check for single wall cases
	if yPos < 0 {
		// vertical wall
		makeVerticalMazeWall(xPos, minY, maxY, playerHeight)
		makeMaze(&Box{minX, minY, xPos - minX, maxY - minY}, depth)
		makeMaze(&Box{xPos, minY, maxX - xPos, maxY - minY}, depth)
	} else if xPos < 0 {
		// horizontal wall
		makeHorizontalMazeWall(yPos, minX, maxX, playerWidth)
		makeMaze(&Box{minX, minY, maxX - minX, yPos - minY}, depth)
		makeMaze(&Box{minX, yPos, maxX - minX, maxY - yPos}, depth)
	} else {
		// two walls created, need to subdivide into four walls
		// skip one chamber randomly
		skip := random(1, 4)
		if skip != 1 {
			makeVerticalMazeWall(xPos, minY, yPos, playerHeight)
		} else {
			wall = append(wall, &Box{xPos, minY, WALL_WIDTH, yPos})
		}
		if skip != 2 {
			makeVerticalMazeWall(xPos, yPos, maxY, playerHeight)
		} else {
			wall = append(wall, &Box{xPos, yPos, WALL_WIDTH, maxY})
		}
		if skip != 3 {
			makeHorizontalMazeWall(yPos, minX, xPos, playerWidth)
		} else {
			wall = append(wall, &Box{minX, yPos, xPos, WALL_WIDTH})
		}
		if skip != 4 {
			makeHorizontalMazeWall(yPos, xPos, maxX, playerWidth)
		} else {
			wall = append(wall, &Box{xPos, yPos, maxX, WALL_WIDTH})
		}

		makeMaze(&Box{minX, minY, xPos - minX, yPos - minY}, depth) // top left
		makeMaze(&Box{xPos, minY, maxX - xPos, yPos - minY}, depth) // top right
		makeMaze(&Box{minX, yPos, xPos - minX, maxY - yPos}, depth) // bottom left
		makeMaze(&Box{xPos, yPos, maxX - minX, maxY - yPos}, depth) // bottom right
	}
}

func makeVerticalMazeWall(xPos int32, minY int32, maxY int32, playerHeight int32) {
	wallOpen := random(minY, maxY)
	if wallOpen > maxY-minY - playerHeight {
		// wall open is at the bottom
		wall = append(wall, &Box{xPos, minY, WALL_WIDTH, maxY - minY - playerHeight})
	} else if wallOpen < minY + playerHeight {
		// wall open is at the top
		wall = append(wall, &Box{xPos, playerHeight + minY, WALL_WIDTH, maxY - minY - playerHeight})
	} else {
		// wall open is in the middle, make two walls
		wall = append(wall, &Box{xPos, minY, WALL_WIDTH, wallOpen - minY})
		wall = append(wall, &Box{xPos, wallOpen + playerHeight, WALL_WIDTH, maxY - wallOpen - playerHeight})
	}
}

func makeHorizontalMazeWall(yPos int32, minX int32, maxX int32, playerWidth int32) {
	wallOpen := random(minX, maxX)
	if wallOpen > maxX-minX - playerWidth {
		// wall open is at the right
		wall = append(wall, &Box{minX, yPos, maxX - minX - playerWidth, WALL_WIDTH})
	} else if wallOpen < minX + playerWidth {
		// wall open is at the left
		wall = append(wall, &Box{playerWidth + minX, yPos, maxX - minX - playerWidth, WALL_WIDTH})
	} else {
		// wall open is in the middle, make two walls
		wall = append(wall, &Box{minX, yPos, wallOpen - minX, WALL_WIDTH})
		wall = append(wall, &Box{wallOpen + playerWidth, yPos, maxX - wallOpen - playerWidth, WALL_WIDTH})
	}
}

func GetPlayer() *Player {
	return player
}

func GetChaser() *Chaser {
	return chaser
}

func GetScore() int32 {
	return score
}

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

func SetClickLocation(click Location, behavior uint8) {
	target = &click

	// calculate the target angle
	y := player.Location.Y - click.Y
	x := click.X - player.Location.X

	if behavior == BEHAVIOR_AVOID {
		x *= -1
		y *= -1
	}

	theta := math.Atan2(float64(y), float64(x))

	player.SetDirection(theta)

	// also reset the speed (in case the sprite stopped)
	player.SetSpeed(1.0)
}

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

func (p *Player) AdjustSpeed(amount float64) {
	if p.Speed > -1.0 && p.Speed < 1.0 {
		p.Speed += amount
	}
}

func (p *Player) SetSpeed(amount float64) {
	p.Speed = amount
}

func (c *Chaser) SetSpeed(amount float64) {
	c.Speed = amount
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

func (p *Player) GetCollisionBox() *Box {
	b := &Box{p.Location.X - p.HalfW, p.Location.Y - p.HalfH,
		p.Bounds.W, p.Bounds.H}
	return b
}

func GetWalls() []*Box {
	return wall
}
