package ui

import (
	"chaser/state"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"os"
	"strconv"
	"time"
)

type Sprite struct {
	name     string
	filename string
	// TODO split frames into types
	// TODO consider overlays on sprites (clothes, equipment, buffs)
	frameCount uint8
	rawFrame   []*sdl.Surface
	frame      []*sdl.Texture
	size       []*sdl.Rect
	offsetX    []float64 // offset is because sprites render center at location
	offsetY    []float64
}

const (
	COLOR_DEPTH_8           uint8  = 8
	COLOR_DEPTH_15          uint8  = 15
	COLOR_DEPTH_16          uint8  = 16
	COLOR_DEPTH_24          uint8  = 24
	SCREEN_ORIGIN_UNDEFINED uint32 = 0xFFFFFFFF
)

// color definitions
var (
	BLACK  sdl.Color = sdl.Color{0, 0, 0, 255}
	GREEN  sdl.Color = sdl.Color{0, 255, 0, 255}
	RED    sdl.Color = sdl.Color{255, 0, 0, 255}
	YELLOW sdl.Color = sdl.Color{255, 255, 0, 255}
)

// window definitions
var (
	window   *sdl.Window
	surface  *sdl.Surface
	renderer *sdl.Renderer
)

// FPS tracking
var (
	lastFpsTime time.Time = time.Now()
	frameCount  uint32    = 0
	currentFps  string    = "0"
	sansFont16  *ttf.Font
	renderFps   bool = true
	spriteDefs       = []string{"player", "resources/PlanetCute PNG/Character Pink Girl.png",
		"chaser", "resources/PlanetCute PNG/Heart.png"}
	sprites = make(map[string]*Sprite)
)

// ==== Public interface

// struct to make a portable screen definition for callers
type Screen struct {
	OriginX  uint32
	OriginY  uint32
	Width    uint32
	Height   uint32
	Depth    uint8 // not used
	Windowed bool
	Title    string
}

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
	inited := img.Init(img.INIT_PNG)
	if inited&img.INIT_PNG != img.INIT_PNG {
		panic(img.GetError())
	}

	// load the sprites
	for idx := 0; idx < len(spriteDefs); idx += 2 {
		spriteName := spriteDefs[idx]
		spriteFilename := spriteDefs[idx+1]

		// TODO at this time, all the sprites are single image
		// we'd need to include sprite information:
		// frame size, frame id/count
		// image strip position and size
		rwOp := sdl.RWFromFile(spriteFilename, "rb")

		if rwOp != nil {
			spriteSurface, err := img.LoadPNG_RW(rwOp)
			//		spriteSurface, err := img.Load(spriteFilename)
			if err != nil {
				panic(err)
			}

			newSprite := Sprite{spriteName, spriteFilename, 1,
				make([]*sdl.Surface, 1, 1), make([]*sdl.Texture, 1, 1),
				make([]*sdl.Rect, 1, 1), make([]float64, 1, 1), make([]float64, 1, 1)}
			newSprite.rawFrame[0] = spriteSurface
			newSprite.frame[0], err = renderer.CreateTextureFromSurface(spriteSurface)
			if err != nil {
				panic(err)
			}
			newSprite.size[0] = &sdl.Rect{0, 0, spriteSurface.W, spriteSurface.H}
			newSprite.offsetX[0] = float64(spriteSurface.W) / -2.0
			newSprite.offsetY[0] = float64(spriteSurface.H) / -2.0
			sprites[spriteName] = &newSprite
		} else {
			panic(spriteFilename + " was nil")
		}
	}

	// cleanup
	img.Quit()
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
	rect := &sdl.Rect{int32(sprite.offsetX[0] + location.X),
		int32(sprite.offsetY[0] + location.Y),
		sprite.size[0].W, sprite.size[0].H}
	err := renderer.Copy(sprite.frame[0], nil, rect)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to copy texture to renderer: %s\n", err)
		panic(err)
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
