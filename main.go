package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tapper/config"
	"tapper/patterns"
	"tapper/version"
	"time"

	"github.com/adrg/xdg"
	evdev "github.com/holoplot/go-evdev"
	flag "github.com/spf13/pflag"
)

type (
	App struct {
		dev        *evdev.InputDevice
		patternMap *patterns.Patterns
		activeKeys patterns.Chord
	}
)

var optConfigPath string
var optListDevices bool
var optListKeys bool
var optVersion bool

func init() {
	configFilePath, ok := os.LookupEnv("TAPPER_CONFIG_FILE")
	if !ok {
		configFilePath = filepath.Join(xdg.ConfigHome, "tapper", "tapper.yaml")
	}

	flag.StringVarP(&optConfigPath, "config", "f", configFilePath, "Path to configuration file")
	flag.BoolVarP(&optListDevices, "list-devices", "L", false, "List available devices")
	flag.BoolVarP(&optListKeys, "list-keys", "K", false, "List available keycodes")
	flag.BoolVarP(&optVersion, "version", "v", false, "Show version")
}

func NewApp() *App {
	return &App{}
}

func (app *App) GetDevice() error {
	var err error

	path := config.Options.GetString("device.path")
	name := config.Options.GetString("device.name")

	if path != "" {
		app.dev, err = evdev.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open device %s: %v", path, err)
		}

		return nil
	}

	app.dev, err = findDeviceByName(name)
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
	return nil
}

func (app *App) LoadPatterns() error {
	app.patternMap = patterns.NewPatterns()
	app.activeKeys = make(patterns.Chord)

	for _, actionconf := range config.Config.Actions {
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

	lastEvent := time.Time{}
	keysDown := make(patterns.Chord)
	cur := make(patterns.Chord)
	seq := patterns.Sequence{}
	interval := time.Duration(config.Options.GetUint64("interval")) * time.Millisecond

	var match *patterns.PatternNode
	var timer *time.Timer

	reset := func() {
		seq = patterns.Sequence{}
		cur = make(patterns.Chord)
		match = nil
	}

	runCommand := func() {
		if match != nil {
			_ = match.RunCommand()
			reset()
		}
	}

	for evt := range events {
		now := time.Now()
		elapsed := now.Sub(lastEvent)
		lastEvent = now

		if timer != nil {
			timer.Stop()
		}

		if elapsed > interval {
			reset()
		}

		if evt.Value == 0 {
			delete(keysDown, evt.Code)
		} else {
			keysDown[evt.Code] = true
		}

		// all keys have been released
		if len(keysDown) == 0 {
			if len(cur) > 0 {
				var found, more bool
				seq = append(seq, cur)
				cur = make(patterns.Chord)
				match, found, more = app.patternMap.FindSequence(seq)
				if found {
					if more {
						timer = time.AfterFunc(interval, runCommand)
					} else {
						runCommand()
					}
				}
			}
		} else {
			for key := range keysDown {
				cur[key] = true
			}
		}
	}

	return nil
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

	if err := config.LoadConfigFromFile(optConfigPath); err != nil {
		log.Fatalf("ERROR: failed to read configuration: %v", err)
	}

	app := NewApp()

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
