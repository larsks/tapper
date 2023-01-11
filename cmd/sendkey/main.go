package main

import (
	"fmt"
	"strings"
	"time"

	"tapper/keys"
	"tapper/virtualkbd"

	"github.com/holoplot/go-evdev"
	flag "github.com/spf13/pflag"
)

var optRepeat int

func init() {
	flag.IntVarP(&optRepeat, "repeat", "r", 1, "Number of times to repeat key")
}

func main() {
	flag.Parse()

	chord := keys.NewChord()

	for _, arg := range flag.Args() {
		key, ok := evdev.KEYFromString[fmt.Sprintf("KEY_%s", strings.ToUpper(arg))]
		if !ok {
			fmt.Printf("invalid key name: %s\n", arg)
		}
		chord.Add(key)
	}

	kbd := virtualkbd.NewKeyboard().WithKeys(chord.Keys())
	if err := kbd.Start(); err != nil {
		panic(err)
	}

	for i := 0; i < optRepeat; i++ {
		if err := kbd.TypeChord(chord); err != nil {
			panic(err)
		}
		time.Sleep(virtualkbd.DEFAULT_KEY_INTERVAL)
	}
}
