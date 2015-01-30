package ui

import (
	"chaser/state"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

// window definitions
var (
	window   *sdl.Window
	surface  *sdl.Surface
	renderer *sdl.Renderer
)

// FPS tracking
var (
	lastFpsTime = time.Now()
	frameCount  = 0
	currentFps  = "0"
	sansFont16  *ttf.Font
	renderFps   = true
	spriteDefs  = []string{"player", "resources/PlanetCute PNG/Character Pink Girl.png",
		"chaser", "resources/PlanetCute PNG/Heart.png"}
	sprites = make(map[string]*Sprite)
)

// ==== Public interface

func InitRenderer(screen *Screen) {
	initScreen(screen)
	initText()
	loadSprites()
}

func UpdateScreen() {
	var err error

	calculateFps()

	// clear the screen
	//err = renderer.SetDrawColor(BLACK.R, BLACK.G, BLACK.B, BLACK.A)
	err = renderer.SetDrawColor(127, 127, 127, 255)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set renderer draw color to black: %s\n", err)
		panic(err)
	}
	err = renderer.Clear()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to clear screen: %s\n", err)
		panic(err)
	}

	var rendererViewport *sdl.Rect = &sdl.Rect{0, 0, 0, 0}
	renderer.GetViewport(rendererViewport)

	if renderFps {
		renderFpsCounter(rendererViewport)
	}

	renderSolidAreas()
	renderWalls()
	renderChaser()
	renderPlayer()

	renderer.Present()
}

func ShutdownRenderer() {
	shutdownText()
	shutdownScreen()
}

// ====== Implementations

func initScreen(screen *Screen) {
	var originX int = sdl.WINDOWPOS_UNDEFINED
	var originY int = sdl.WINDOWPOS_UNDEFINED
	var flags uint32 = sdl.WINDOW_SHOWN

	if screen.OriginX != SCREEN_ORIGIN_UNDEFINED {
		originX = int(screen.OriginX)
	}

	if screen.OriginY != SCREEN_ORIGIN_UNDEFINED {
		originY = int(screen.OriginY)
	}

	if !screen.Windowed {
		flags = flags | sdl.WINDOW_FULLSCREEN
	}

	var err error
	window, err = sdl.CreateWindow(screen.Title, originX, originY,
		int(screen.Width), int(screen.Height), flags)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	err = renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	surface = window.GetSurface()
}

func initText() {
	success := ttf.Init()
	if success == -1 {
		panic("SDL TTF init failed")
	}

	var err error
	sansFont16, err = ttf.OpenFont("resources/ubuntu.ttf", 16)
	if err != nil {
		panic(err)
	}
	// 1 px outline
	sansFont16.SetOutline(1)
}

func loadSprites() {
	LoadTiles("tiles.json")

	// set the Player and Chaser bounds
	playerSprite := sprites["player"].size[0]
	chaserSprite := sprites["chaser"].size[0]
	state.GetPlayer().Bounds = &state.Box{playerSprite.X,
		playerSprite.Y, playerSprite.W / 2,
		playerSprite.H / 2, 0}
	state.GetPlayer().HalfW = playerSprite.W / 4
	state.GetPlayer().HalfH = playerSprite.H / 4 // the image is huge!
	state.GetChaser().Bounds = &state.Box{chaserSprite.X,
		chaserSprite.Y, chaserSprite.W,
		chaserSprite.H, 0}
}

func calculateFps() {
	var lastFpsDuration time.Duration = time.Since(lastFpsTime)
	frameCount += 1
	if lastFpsDuration > time.Second {
		lastFpsTime = time.Now()
		currentFps = strconv.Itoa(int(frameCount))
		fmt.Println(time.Now(), "FPS", currentFps) // TODO remove the following line when the on-screen rendering works
		frameCount = 0
	}
}

func renderFpsCounter(rendererViewport *sdl.Rect) {
	var err error

	// write the FPS counter
	fpsSurface := sansFont16.RenderText_Solid(currentFps, YELLOW)
	//fpsSurface := sansFont16.RenderText_Solid("Owen", YELLOW)
	fpsRect := sdl.Rect{rendererViewport.W - 50, 5, fpsSurface.W, fpsSurface.H}
	// fpsRect := sdl.Rect{10, 10, fpsSurface.W, fpsSurface.H}
	var fpsTexture *sdl.Texture
	fpsTexture, err = renderer.CreateTextureFromSurface(fpsSurface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture from surface: %s\n", err)
		panic(err)
	}

	err = renderer.Copy(fpsTexture, nil, &fpsRect)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to copy texture to renderer: %s\n", err)
		panic(err)
	}
	fpsTexture.Destroy()
	fpsSurface.Free()
}

