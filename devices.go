package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	evdev "github.com/holoplot/go-evdev"
)

func listInputDevices() []string {
	var devices []string

	files, err := os.ReadDir(optDeviceBasePath)
	if err != nil {
		log.Printf("failed to read %s: %v", optDeviceBasePath, err)
		return nil
	}

	for _, fileName := range files {
		if fileName.Type()&fs.ModeCharDevice == 0 {
			continue
		}

		full := fmt.Sprintf("%s/%s", optDeviceBasePath, fileName.Name())
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
