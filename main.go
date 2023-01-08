package main

import (
	"fmt"
	"log"
	"time"
	"tippytap/version"

	evdev "github.com/holoplot/go-evdev"
	flag "github.com/spf13/pflag"
)

type (
	Chord    map[evdev.EvCode]bool
	Sequence []Chord

	App struct {
		*TippyTapConfig
		dev *evdev.InputDevice
	}
)

var optConfigPath string
var optDeviceBasePath string
var optListDevices bool
var optListKeys bool
var optVersion bool
var optDebug bool

func init() {
	flag.StringVarP(&optConfigPath, "config", "f", "tippytap.yaml", "Path to configuration file")
	flag.StringVarP(&optDeviceBasePath, "device-path", "D", "/dev/input", "Base device path")
	flag.BoolVarP(&optDebug, "debug", "", false, "Show debug output")
	flag.BoolVarP(&optListDevices, "list-devices", "L", false, "List available devices")
	flag.BoolVarP(&optListKeys, "list-keys", "K", false, "List available keycodes")
	flag.BoolVarP(&optVersion, "version", "v", false, "Show version")
}

func NewApp(config *TippyTapConfig) *App {
	return &App{
		TippyTapConfig: config,
	}
}

func (app *App) GetDevice() error {
	var err error

	if app.Device.Path != "" {
		dev, err := evdev.Open(app.Device.Path)
		if err != nil {
			return fmt.Errorf("failed to open device %s: %v", app.Device.Path, err)
		}

		app.dev = dev
		return nil
	}

	app.dev, err = findDeviceByName(app.Device.Name)
	return err
}

func (app *App) KeyEvents() chan *evdev.InputEvent {
	events := make(chan *evdev.InputEvent)

	go func() {
		for {
			evt, err := app.dev.ReadOne()
			if err != nil {
				panic(fmt.Errorf("failed to read event: %w", err))
			}

			if evt.Type == evdev.EV_KEY {
				events <- evt
			}
		}
	}()

	return events
}

func (app *App) OpenDevice() error {
	err := app.GetDevice()
	if err != nil {
		return fmt.Errorf("failed to get device path: %w", err)
	}
	log.Printf("using device %s\n", app.dev.Path())
	return nil
}

func (app *App) KeyLoop() error {
	events := app.KeyEvents()
	timer := time.NewTimer(0)

	lastEvent := time.Time{}
	keysDown := make(Chord)
	cur := make(Chord)
	seq := Sequence{}

	for {
		select {
		case <-timer.C:
			fmt.Printf("timer\n")
		case evt := <-events:
			now := time.Now()
			elapsed := now.Sub(lastEvent)
			lastEvent = now

			if elapsed.Milliseconds() > app.Options.Interval {
				fmt.Printf("reset! %d %d\n", elapsed.Milliseconds(), app.Options.Interval)
				seq = Sequence{}
				cur = make(Chord)
			}

			if evt.Value == 0 {
				delete(keysDown, evt.Code)
			} else {
				keysDown[evt.Code] = true
			}

			fmt.Printf("keysdown %v cur %+v seq %v\n",
				keysDown, cur, seq)

			// all keys have been released
			if len(keysDown) == 0 {
				if len(cur) > 0 {
					seq = append(seq, cur)
					timer.Reset(time.Duration(app.Options.Interval) * time.Millisecond)
					fmt.Printf("seq %v\n", seq)
				}
				cur = make(Chord)
			} else {
				for key := range keysDown {
					cur[key] = true
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	if optVersion {
		printVersion()
		return
	}

	if optListDevices {
		printDevices()
		return
	}

	if optListKeys {
		printKeys()
		return
	}

	log.Printf("read config from %s\n", optConfigPath)
	config, err := ReadConfiguration(optConfigPath)
	if err != nil {
		log.Fatalf("failed to read configuration: %v", err)
	}

	app := NewApp(config)

	if optDebug {
		fmt.Printf("config\n")
		fmt.Printf("device: %+v\n", config.Device)
		fmt.Printf("options: %+v\n", config.Options)
		fmt.Printf("actions: %+v\n", config.Actions)
	}

	if err := app.OpenDevice(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	if err := app.KeyLoop(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}

func printDevices() {
	for _, fileName := range listInputDevices() {
		d, err := evdev.Open(fileName)
		if err == nil {
			name, _ := d.Name()

			if err == nil {
				fmt.Printf("%s:\t%s\n", d.Path(), name)
			}
		}
	}
}

func printVersion() {
	fmt.Printf("Version: %s\n", version.Version)
	fmt.Printf("BuildDate: %s\n", version.BuildDate)
	fmt.Printf("BuildRef: %s\n", version.BuildRef)
}

func printKeys() {
	for _, key := range evdev.KEYNames {
		fmt.Printf("%s\n", key)
	}
}
