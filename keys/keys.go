package keys

import (
	"fmt"
	"strings"

	"github.com/holoplot/go-evdev"
)

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

func (chord Chord) Update(ref Chord) {
	for key := range ref {
		chord.Add(key)
	}
}

func (chord Chord) Remove(key evdev.EvCode) {
	delete(chord, key)
}

func (chord Chord) String() string {
	var keys []string

	for key := range chord {
		keys = append(keys, evdev.KEYNames[key])
	}

	return strings.Join(keys, ":")
}

func (seq Sequence) String() string {
	var chords []string

	for _, chord := range seq {
		chords = append(chords, chord.String())
	}

	return strings.Join(chords, " ")
}

func (seq Sequence) Add(chord Chord) {
	seq = append(seq, chord)
}

func ChordFromString(s string) (Chord, error) {
	keys := NewChord()

	for _, tok := range strings.Split(s, ":") {
		tok = fmt.Sprintf("KEY_%s", strings.ToUpper(tok))
		code, ok := evdev.KEYFromString[tok]
		if !ok {
			return keys, fmt.Errorf("%s: unknown key name", tok)
		}

		keys[code] = true
	}

	return keys, nil
}

func SequenceFromString(s string) (Sequence, error) {
	seq := Sequence{}

	for _, tok := range strings.Fields(s) {
		chord, err := ChordFromString(tok)
		if err != nil {
			return nil, err
		}

		seq = append(seq, chord)
	}

	return seq, nil
}
