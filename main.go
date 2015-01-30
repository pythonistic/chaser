package main

import (
	"chaser/state"
	"chaser/ui"
	"runtime"
	"time"
)

// ideal frames per second
const Framerate int64 = 60

// ideal render + sleep time based on FPS
const FramerateSleepNs = int64(time.Second) / Framerate

// sleep a minimum of 10 ms per loop
const MinimumSleepDuration time.Duration = time.Duration(10 * 1000 * 1000)

// playfield size
var playfield = state.Playfield{Width: 1024, Height: 768}

var lastSleep time.Time

func main() {
	runtime.LockOSThread()

	Init()
	var running = true
	lastSleep = time.Now()

	for running {
		ui.UpdateScreen() // game timer lives here
		running = ui.HandleEvents()
		state.UpdateState()
		Sleep()
	}

	CleanUp()
}

// Init intializes the game state in the correct order.
func Init() {
	state.InitState(&playfield)
	var screen = ui.Screen{OriginX: ui.ScreenOriginUndefined,
		OriginY: ui.ScreenOriginUndefined, Width: playfield.Width,
		Height: playfield.Height, Depth: ui.ColorDepth24, Windowed: true,
		Title: "Chaser"}
	ui.InitRenderer(&screen)
}

// Sleep calculates how long to sleep between frames and pauses the program.
func Sleep() {
	var lastDuration = time.Since(lastSleep)

	// we want to sleep to try to match 60 FPS
	var sleepDuration = time.Duration(FramerateSleepNs -
		lastDuration.Nanoseconds())

	if sleepDuration < MinimumSleepDuration {
		sleepDuration = MinimumSleepDuration
	}

	// capture the last sleep so we can adjust to match framerate
	lastSleep = time.Now()

	time.Sleep(sleepDuration)
}

// CleanUp cleans up as the game is shutting down.
func CleanUp() {
	ui.ShutdownRenderer()
}
