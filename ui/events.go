package ui

import (
	"chaser/state"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

func HandleEvents() bool {
	var event sdl.Event
	var keepRunning = true
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			keepRunning = false
		case *sdl.MouseMotionEvent:
			fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
				t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
		case *sdl.MouseButtonEvent:
			fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
				t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
		case *sdl.MouseWheelEvent:
			fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
				t.Timestamp, t.Type, t.Which, t.X, t.Y)
		case *sdl.KeyDownEvent:
			switch t.Keysym.Sym {
			case sdl.K_LEFT:
				state.GetPlayer().AdjustDirection(0.1)
			case sdl.K_RIGHT:
				state.GetPlayer().AdjustDirection(-0.1)
			case sdl.K_UP:
				state.GetPlayer().AdjustSpeed(0.1)
			case sdl.K_DOWN:
				state.GetPlayer().AdjustSpeed(0.1)
			}
		case *sdl.KeyUpEvent:
			fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			if t.Keysym.Sym == 'q' {
				keepRunning = false
			}
		}
	}

	return keepRunning
}
