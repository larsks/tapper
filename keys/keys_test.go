package keys

import (
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/assert"
)

func TestChordUnique(t *testing.T) {
	chord := NewChord()
	chord.Add(evdev.KEYFromString["KEY_A"])
	chord.Add(evdev.KEYFromString["KEY_A"])

	assert.Len(t, chord.Keys(), 1)
}

func TestChordFromStringSucceeds(t *testing.T) {
	chord, err := ChordFromString("a")
	assert.Nil(t, err)
	assert.ElementsMatch(t, []evdev.EvCode{evdev.KEY_A}, chord.Keys())
}

func TestChordFromStringFails(t *testing.T) {
	chord, err := ChordFromString("does_not_exist")
	assert.NotNil(t, err)
	assert.Len(t, chord.Keys(), 0)
}

func TestSequenceFromString(t *testing.T) {
	expected := []*Chord{
		NewChord(evdev.KEY_LEFTSHIFT, evdev.KEY_T),
		NewChord(evdev.KEY_A),
		NewChord(evdev.KEY_P),
		NewChord(evdev.KEY_P),
		NewChord(evdev.KEY_E),
		NewChord(evdev.KEY_RIGHTSHIFT, evdev.KEY_R),
	}

	seq, err := SequenceFromString("leftshift:t a p p e rightshift:r")
	assert.Nil(t, err)
	assert.Len(t, seq.Chords(), 6)
	assert.ElementsMatch(t, expected, seq.Chords())
}

func TestChordRemoveKey(t *testing.T) {
	chord, err := ChordFromString("a:b:c")
	assert.Nil(t, err)
	chord.Remove(evdev.KEY_B)
	assert.Len(t, chord.Keys(), 2)
	assert.ElementsMatch(t, []evdev.EvCode{evdev.KEY_A, evdev.KEY_C}, chord.Keys())
}

func TestChordUpdate(t *testing.T) {
	chord1 := NewChord(evdev.KEY_A)
	chord2 := NewChord(evdev.KEY_B, evdev.KEY_C)

	chord1.Update(chord2)

	expected := []evdev.EvCode{evdev.KEY_A, evdev.KEY_B, evdev.KEY_C}
	assert.Len(t, chord1.Keys(), 3)
	assert.ElementsMatch(t, expected, chord1.Keys())
}
