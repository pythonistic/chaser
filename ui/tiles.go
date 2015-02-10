package ui

import (
	"encoding/json"
	"io/ioutil"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

// LoadTiles loads the tiles in the given JSDON filename, typically tiles.json.
func LoadTiles(filename string) {
	// load the sprite definitons
	tilesFiles := ParseTileDefinitions(filename)

	inited := img.Init(img.INIT_PNG)
	if inited&img.INIT_PNG != img.INIT_PNG {
		panic(img.GetError())
	}

	// load the sprites
	for _, spriteFile := range tilesFiles.Files {
		// we need to include sprite information:
		// frame size, frame id/count
		// image strip position and size
		rwOp := sdl.RWFromFile(spriteFile.Filename, "rb")

		if rwOp != nil {
			spriteSurface, err := img.LoadPNG_RW(rwOp)
			//		spriteSurface, err := img.Load(spriteFilename)
			if err != nil {
				panic(err)
			}
			spriteTexture, err := renderer.CreateTextureFromSurface(spriteSurface)
			if err != nil {
				panic(err)
			}

			for _, spriteDef := range spriteFile.Sprites {
				frameCount := len(spriteDef.Frames)
				newSprite := Sprite{name: spriteDef.Name, filename: spriteFile.Filename,
					frameCount: frameCount,
					rawFrame:   make([]*sdl.Surface, frameCount),
					frame:      make([]*sdl.Texture, frameCount),
					size:       make([]*sdl.Rect, frameCount),
					offsetX:    make([]int32, frameCount),
					offsetY:    make([]int32, frameCount)}
				for idx, frameDef := range spriteDef.Frames {
					newSprite.rawFrame[idx] = spriteSurface
					newSprite.frame[idx] = spriteTexture
					newSprite.size[idx] = &sdl.Rect{X: frameDef.X, Y: frameDef.Y,
						W: frameDef.W, H: frameDef.H}
					newSprite.offsetX[idx] = frameDef.W / -2
					newSprite.offsetY[idx] = frameDef.H / -2
				}
				sprites[spriteDef.Name] = &newSprite
			}
		} else {
			panic(spriteFile.Filename + " was nil")
		}

	}

	// cleanup
	img.Quit()
}

// ParseTileDefinitions parses a given filename and turns it into a TilesFile structure.
func ParseTileDefinitions(filename string) *TilesFile {
	var rawFile []byte
	var err error

	rawFile, err = ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var tilesFile TilesFile
	err = json.Unmarshal(rawFile, &tilesFile)
	if err != nil {
		panic(err)
	}

	return &tilesFile
}
