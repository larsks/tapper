package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"tapper/environ"

	evdev "github.com/holoplot/go-evdev"
)

var env *environ.Environ = environ.NewEnviron("TAPPER")

var DEVICE_BASE_PATH string = env.Get("DEVICE_BASE_PATH", "/dev/input")

func listInputDevices() []string {
	var devices []string

	files, err := os.ReadDir(DEVICE_BASE_PATH)
	if err != nil {
		log.Printf("failed to read %s: %v", DEVICE_BASE_PATH, err)
		return nil
	}

	for _, fileName := range files {
		if fileName.Type()&fs.ModeCharDevice == 0 {
			continue
		}

		full := fmt.Sprintf("%s/%s", DEVICE_BASE_PATH, fileName.Name())
		devices = append(devices, full)
	}

	return devices
}

func findDeviceByName(want string) (*evdev.InputDevice, error) {
	for _, fileName := range listInputDevices() {
		if dev, err := evdev.Open(fileName); err == nil {
			if have, err := dev.Name(); err == nil {
				if match, err := filepath.Match(want, have); err == nil && match {
					return dev, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("failed to find device matching \"%s\"", want)
}
