package main

import (
	"fmt"
	"log"
	"time"
	"tippytap/patterns"
	"tippytap/version"

	evdev "github.com/holoplot/go-evdev"
	flag "github.com/spf13/pflag"
)

type (
	App struct {
		*TippyTapConfig
		dev        *evdev.InputDevice
		patternMap *patterns.Patterns
		activeKeys patterns.Chord
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
				_, ok := app.activeKeys[evt.Code]
				if ok {
					events <- evt
				}
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

func (app *App) LoadPatterns() error {
	app.patternMap = patterns.NewPatterns()
	app.activeKeys = make(patterns.Chord)

	for _, actionconf := range app.Actions {
		seq := patterns.Sequence{}
		for _, keys := range actionconf.Pattern {
			chord, err := patterns.ChordFromString(keys)
			if err != nil {
				return fmt.Errorf("failed to read keys: %w", err)
			}
			seq = append(seq, chord)

			for key := range chord {
				app.activeKeys[key] = true
			}
		}
		app.patternMap.AddSequence(seq, actionconf.Command)
	}

	return nil
}

func (app *App) Init() error {
	if err := app.OpenDevice(); err != nil {
		return err
	}

	if err := app.LoadPatterns(); err != nil {
		return fmt.Errorf("failed to load patterns: %w", err)
	}

	return nil
}

func (app *App) KeyLoop() error {
	events := app.KeyEvents()
	timer := time.NewTimer(0)

	lastEvent := time.Time{}
	keysDown := make(patterns.Chord)
	cur := make(patterns.Chord)
	seq := patterns.Sequence{}

	var match *patterns.PatternNode

	reset := func() {
		fmt.Printf("reset!\n")
		seq = patterns.Sequence{}
		cur = make(patterns.Chord)
		match = nil
	}

	for {
		select {
		case <-timer.C:
			fmt.Printf("timer match %+v\n", match)
			if match != nil {
				fmt.Printf("timer execute: %+v\n", match)
				_ = match.RunCommand()
				reset()
			}
		case evt := <-events:
			now := time.Now()
			elapsed := now.Sub(lastEvent)
			lastEvent = now

			if elapsed.Milliseconds() > app.Options.Interval {
				reset()
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
					var found, more bool
					seq = append(seq, cur)
					match, found, more = app.patternMap.FindSequence(seq)
					if found {
						if more {
							timer.Reset(time.Duration(app.Options.Interval) * time.Millisecond)
						} else {
							fmt.Printf("execute: %+v\n", match)
							_ = match.RunCommand()
							reset()
						}
					}
				}
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

	if err := app.Init(); err != nil {
		log.Fatalf("ERROR: failed to initialize: %v", err)
	}

	if err := app.KeyLoop(); err != nil {
		log.Fatalf("ERROR: failed to process events: %v", err)
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
	for code, key := range evdev.KEYNames {
		fmt.Printf("[%03d] %s\n", code, key)
	}
}
