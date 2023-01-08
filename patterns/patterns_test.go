package patterns

import (
	"fmt"
	"testing"

	evdev "github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/assert"
)

func TestFindSequence(t *testing.T) {
	pattern_ls := Sequence{
		Chord{evdev.KEY_LEFTSHIFT: true},
	}

	pattern_rs := Sequence{
		Chord{evdev.KEY_RIGHTSHIFT: true},
	}

	pattern_lsls := Sequence{
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true},
	}

	pattern_lslsls := Sequence{
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true},
	}

	pattern_lsrsls := Sequence{
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true},
	}

	pattern_rsrsrs := Sequence{
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
	}

	pattern_rsls := Sequence{
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_LEFTSHIFT: true},
	}

	pattern_lsrs := Sequence{
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
	}

	pattern_lsrsrs := Sequence{
		Chord{evdev.KEY_LEFTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
	}

	pattern_rsrs := Sequence{
		Chord{evdev.KEY_RIGHTSHIFT: true},
		Chord{evdev.KEY_RIGHTSHIFT: true},
	}

	p := NewPatterns()
	p.AddSequence(pattern_lsls, []string{})
	p.AddSequence(pattern_lsrsls, []string{})
	p.AddSequence(pattern_lsrsrs, []string{})
	p.AddSequence(pattern_rsrsrs, []string{})
	p.AddSequence(pattern_rsrs, []string{})

	seqlist := []struct {
		seq   Sequence
		found bool
		more  bool
	}{
		{pattern_ls, false, true},
		{pattern_rs, false, true},
		{pattern_lsls, true, false},
		{pattern_lsrs, false, true},
		{pattern_lsrsls, true, false},
		{pattern_lsrsrs, true, false},
		{pattern_lslsls, false, false},
		{pattern_rsls, false, true},
		{pattern_rsrs, true, true},
	}

	for _, check := range seqlist {
		node, found, more := p.FindSequence(check.seq)
		fmt.Printf("---\nseq %v\nhave found (%+v) %t more %t\nwant found %t more %t\n",
			check.seq, node, found, more, check.found, check.more)

		assert.Equal(t, check.found, found)
		if check.found {
			assert.NotNil(t, node)
		} else {
			assert.Nil(t, node)
		}
		assert.Equal(t, check.more, more)
	}
}
