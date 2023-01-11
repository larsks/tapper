package virtualkbd

import (
	"time"

	"github.com/holoplot/go-evdev"
)

type (
	Keyboard struct {
		Dev         *evdev.InputDevice
		KeyDownTime time.Duration
		KeyInterval time.Duration

		name         string
		id           evdev.InputID
		capabilities map[evdev.EvType][]evdev.EvCode
	}
)
