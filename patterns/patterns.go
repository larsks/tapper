package patterns

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"tapper/keys"

	evdev "github.com/holoplot/go-evdev"
)

type (
	PatternNode struct {
		Chord       *keys.Chord
		Command     []string
		KeySequence keys.Sequence
		Next        []*PatternNode
		Terminal    bool
	}

	Patterns struct {
		PatternNode
	}
)

func ChordFromString(s string) (*keys.Chord, error) {
	keys := keys.NewChord()

	for _, tok := range strings.Fields(s) {
		tok = fmt.Sprintf("KEY_%s", tok)
		code, ok := evdev.KEYFromString[tok]
		if !ok {
			return keys, fmt.Errorf("%s: unknown key name", tok)
		}

		keys.Add(code)
	}

	return keys, nil
}

func NewPatterns() *Patterns {
	return new(Patterns)
}

func (patterns *Patterns) AddSequence(seq *keys.Sequence, Command []string) {
	node := &patterns.PatternNode
outer:
	for _, chord := range seq.Chords() {
		for _, next := range node.Next {
			if chord.Equal(next.Chord) {
				node = next
				continue outer
			}
		}

		p := PatternNode{
			Chord: chord,
		}

		node.Next = append(node.Next, &p)
		node = &p
	}

	node.Terminal = true
	node.Command = Command
}

func (patterns *Patterns) FindSequence(seq *keys.Sequence) (*PatternNode, bool, bool) {
	node := &patterns.PatternNode
	chords := seq.Chords()

	for len(chords) > 0 {
		for _, next := range node.Next {
			if next.Chord.Equal(chords[0]) {
				node = next
				if len(chords) == 1 && node.Terminal {
					return node, true, len(node.Next) > 0
				}
				break
			}
		}

		chords = chords[1:]
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
	if len(node.KeySequence.Chords()) == 0 {
		return fmt.Errorf("no key sequence")
	}
	return nil
}
