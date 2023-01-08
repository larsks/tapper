package patterns

import (
	"fmt"
	"strings"

	evdev "github.com/holoplot/go-evdev"
	"golang.org/x/exp/maps"
)

type (
	Chord    map[evdev.EvCode]bool
	Sequence []Chord

	PatternNode struct {
		Keys     Chord
		Command  []string
		Next     []*PatternNode
		Terminal bool
	}

	Patterns struct {
		PatternNode
	}
)

func ChordFromString(s string) (Chord, error) {
	keys := make(Chord)

	for _, tok := range strings.Fields(s) {
		tok = fmt.Sprintf("KEY_%s", tok)
		code, ok := evdev.KEYFromString[tok]
		if !ok {
			return keys, fmt.Errorf("%s: unknown key name", tok)
		}

		keys[code] = true
	}

	return keys, nil
}

func NewPatterns() *Patterns {
	return new(Patterns)
}

func (patterns *Patterns) AddSequence(seq Sequence, Command []string) {
	node := &patterns.PatternNode
outer:
	for _, keys := range seq {
		for _, next := range node.Next {
			if maps.Equal(next.Keys, keys) {
				node = next
				continue outer
			}
		}

		p := PatternNode{
			Keys: keys,
		}

		node.Next = append(node.Next, &p)
		node = &p
	}

	node.Terminal = true
	node.Command = Command
}

func (patterns *Patterns) FindSequence(seq Sequence) (*PatternNode, bool) {
	node := &patterns.PatternNode

	for len(seq) > 0 {
		for _, next := range node.Next {
			if maps.Equal(next.Keys, seq[0]) {
				node = next
				if len(seq) == 1 && node.Terminal {
					return node, len(node.Next) > 0
				}
				break
			}
		}

		seq = seq[1:]
	}

	return nil, len(node.Next) > 0
}
