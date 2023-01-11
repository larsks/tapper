package keys

import "github.com/holoplot/go-evdev"

type (
	Chord struct {
		keys []evdev.EvCode
	}

	Sequence struct {
		chords []*Chord
	}
)
