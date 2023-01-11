package keys

import "github.com/holoplot/go-evdev"

func NewChord(keys ...evdev.EvCode) Chord {
	chord := make(Chord)

	for _, key := range keys {
		chord[key] = true
	}

	return chord
}

func (chord Chord) Keys() (keys []evdev.EvCode) {
	for k := range chord {
		keys = append(keys, k)
	}
	return
}

func (chord Chord) Add(key evdev.EvCode) {
	chord[key] = true
}

func (chord Chord) Remove(key evdev.EvCode) {
	delete(chord, key)
}
