package virtualkbd

import (
	"fmt"
	"tapper/keys"
	"time"

	"github.com/holoplot/go-evdev"
)

var DEFAULT_KEY_DOWN_TIME time.Duration = 10 * time.Millisecond
var DEFAULT_KEY_INTERVAL time.Duration = 100 * time.Millisecond

func NewKeyboard() *Keyboard {
	k := Keyboard{}
	k.name = "Tapper Virtual Keyboard"
	k.id = evdev.InputID{
		BusType: evdev.BUS_USB,
		Vendor:  0x03f0,
		Product: 0x034a,
		Version: 1,
	}
	k.capabilities = make(map[evdev.EvType][]evdev.EvCode)
	k.KeyDownTime = DEFAULT_KEY_DOWN_TIME
	k.KeyInterval = DEFAULT_KEY_INTERVAL

	return &k
}

func (kbd *Keyboard) WithInputID(id evdev.InputID) *Keyboard {
	kbd.id = id
	return kbd
}

func (kbd *Keyboard) WithCapabilities(caps map[evdev.EvType][]evdev.EvCode) *Keyboard {
	kbd.capabilities = caps
	return kbd
}

func (kbd *Keyboard) WithName(name string) *Keyboard {
	kbd.name = name
	return kbd
}

func (kbd *Keyboard) WithKeys(keyEvents []evdev.EvCode) *Keyboard {
	chord := make(keys.Chord)

	for _, key := range kbd.capabilities[evdev.EV_KEY] {
		chord[key] = true
	}

	for _, key := range keyEvents {
		chord[key] = true
	}

	updated := []evdev.EvCode{}
	for key := range chord {
		updated = append(updated, key)
	}

	kbd.capabilities[evdev.EV_KEY] = updated

	return kbd
}

func (kbd *Keyboard) Start() error {
	dev, err := evdev.CreateDevice(
		kbd.name, kbd.id, kbd.capabilities,
	)
	if err != nil {
		return fmt.Errorf("failed to create input device: %w", err)
	}

	kbd.Dev = dev
	return nil
}

func (kbd *Keyboard) KeyEvent(chord keys.Chord, value int32) error {
	for key := range chord {
		evt := evdev.InputEvent{
			Type:  evdev.EV_KEY,
			Code:  key,
			Value: value,
		}
		if err := kbd.Dev.WriteOne(&evt); err != nil {
			return fmt.Errorf("failed to send key %d: %w", key, err)
		}
	}

	evt := evdev.InputEvent{
		Type:  evdev.EV_SYN,
		Code:  evdev.SYN_REPORT,
		Value: 0,
	}

	if err := kbd.Dev.WriteOne(&evt); err != nil {
		return fmt.Errorf("failed to send syn report: %w", err)
	}
	return nil
}

func (kbd *Keyboard) KeyDown(chord keys.Chord) error {
	return kbd.KeyEvent(chord, 1)
}

func (kbd *Keyboard) KeyUp(chord keys.Chord) error {
	return kbd.KeyEvent(chord, 0)
}

func (kbd *Keyboard) TypeChord(chord keys.Chord) error {
	if err := kbd.KeyDown(chord); err != nil {
		return err
	}

	time.Sleep(kbd.KeyDownTime)
	if err := kbd.KeyUp(chord); err != nil {
		return err
	}

	return nil
}

func (kbd *Keyboard) TypeSequence(seq keys.Sequence) error {
	for _, chord := range seq {
		if err := kbd.TypeChord(chord); err != nil {
			return err
		}
		time.Sleep(kbd.KeyInterval)
	}

	return nil
}
