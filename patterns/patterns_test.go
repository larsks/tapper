package patterns

import (
	"fmt"
	"testing"

	evdev "github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/assert"
)

type (
	SequenceMap map[string]Sequence
	Expected    struct {
		name  string
		found bool
		more  bool
	}
)

var sequenceMap = SequenceMap{
	"ls": {
		Chord{evdev.KEY_LEFTSHIFT: true},
	},
	"rs": {
		Chord{evdev.KEY_RIGHTSHIFT: true},
	},
	"lsls": {
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true},
	},
	"lslsls": {
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true},
	},
	"lsrsls": {
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true},
	},
	"rsrsrs": {
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
	},
	"rsls": {
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true},
	},
	"lsrs": {
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
	},
	"lsrsrs": {
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
	},
	"rsrs": {
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
	},
	"chord:lsrs": {
		Chord{evdev.KEY_LEFTSHIFT: true, evdev.KEY_RIGHTSHIFT: true},
	},
	"chord:lsrs:lsrs": {
		Chord{evdev.KEY_LEFTSHIFT: true, evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true, evdev.KEY_RIGHTSHIFT: true},
	},
}

func createPatterns(names ...string) *Patterns {
	p := NewPatterns()

	for _, name := range names {
		seq, ok := sequenceMap[name]
		if !ok {
			panic("unknown sequence name")
		}

		p.AddSequence(seq, []string{name})
	}

	return p
}

func checkResults(t *testing.T, patterns *Patterns, expected []Expected) {
	for _, check := range expected {
		seq := sequenceMap[check.name]
		node, found, more := patterns.FindSequence(seq)

		/*
			fmt.Printf("---\nsequence: %s\nhave: found %t more %t\nwant: found %t have %t\n",
				check.name, found, more, check.found, check.more,
			)
		*/

		assert.Equal(t, check.found, found, fmt.Sprintf("sequence %s field found", check.name))
		assert.Equal(t, check.more, more, fmt.Sprintf("sequence %s field more", check.name))

		if check.found {
			assert.Equal(t, check.name, node.Command[0])
		}
	}
}

func TestFindSequence_chords(t *testing.T) {
	patterns := createPatterns(
		"lsrs",
		"chord:lsrs:lsrs",
	)

	expected := []Expected{
		{"lsrs", true, false},
		{"chord:lsrs", false, true},
		{"chord:lsrs:lsrs", true, false},
	}

	checkResults(t, patterns, expected)
}

func TestFindSequence_simple(t *testing.T) {
	patterns := createPatterns(
		"ls",
		"rs",
		"lsls",
		"lsrs",
	)

	expected := []Expected{
		{"ls", true, true},
		{"rs", true, false},
		{"rsrs", false, false},
		{"lsls", true, false},
		{"lsrs", true, false},
		{"lslsls", false, false},
	}

	checkResults(t, patterns, expected)
}

func TestFindSequence_empty(t *testing.T) {
	patterns := createPatterns()

	expected := []Expected{
		{"ls", false, false},
	}

	checkResults(t, patterns, expected)
}
