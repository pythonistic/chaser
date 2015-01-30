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
	wall = make([]*Box, 0, 100)
	opening = make([]*Box, 0, 100)

	makeMaze(&Box{0, 0, playfield.Width, playfield.Height, 0})
}

func makeMaze(box *Box) {
	depth := box.Z + 1
	fmt.Println("DEPTH ", depth)
	fmt.Println("Box ", box)
	if depth > 4 {
		return
	}

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
	var minimumWidth = playerWidth*2 + WallWidth
	var minumumHeight = playerHeight*2 + WallWidth

	// fail early if nothing to do
	if (box.W < minimumWidth &&
		box.H < minumumHeight) ||
		box.W+box.X > playfield.Width ||
		box.H+box.Y > playfield.Height {
		return
	}

	area = append(area, &Box{box.X, box.Y, box.W, box.H, depth})

	minX := box.X
	minY := box.Y
	maxX := box.W + minX
	maxY := box.H + minY
	var xPos, yPos int32 = -1, -1

	// make sure there's enough room to make a horizontal wall
	if box.W > minimumWidth {
		xPos = random(minX+playerHeight, maxX-playerHeight)
	}

	// make sure there's enough room to make a veritcal wall
	if box.H > minumumHeight {
		yPos = random(minY+playerWidth, maxY-playerWidth)
	}

	// check for intersection with an existing opening
	for _, open := range opening {
		if open.Contains(&Point{minX, yPos}) ||
			open.Contains(&Point{maxX, yPos}) {
			// fix the y pos
			if yPos-minY > maxY-yPos {
				// closer to bottom
				fmt.Println("A) adjusting yPos from", yPos, "to", (open.Y + open.H + 1))
				yPos = open.Y + open.H + 1
			} else {
				fmt.Println("B) adjusting yPos from", yPos, "to", (open.Y - 1))
				yPos = open.Y - 1
			}
			if yPos < minY || yPos > maxY {
				// area too small to really subdivide
				yPos = -1
			}
		}

		if open.Contains(&Point{xPos, minY}) ||
			open.Contains(&Point{xPos, maxY}) {
			// fix the x pos
			if xPos-minX > maxX-xPos {
				// closer to right
				fmt.Println("C) adjusting xPos from", xPos, "to", (open.X + open.W + 1))
				xPos = open.X + open.W + 1
			} else {
				fmt.Println("D) adjusting xPos from", xPos, "to", (open.X - 1))
				xPos = open.X - 1
			}
			if xPos < minX || xPos > maxX {
				// area too small to subdivide
				xPos = -1
			}
		}
	}

	// check for single wall cases
	if yPos < 0 {
		fmt.Println("vertical wall case - xPos", xPos)
		// vertical wall
		makeVerticalMazeWall(xPos, minY, maxY, playerWidth, playerHeight, depth)
		makeMaze(&Box{minX, minY, xPos - minX, maxY - minY, depth})
		makeMaze(&Box{xPos, minY, maxX - xPos, maxY - minY, depth})
	} else if xPos < 0 {
		// horizontal wall
		fmt.Println("Horizonal wall case - yPos", yPos)
		makeHorizontalMazeWall(yPos, minX, maxX, playerWidth, playerHeight, depth)
		makeMaze(&Box{minX, minY, maxX - minX, yPos - minY, depth})
		makeMaze(&Box{minX, yPos, maxX - minX, maxY - yPos, depth})
	} else {
		// two walls created, need to subdivide into four walls
		// skip one chamber randomly
		skip := random(1, 4)
		fmt.Println("skip ", skip, " xPos", xPos, " yPos", yPos)
		if skip != 1 {
			makeVerticalMazeWall(xPos, minY, yPos, playerWidth, playerHeight, depth)
		} else {
			wall = append(wall, &Box{xPos, minY, WallWidth, yPos, depth})
		}
		if skip != 2 {
			makeVerticalMazeWall(xPos, yPos, maxY, playerWidth, playerHeight, depth)
		} else {
			wall = append(wall, &Box{xPos, yPos, WallWidth, maxY, depth})
		}
		if skip != 3 {
			makeHorizontalMazeWall(yPos, minX, xPos, playerWidth, playerHeight, depth)
		} else {
			wall = append(wall, &Box{minX, yPos, xPos, WallWidth, depth})
		}
		if skip != 4 {
			makeHorizontalMazeWall(yPos, xPos, maxX, playerWidth, playerHeight, depth)
		} else {
			wall = append(wall, &Box{xPos, yPos, maxX, WallWidth, depth})
		}

		fmt.Println("make top left", Box{minX, minY, xPos - minX, yPos - minY, depth})
		makeMaze(&Box{minX, minY, xPos - minX, yPos - minY, depth}) // top left
		fmt.Println("make top right", Box{xPos, minY, maxX - xPos, yPos - minY, depth})
		makeMaze(&Box{xPos, minY, maxX - xPos, yPos - minY, depth}) // top right
		fmt.Println("make bottom left", Box{minX, yPos, xPos - minX, maxY - yPos, depth})
		makeMaze(&Box{minX, yPos, xPos - minX, maxY - yPos, depth}) // bottom left
		fmt.Println("make bottom right", Box{xPos, yPos, maxX - xPos, maxY - yPos, depth})
		makeMaze(&Box{xPos, yPos, maxX - xPos, maxY - yPos, depth}) // bottom right
	}
}