func renderPlayer() {
	player := state.GetPlayer()

	bounds := player.GetCollisionBox()
	rect := &sdl.Rect{bounds.X, bounds.Y, bounds.W, bounds.H}
	renderer.SetDrawColor(YELLOW.R, YELLOW.G, YELLOW.B, YELLOW.A)
	renderer.DrawRect(rect)
	//renderTriangle(player.Location, GREEN)
	renderSprite(sprites["player"], player.Location)
}

func renderChaser() {
	chaser := state.GetChaser()
	//renderTriangle(chaser.Location, RED)
	renderSprite(sprites["chaser"], chaser.Location)
}

func renderTriangle(o state.Location, c sdl.Color) {
	renderer.SetDrawColor(c.R, c.G, c.B, c.A)
	points := []sdl.Point{{int32(o.X) - 2, int32(o.Y) + 2},
		{int32(o.X) + 2, int32(o.Y) + 2}, {int32(o.X), int32(o.Y) - 2},
		{int32(o.X) - 2, int32(o.Y) + 2}}
	renderer.DrawLines(points)
}

func renderSprite(sprite *Sprite, location state.Location) {
	rect := &sdl.Rect{sprite.offsetX[0] + location.X,
		sprite.offsetY[0] + location.Y,
		sprite.size[0].W, sprite.size[0].H}
	err := renderer.Copy(sprite.frame[0], nil, rect)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to copy texture to renderer: %s\n", err)
		panic(err)
	}
}

func renderWalls() {
	walls := state.GetWalls()
	for i := 0; i < len(walls); i++ {
		wall := walls[i]
		rect := &sdl.Rect{wall.X, wall.Y, wall.W, wall.H}
		renderer.SetDrawColor(RED.R, RED.G, RED.B, RED.A)
		switch wall.Z {
		case 0:
			renderer.SetDrawColor(WHITE.R, WHITE.G, WHITE.B, WHITE.A)
		case 1:
			renderer.SetDrawColor(BLUE.R, BLUE.G, BLUE.B, BLUE.A)
		case 2:
			renderer.SetDrawColor(GREEN.R, GREEN.G, GREEN.B, GREEN.A)
		case 3:
			renderer.SetDrawColor(YELLOW.R, YELLOW.G, YELLOW.B, YELLOW.A)
		case 4:
			renderer.SetDrawColor(PURPLE.R, PURPLE.G, PURPLE.B, PURPLE.A)
		case 5:
			renderer.SetDrawColor(BLACK.R, BLACK.G, BLACK.B, BLACK.A)
		}
		err := renderer.DrawRect(rect)
		if err != nil {
			panic(err)
		}
	}
}

func renderSolidAreas() {
	areas := state.GetAreas()
	for _, area := range areas {
		rect := &sdl.Rect{area.X, area.Y, area.W, area.H}
		if area.Z == 3 {

			renderer.SetDrawColor(RED.R, RED.G, RED.B, 32)
			switch area.Z {
			case 0:
				renderer.SetDrawColor(WHITE.R, WHITE.G, WHITE.B, 32)
			case 1:
				renderer.SetDrawColor(BLUE.R, BLUE.G, BLUE.B, 32)
			case 2:
				renderer.SetDrawColor(GREEN.R, GREEN.G, GREEN.B, 32)
			case 3:
				renderer.SetDrawColor(YELLOW.R, YELLOW.G, YELLOW.B, 32)
			case 4:
				renderer.SetDrawColor(PURPLE.R, PURPLE.G, PURPLE.B, 32)
			case 5:
				renderer.SetDrawColor(BLACK.R, BLACK.G, BLACK.B, 32)
			}
			err := renderer.FillRect(rect)
			if err != nil {
				panic(err)
			}
			renderer.SetDrawColor(0, 0, 0, 128)
			renderer.DrawRect(rect)
		}

		for _, opening := range state.GetOpenings() {
			rect := &sdl.Rect{opening.X, opening.Y, opening.W, opening.H}
			renderer.SetDrawColor(WHITE.R, WHITE.G, WHITE.B, 64)
			renderer.DrawRect(rect)
		}

	}
}

func shutdownScreen() {
	window.Destroy()
	sdl.Quit()
}

func shutdownText() {
	sansFont16.Close()
	ttf.Quit()
}
