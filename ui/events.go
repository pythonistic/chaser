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
			/*			fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
						t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel) */
		case *sdl.MouseButtonEvent:
			fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
				t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			/*
			 * To check for click events:
			 *  1) record mouse button down location and time
			 *  2) if mouse button up location within 2 px of mouse button down and
			 *  3) if mouse buttown down time within 150 ms of mouse button up type == click
			 *  4) else type = drag
			 */
			switch t.Type {
			case sdl.MOUSEBUTTONUP:
				switch t.Button {
				case sdl.BUTTON_LEFT:
					state.SetClickLocation(state.Location{float64(t.X), float64(t.Y), 0.0},
						state.BEHAVIOR_ATTRACT)
				case sdl.BUTTON_RIGHT:
					state.SetClickLocation(state.Location{float64(t.X), float64(t.Y), 0.0},
						state.BEHAVIOR_AVOID)
				}
			}
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
