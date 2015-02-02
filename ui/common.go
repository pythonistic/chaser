package ui

import "github.com/veandco/go-sdl2/sdl"

// SpriteFrame represents a single frame of a sprite definition.
type SpriteFrame struct {
	X int32
	Y int32
	W int32
	H int32
}

// SpriteDef represents a sprite definition as loaded from tiles.json.
type SpriteDef struct {
	Name   string
	Frames []*SpriteFrame
}

// SpriteFile represents the sprites in a single graphics file as loaded from tiles.json.
type SpriteFile struct {
	Filename string
	Sprites  []SpriteDef
}

// TilesFile represents the sprite definitions in tiles.json.
type TilesFile struct {
	Files []SpriteFile
}

// Sprite represents a renderable, animated sprite.
type Sprite struct {
	name     string
	filename string
	// TODO split frames into types
	// TODO consider overlays on sprites (clothes, equipment, buffs)
	frameCount uint8
	rawFrame   []*sdl.Surface
	frame      []*sdl.Texture
	size       []*sdl.Rect
	offsetX    []int32 // offset is because sprites render center at location
	offsetY    []int32
}

// Screen provides a portable screen definition for callers.
type Screen struct {
	OriginX  int32
	OriginY  int32
	Width    int32
	Height   int32
	Depth    uint8 // not used
	Windowed bool
	Title    string
}

const (
	// ColorDepth8 indicates 8-bit color.
	ColorDepth8 = 8
	// ColorDepth15 indicates 15-bit color.
	ColorDepth15 = 15
	// ColorDepth16 indicates 16-bit color.
	ColorDepth16 = 16
	// ColorDepth24 indicates 24-bit color.
	ColorDepth24 = 24
	// ScreenOriginUndefined is a magic value that the program to render the game window in the center.
	ScreenOriginUndefined = 0xFFFFFF
)

// color definitions
var (
	Black  = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	Green  = sdl.Color{R: 0, G: 255, B: 0, A: 255}
	Red    = sdl.Color{R: 255, G: 0, B: 0, A: 255}
	Yellow = sdl.Color{R: 255, G: 255, B: 0, A: 255}
	White  = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	Blue   = sdl.Color{R: 0, G: 0, B: 255, A: 255}
	Purple = sdl.Color{R: 255, G: 0, B: 255, A: 255}
)
