package drawstate_test

import "time"
import "drawstate"

func Example() {
	// Open 512x512 window. Set number of states to 128x128.
	drawstate.Open(512, 512, 128, 128)
	defer drawstate.Close()

	// Allocate states slice (all states initialized to 0)
	var state = make([]uint32, 128*128)

	// Draw state and wait 2 seconds
	drawstate.Draw(state)
	time.Sleep(time.Second)

	// Set upper half of states to 1
	for i := range state[0 : len(state)/2] {
		state[i] = 1
	}

	// Draw changed state and wait 2 seconds
	drawstate.Draw(state)
	time.Sleep(time.Second)
}
