// Package drawstate provides a simple way of drawing the state of a 2D uint32 slice in a window.
package drawstate

import (
	"sync"
)

type config struct {
	windowWidth, windowHeight int
	stateWidth, stateHeight   int32
}

var drawChannel = make(chan []uint32)
var configChannel = make(chan config)
var drawChannelClosed bool
var drawChannelLock sync.Mutex

func init() {
	go worker()
}

func Open(windowWidth, windowHeight, stateWidth, stateHeight uint16) {
	configChannel <- config{
		windowWidth:  int(windowWidth),
		windowHeight: int(windowHeight),
		stateWidth:   int32(stateWidth),
		stateHeight:  int32(stateHeight),
	}
}

func Draw(state []uint32) {
	drawChannelLock.Lock()
	if !drawChannelClosed {
		drawChannel <- state
	}
	drawChannelLock.Unlock()
}

func Close() {
	drawChannelLock.Lock()
	if !drawChannelClosed {
		close(drawChannel)
		drawChannelClosed = true
	}
	drawChannelLock.Unlock()
}

func Closed() bool {
	drawChannelLock.Lock()
	var ret = drawChannelClosed
	drawChannelLock.Unlock()
	return ret
}
