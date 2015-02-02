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
				make([]*sdl.Rect, 1, 1), make([]int32, 1, 1), make([]int32, 1, 1)}
			newSprite.rawFrame[0] = spriteSurface
			newSprite.frame[0], err = renderer.CreateTextureFromSurface(spriteSurface)
			if err != nil {
				panic(err)
			}
			newSprite.size[0] = &sdl.Rect{X: 0, Y: 0, W: spriteSurface.W, H: spriteSurface.H}
			newSprite.offsetX[0] = spriteSurface.W / -2
			newSprite.offsetY[0] = spriteSurface.H / -2
			sprites[spriteName] = &newSprite
		} else {
			panic(spriteFilename + " was nil")
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
