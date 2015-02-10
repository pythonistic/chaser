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

// InitRenderer initializes the game rendering structures and window.
func InitRenderer(screen *Screen) {
	initScreen(screen)
	initText()
	loadSprites()
}

// UpdateScreen is called once per frame or tick to redraw the screen.
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

	var rendererViewport = &sdl.Rect{X: 0, Y: 0, W: 0, H: 0}
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

// ShutdownRenderer is called to release window and UI resources when the game is exiting.
func ShutdownRenderer() {
	shutdownText()
	shutdownScreen()
}

// ====== Implementations

func initScreen(screen *Screen) {
	var originX = sdl.WINDOWPOS_UNDEFINED
	var originY = sdl.WINDOWPOS_UNDEFINED
	var flags uint32 = sdl.WINDOW_SHOWN

	if screen.OriginX != ScreenOriginUndefined {
		originX = int(screen.OriginX)
	}

	if screen.OriginY != ScreenOriginUndefined {
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
	LoadTiles("resources/tiles.json")

	fmt.Println("sprites", sprites["player"])

	// set the Player and Chaser bounds
	playerSprite := sprites["player"].size[0]
	chaserSprite := sprites["orc"].size[0]
	state.GetPlayer().Bounds = &state.Box{X: playerSprite.X,
		Y: playerSprite.Y, W: playerSprite.W / 2,
		H: playerSprite.H / 2, Z: 0}
	state.GetPlayer().HalfW = playerSprite.W / 4
	state.GetPlayer().HalfH = playerSprite.H / 4 // the image is huge!
	state.GetChaser().Bounds = &state.Box{X: chaserSprite.X,
		Y: chaserSprite.Y, W: chaserSprite.W,
		H: chaserSprite.H, Z: 0}
}

func calculateFps() {
	var lastFpsDuration = time.Since(lastFpsTime)
	frameCount++
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
	fpsSurface := sansFont16.RenderText_Solid(currentFps, Yellow)
	//fpsSurface := sansFont16.RenderText_Solid("Owen", YELLOW)
	fpsRect := sdl.Rect{X: rendererViewport.W - 50, Y: 5, W: fpsSurface.W, H: fpsSurface.H}
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
	rect := &sdl.Rect{X: bounds.X, Y: bounds.Y, W: bounds.W, H: bounds.H}
	renderer.SetDrawColor(Yellow.R, Yellow.G, Yellow.B, Yellow.A)
	renderer.DrawRect(rect)
	//renderTriangle(player.Location, GREEN)
	renderSprite(sprites["player"], player.Location)
}

func renderChaser() {
	chaser := state.GetChaser()
	//renderTriangle(chaser.Location, RED)
	renderSprite(sprites["orc"], chaser.Location)
}

func renderTriangle(o state.Location, c sdl.Color) {
	renderer.SetDrawColor(c.R, c.G, c.B, c.A)
	points := []sdl.Point{{int32(o.X) - 2, int32(o.Y) + 2},
		{int32(o.X) + 2, int32(o.Y) + 2}, {int32(o.X), int32(o.Y) - 2},
		{int32(o.X) - 2, int32(o.Y) + 2}}
	renderer.DrawLines(points)
}

func renderSprite(sprite *Sprite, location state.Location) {
	// TODO also include a frame indicator
	frame := 0
	rect := &sdl.Rect{X: sprite.offsetX[frame] + location.X,
		Y: sprite.offsetY[frame] + location.Y,
		W: sprite.size[frame].W, H: sprite.size[frame].H}
	// use nil for the source rect to copy the whole texture
	// use nil for the dest to paint the whole renderer
	// here, we need to specify the source rect based on the frame
	srcRect := sprite.size[frame]
	err := renderer.Copy(sprite.frame[frame], srcRect, rect)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to copy texture to renderer: %s\n", err)
		panic(err)
	}
}

func renderWalls() {
	walls := state.GetWalls()
	for i := 0; i < len(walls); i++ {
		wall := walls[i]
		rect := &sdl.Rect{X: wall.X, Y: wall.Y, W: wall.W, H: wall.H}
		renderer.SetDrawColor(Red.R, Red.G, Red.B, Red.A)
		switch wall.Z {
		case 0:
			renderer.SetDrawColor(White.R, White.G, White.B, White.A)
		case 1:
			renderer.SetDrawColor(Blue.R, Blue.G, Blue.B, Blue.A)
		case 2:
			renderer.SetDrawColor(Green.R, Green.G, Green.B, Green.A)
		case 3:
			renderer.SetDrawColor(Yellow.R, Yellow.G, Yellow.B, Yellow.A)
		case 4:
			renderer.SetDrawColor(Purple.R, Purple.G, Purple.B, Purple.A)
		case 5:
			renderer.SetDrawColor(Black.R, Black.G, Black.B, Black.A)
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
		rect := &sdl.Rect{X: area.X, Y: area.Y, W: area.W, H: area.H}
		if area.Z == 3 {

			renderer.SetDrawColor(Red.R, Red.G, Red.B, Red.A)
			switch area.Z {
			case 0:
				renderer.SetDrawColor(White.R, White.G, White.B, White.A)
			case 1:
				renderer.SetDrawColor(Blue.R, Blue.G, Blue.B, Blue.A)
			case 2:
				renderer.SetDrawColor(Green.R, Green.G, Green.B, Green.A)
			case 3:
				renderer.SetDrawColor(Yellow.R, Yellow.G, Yellow.B, Yellow.A)
			case 4:
				renderer.SetDrawColor(Purple.R, Purple.G, Purple.B, Purple.A)
			case 5:
				renderer.SetDrawColor(Black.R, Black.G, Black.B, Black.A)
			}
			err := renderer.FillRect(rect)
			if err != nil {
				panic(err)
			}
			renderer.SetDrawColor(0, 0, 0, 128)
			renderer.DrawRect(rect)
		}

		for _, opening := range state.GetOpenings() {
			rect := &sdl.Rect{X: opening.X, Y: opening.Y, W: opening.W, H: opening.H}
			renderer.SetDrawColor(White.R, White.G, White.B, 64)
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
