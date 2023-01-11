package keys

import "github.com/holoplot/go-evdev"

type (
	Chord    map[evdev.EvCode]bool
	Sequence []Chord
)
