package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configText string = `
device:
  name: "MY KEYBOARD"
interval: 999
actions:
  - pattern:
    - LEFTSHIFT
    - LEFTSHIFT
    command:
      - "false"
`

func TestParseConfig(t *testing.T) {
	reader := strings.NewReader(configText)
	err := LoadConfig(reader)
	assert.Nil(t, err)
	assert.Equal(t, "999", Options.GetString("interval"))
	assert.Equal(t, 1, len(Config.Actions))
}

func TestEnvironmentOverride(t *testing.T) {
	reader := strings.NewReader(configText)
	err := LoadConfig(reader)
	assert.Nil(t, err)

	os.Setenv("TAPPER_INTERVAL", "222")
	assert.Equal(t, "222", Options.GetString("interval"))

	os.Setenv("TAPPER_DEVICE_NAME", "ALICE")
	assert.Equal(t, "ALICE", Options.GetString("device.name"))
}
