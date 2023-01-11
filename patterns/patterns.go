package patterns

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"tapper/keys"

	evdev "github.com/holoplot/go-evdev"
	"golang.org/x/exp/maps"
)

type (
	PatternNode struct {
		Chord       keys.Chord
		Command     []string
		KeySequence keys.Sequence
		Next        []*PatternNode
		Terminal    bool
	}

	Patterns struct {
		PatternNode
	}
)

func ChordFromString(s string) (keys.Chord, error) {
	keys := keys.NewChord()

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

func (patterns *Patterns) AddSequence(seq keys.Sequence, Command []string) {
	node := &patterns.PatternNode
outer:
	for _, keys := range seq {
		for _, next := range node.Next {
			if maps.Equal(next.Chord, keys) {
				node = next
				continue outer
			}
		}

		p := PatternNode{
			Chord: keys,
		}

		node.Next = append(node.Next, &p)
		node = &p
	}

	node.Terminal = true
	node.Command = Command
}

func (patterns *Patterns) FindSequence(seq keys.Sequence) (*PatternNode, bool, bool) {
	node := &patterns.PatternNode

	for len(seq) > 0 {
		for _, next := range node.Next {
			if maps.Equal(next.Chord, seq[0]) {
				node = next
				if len(seq) == 1 && node.Terminal {
					return node, true, len(node.Next) > 0
				}
				break
			}
		}

		seq = seq[1:]
	}

	return nil, false, len(node.Next) > 0
}

func (node *PatternNode) RunCommand() error {
	if len(node.Command) == 0 {
		return fmt.Errorf("no command")
	}

	log.Printf("running command %#v", node.Command)
	cmd := exec.Command(node.Command[0], node.Command[1:]...)
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Printf("failed to run command: %v", err)
		}
	}()
	return nil
}

func (node *PatternNode) SendKeySequence() error {
	if len(node.KeySequence) == 0 {
		return fmt.Errorf("no key sequence")
	}
	return nil
}
