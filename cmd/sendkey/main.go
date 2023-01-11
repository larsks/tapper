package main

import (
	"log"
	"time"

	"tapper/keys"
	"tapper/virtualkbd"

	flag "github.com/spf13/pflag"
)

var optRepeat int

var DEFAULT_SEQUENCE []string = []string{
	"e", "c", "h", "o", "space",
	"leftshift:t", "leftshift:a", "leftshift:p", "leftshift:p", "leftshift:e", "leftshift:r", "space",
	"i", "s", "space", "c", "o", "o", "l", "enter",
}

func init() {
	flag.IntVarP(&optRepeat, "repeat", "r", 1, "Number of times to repeat key")
}

func main() {
	flag.Parse()

	seq := keys.Sequence{}
	allKeys := keys.NewChord()

	seqspec := flag.Args()
	if len(seqspec) == 0 {
		seqspec = DEFAULT_SEQUENCE
	}

	for _, arg := range seqspec {
		chord, err := keys.ChordFromString(arg)
		if err != nil {
			log.Fatalf("%s: invalid key: %v", arg, err)
		}

		seq = append(seq, chord)
		allKeys.Update(chord)
	}

	kbd := virtualkbd.NewKeyboard().WithKeys(allKeys.Keys())
	if err := kbd.Start(); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)
	log.Printf("sending sequence: %s", seq)
	for i := 0; i < optRepeat; i++ {
		if err := kbd.TypeSequence(seq); err != nil {
			panic(err)
		}
	}
}
