package ui

import (
	"testing"
)

func TestParseTileDefintions(t *testing.T) {
	filename := "../resources/tiles.json"
	tilesFile := ParseTileDefinitions(filename)
	if tilesFile == nil {
		t.Error("tilesFile was nil")
	}
	if len(tilesFile.Files) == 0 {
		t.Error("no file structures parsed")
	}
	if tilesFile.Files[0].Filename != "resources/tilecrusader-art/characters-32x32.png" {
		t.Error("expected filename incorrect, not characters-32x32.png but was", tilesFile.Files[0].Filename)
	}
	if 10 != len(tilesFile.Files[0].Sprites) {
		t.Error("expected 10 sprites but found", len(tilesFile.Files[0].Sprites))
	}
	if "player" != tilesFile.Files[0].Sprites[0].Name {
		t.Error("expected sprite name player but found", tilesFile.Files[0].Sprites[0].Name)
	}
	if 0 != tilesFile.Files[0].Sprites[0].Frames[0].X {
		t.Error("expected frame 0 X coordinate to be 0 but was", tilesFile.Files[0].Sprites[0].Frames[0].X)
	}
	if 0 != tilesFile.Files[0].Sprites[0].Frames[0].Y {
		t.Error("expected frame 0 Y coordinate to be 0 but was", tilesFile.Files[0].Sprites[0].Frames[0].Y)
	}
	if 32 != tilesFile.Files[0].Sprites[0].Frames[0].W {
		t.Error("expected frame 0 W to be 32 but was", tilesFile.Files[0].Sprites[0].Frames[0].W)
	}
	if 32 != tilesFile.Files[0].Sprites[0].Frames[0].H {
		t.Error("expected frame 0 H to be 32 but was", tilesFile.Files[0].Sprites[0].Frames[0].H)
	}
}
