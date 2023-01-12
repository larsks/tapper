package keys

import (
	"fmt"
	"strings"

	"github.com/holoplot/go-evdev"
	"golang.org/x/exp/slices"
)

func NewChord(keys ...evdev.EvCode) *Chord {
	chord := Chord{}

	for _, key := range keys {
		chord.keys = append(chord.keys, key)
	}

	return &chord
}

// A Chord is really an ordered set -- that is, we want to send keys
// events in the order in which they are added to the chord, but we don't
// want to allow duplicate events. Rather than trying to maintain a unique
// list, we uniq-ify the list in this Keys() method.
func (chord *Chord) Keys() (keys []evdev.EvCode) {
	uniqueKeys := make(map[evdev.EvCode]bool)

	for _, k := range chord.keys {
		exists := uniqueKeys[k]

		if !exists {
			keys = append(keys, k)
			uniqueKeys[k] = true
		}
	}

	return
}

func (chord *Chord) Equal(other *Chord) bool {
	return slices.Equal(chord.Keys(), other.Keys())
}

func (chord *Chord) Add(key evdev.EvCode) {
	chord.keys = append(chord.keys, key)
}

func (chord *Chord) Update(ref *Chord) {
	for _, key := range ref.Keys() {
		chord.Add(key)
	}
}

func (chord *Chord) Remove(deleteKey evdev.EvCode) {
	newKeys := []evdev.EvCode{}

	for _, key := range chord.Keys() {
		if key != deleteKey {
			newKeys = append(newKeys, key)
		}
	}

	chord.keys = newKeys
}

func (chord *Chord) String() string {
	var keys []string

	for _, key := range chord.Keys() {
		keys = append(keys, evdev.KEYNames[key])
	}

	return strings.Join(keys, ":")
}

func NewSequence(chords ...*Chord) *Sequence {
	seq := Sequence{}

	for _, chord := range chords {
		seq.Add(chord)
	}

	return &seq
}

func (seq *Sequence) Chords() []*Chord {
	return seq.chords
}

func (seq *Sequence) String() string {
	var chords []string

	for _, chord := range seq.chords {
		chords = append(chords, chord.String())
	}

	return strings.Join(chords, " ")
}

func (seq *Sequence) Add(chord *Chord) {
	seq.chords = append(seq.chords, chord)
}

func ChordFromString(s string) (*Chord, error) {
	keys := NewChord()

	for _, tok := range strings.Split(s, ":") {
		tok = fmt.Sprintf("KEY_%s", strings.ToUpper(tok))
		code, ok := evdev.KEYFromString[tok]
		if !ok {
			return keys, fmt.Errorf("%s: unknown key name", tok)
		}

		keys.Add(code)
	}

	return keys, nil
}

func SequenceFromString(s string) (*Sequence, error) {
	seq := Sequence{}

	for _, tok := range strings.Fields(s) {
		chord, err := ChordFromString(tok)
		if err != nil {
			return nil, err
		}

		seq.chords = append(seq.chords, chord)
	}

	return &seq, nil
}
