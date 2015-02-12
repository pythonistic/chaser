package state

import "fmt"

func makeRandomizedPrimsMaze(box *Box) [][]bool {
	grid := make([][]bool, box.H)
	for row := int32(0); row < box.H; row++ {
		grid[row] = make([]bool, box.W)
	}
	fmt.Println("dimensions", box.W, box.H)
	walls := make([]Point, box.H*box.W)
	maze := make(map[Point]Point)
	cell := Point{X: 0, Y: 0}
	grid, walls = primsAddCellToMaze(grid, walls, maze, cell)

	for len(walls) > 0 {
		idx := random(0, int32(len(walls)))

		// identify the grid cell opposite the current cell
		wall := walls[idx] // this is the maze wall
		walls = append(walls[:idx], walls[idx+1:]...)
		cell = maze[wall] // this is the maze floor adjacent to the wall
		delete(maze, wall)
		var opposite Point
		if wall.X < cell.X {
			opposite = Point{X: cell.X - 2, Y: cell.Y}
		} else if wall.X > cell.X {
			opposite = Point{X: cell.X + 2, Y: cell.Y}
		} else if wall.Y < cell.Y {
			opposite = Point{X: cell.X, Y: cell.Y - 2}
		} else {
			opposite = Point{X: cell.X, Y: cell.Y + 2}
		}
		if box.Contains(&opposite) {
			if !grid[opposite.Y][opposite.X] {
				grid[wall.Y][wall.X] = true
				grid, walls = primsAddCellToMaze(grid, walls, maze, opposite)
			}
		}
	}

	return grid
}

func primsAddCellToMaze(grid [][]bool, walls []Point, maze map[Point]Point, cell Point) ([][]bool, []Point) {
	grid[cell.Y][cell.X] = true
	if cell.X > 0 {
		pt := Point{X: cell.X - 1, Y: cell.Y}
		if !grid[pt.Y][pt.X] {
			walls = append(walls, pt)
			maze[pt] = cell
		}
	}
	if cell.X < int32(len(grid[0])-1) {
		pt := Point{X: cell.X + 1, Y: cell.Y}
		if !grid[pt.Y][pt.X] {
			walls = append(walls, pt)
			maze[pt] = cell
		}
	}
	if cell.Y > 0 {
		pt := Point{X: cell.X, Y: cell.Y - 1}
		if !grid[pt.Y][pt.X] {
			walls = append(walls, pt)
			maze[pt] = cell
		}
	}
	if cell.Y < int32(len(grid)-1) {
		pt := Point{X: cell.X, Y: cell.Y + 1}
		if !grid[pt.Y][pt.X] {
			walls = append(walls, pt)
			maze[pt] = cell
		}
	}
	fmt.Println(walls)
	return grid, walls
}

// The recursive maze subdivides the area (in box bounds) with walls
// The player width and height are fixed, although they should be derived
// from the sprite.
// This implementation has problems beyond two iterations:  the maze walls
// don't respect the already drawn or reserved areas.
func makeRecursiveMaze(box *Box) {
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
		makeRecursiveMaze(&Box{minX, minY, xPos - minX, maxY - minY, depth})
		makeRecursiveMaze(&Box{xPos, minY, maxX - xPos, maxY - minY, depth})
	} else if xPos < 0 {
		// horizontal wall
		fmt.Println("Horizonal wall case - yPos", yPos)
		makeHorizontalMazeWall(yPos, minX, maxX, playerWidth, playerHeight, depth)
		makeRecursiveMaze(&Box{minX, minY, maxX - minX, yPos - minY, depth})
		makeRecursiveMaze(&Box{minX, yPos, maxX - minX, maxY - yPos, depth})
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
		makeRecursiveMaze(&Box{minX, minY, xPos - minX, yPos - minY, depth}) // top left
		fmt.Println("make top right", Box{xPos, minY, maxX - xPos, yPos - minY, depth})
		makeRecursiveMaze(&Box{xPos, minY, maxX - xPos, yPos - minY, depth}) // top right
		fmt.Println("make bottom left", Box{minX, yPos, xPos - minX, maxY - yPos, depth})
		makeRecursiveMaze(&Box{minX, yPos, xPos - minX, maxY - yPos, depth}) // bottom left
		fmt.Println("make bottom right", Box{xPos, yPos, maxX - xPos, maxY - yPos, depth})
		makeRecursiveMaze(&Box{xPos, yPos, maxX - xPos, maxY - yPos, depth}) // bottom right
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
