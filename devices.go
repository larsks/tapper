package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"tapper/config"

	evdev "github.com/holoplot/go-evdev"
)

func listInputDevices() []string {
	var devices []string

	basePath := config.Options.GetString("device_base_path")
	files, err := os.ReadDir(basePath)
	if err != nil {
		log.Printf("failed to read %s: %v", basePath, err)
		return nil
	}

	for _, fileName := range files {
		if fileName.Type()&fs.ModeCharDevice == 0 {
			continue
		}

		full := fmt.Sprintf("%s/%s", basePath, fileName.Name())
		devices = append(devices, full)
	}

	return devices
}

func findDeviceByName(want string) (*evdev.InputDevice, error) {
	log.Printf("looking for device matching %s", want)
	for _, fileName := range listInputDevices() {
		if dev, err := evdev.Open(fileName); err == nil {
			if have, err := dev.Name(); err == nil {
				if match, err := filepath.Match(want, have); err == nil && match {
					log.Printf("found device %s", fileName)
					return dev, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("failed to find device matching \"%s\"", want)
}