func makeVerticalMazeWall(xPos int32, minY int32, maxY int32, playerWidth int32, playerHeight int32, depth int32) {
	fmt.Println("makeVerticalMazeWall", xPos, minY, maxY, playerHeight)
	wallOpen := random(minY, maxY)
	if wallOpen > maxY-minY-playerHeight {
		// wall open is at the bottom
		fmt.Println("open bottom", Box{xPos, minY, WallWidth, maxY - minY - playerHeight, depth})
		wall = append(wall, &Box{xPos, minY, WallWidth, maxY - minY - playerHeight, depth})
		opening = append(opening, &Box{xPos - playerWidth/2, maxY - playerHeight, playerWidth, playerHeight, depth})
	} else if wallOpen < minY+playerHeight {
		// wall open is at the top
		fmt.Println("open top", Box{xPos, playerHeight + minY, WallWidth, maxY - minY - playerHeight, depth})
		wall = append(wall, &Box{xPos, playerHeight + minY, WallWidth, maxY - minY - playerHeight, depth})
		opening = append(opening, &Box{xPos - playerWidth/2, minY, playerWidth, playerHeight, depth})
	} else {
		// wall open is in the middle, make two walls
		fmt.Println("open midtop", Box{xPos, minY, WallWidth, wallOpen - minY, depth})
		fmt.Println("open midbot", Box{xPos, minY, WallWidth, wallOpen - minY, depth})
		wall = append(wall, &Box{xPos, minY, WallWidth, wallOpen - minY, depth})
		wall = append(wall, &Box{xPos, wallOpen + playerHeight, WallWidth, maxY - wallOpen - playerHeight, depth})
		opening = append(opening, &Box{xPos - playerWidth/2, wallOpen, playerWidth, playerHeight, depth})
	}
}

func makeHorizontalMazeWall(yPos int32, minX int32, maxX int32, playerWidth int32, playerHeight int32, depth int32) {
	wallOpen := random(minX, maxX)
	if wallOpen > maxX-minX-playerWidth {
		// wall open is at the right
		wall = append(wall, &Box{minX, yPos, maxX - minX - playerWidth, WallWidth, depth})
		opening = append(opening, &Box{maxX - playerWidth, yPos - playerHeight/2, playerWidth, playerHeight, depth})
	} else if wallOpen < minX+playerWidth {
		// wall open is at the left
		wall = append(wall, &Box{playerWidth + minX, yPos, maxX - minX - playerWidth, WallWidth, depth})
		opening = append(opening, &Box{minX, yPos - playerHeight/2, playerWidth, playerHeight, depth})
	} else {
		// wall open is in the middle, make two walls
		wall = append(wall, &Box{minX, yPos, wallOpen - minX, WallWidth, depth})
		wall = append(wall, &Box{wallOpen + playerWidth, yPos, maxX - wallOpen - playerWidth, WallWidth, depth})
		opening = append(opening, &Box{wallOpen, yPos - playerHeight/2, playerWidth, playerHeight, depth})
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
