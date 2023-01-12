package patterns

import (
	"fmt"
	"tapper/keys"
	"testing"

	evdev "github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/assert"
)

type (
	SequenceMap map[string]*keys.Sequence
	Expected    struct {
		name  string
		found bool
		more  bool
	}
)

var sequenceMap = SequenceMap{
	"ls": keys.NewSequence(
		keys.NewChord(evdev.KEY_LEFTSHIFT),
	),
	"rs": keys.NewSequence(
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
	),
	"lsls": keys.NewSequence(
		keys.NewChord(evdev.KEY_LEFTSHIFT),
		keys.NewChord(evdev.KEY_LEFTSHIFT),
	),
	"lslsls": keys.NewSequence(
		keys.NewChord(evdev.KEY_LEFTSHIFT),
		keys.NewChord(evdev.KEY_LEFTSHIFT),
		keys.NewChord(evdev.KEY_LEFTSHIFT),
	),
	"lsrsls": keys.NewSequence(
		keys.NewChord(evdev.KEY_LEFTSHIFT),
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
		keys.NewChord(evdev.KEY_LEFTSHIFT),
	),
	"rsrsrs": keys.NewSequence(
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
	),
	"rsls": keys.NewSequence(
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
		keys.NewChord(evdev.KEY_LEFTSHIFT),
	),
	"lsrs": keys.NewSequence(
		keys.NewChord(evdev.KEY_LEFTSHIFT),
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
	),
	"lsrsrs": keys.NewSequence(
		keys.NewChord(evdev.KEY_LEFTSHIFT),
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
	),
	"rsrs": keys.NewSequence(
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
		keys.NewChord(evdev.KEY_RIGHTSHIFT),
	),
	"chord:lsrs": keys.NewSequence(
		keys.NewChord(evdev.KEY_LEFTSHIFT, evdev.KEY_RIGHTSHIFT),
	),
	"chord:lsrs:lsrs": keys.NewSequence(
		keys.NewChord(evdev.KEY_LEFTSHIFT, evdev.KEY_RIGHTSHIFT),
		keys.NewChord(evdev.KEY_LEFTSHIFT, evdev.KEY_RIGHTSHIFT),
	),
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
