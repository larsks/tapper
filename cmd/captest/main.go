package main

import (
	"fmt"
	"os"
	"syscall"
	"tapper/virtualkbd"

	"github.com/holoplot/go-evdev"
)

func main() {
	proto, err := evdev.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	kbd := virtualkbd.NewKeyboard().WithPrototype(proto)
	kbd.Start()
	fmt.Printf("kbd: %+v\n", kbd)
	syscall.Pause()
}
