package main

import (
	"chaser/state"
	"chaser/ui"
	"runtime"
	"time"
)

// ideal frames per second
const FRAMERATE int64 = 60

// ideal render + sleep time based on FPS
const FRAMERATE_SLEEP_NS int64 = int64(time.Second) / FRAMERATE

// sleep a minimum of 10 ms per loop
const MINIMUM_SLEEP_DURATION time.Duration = time.Duration(10 * 1000 * 1000)

// playfield size
var playfield state.Playfield = state.Playfield{800.0, 600.0}

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

func Init() {
	state.InitState(&playfield)
	var screen = ui.Screen{ui.SCREEN_ORIGIN_UNDEFINED, ui.SCREEN_ORIGIN_UNDEFINED,
		uint32(playfield.Width), uint32(playfield.Height), ui.COLOR_DEPTH_24, true, "Chaser"}
	ui.InitRenderer(&screen)
}

func Sleep() {
	var lastDuration time.Duration = time.Since(lastSleep)

	// we want to sleep to try to match 60 FPS
	var sleepDuration time.Duration = time.Duration(FRAMERATE_SLEEP_NS -
		lastDuration.Nanoseconds())

	if sleepDuration < MINIMUM_SLEEP_DURATION {
		sleepDuration = MINIMUM_SLEEP_DURATION
	}

	// capture the last sleep so we can adjust to match framerate
	lastSleep = time.Now()

	time.Sleep(sleepDuration)
}

func CleanUp() {
	ui.ShutdownRenderer()
}
