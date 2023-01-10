package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type (
	Device struct {
		Name string
		Path string
	}

	Action struct {
		Pattern []string
		Command []string
		Keys    []string
	}

	ConfigFile struct {
		Device         Device
		Interval       uint64
		DeviceBasePath string `mapstructure:"device_base_path"`
		Actions        []Action
	}
)

func (action *Action) String() string {
	return strings.Join(action.Pattern, " ")
}

var Options = viper.NewWithOptions(
	viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
)
var Config ConfigFile

func init() {
	Options.SetConfigType("yaml")
	Options.SetEnvPrefix("TAPPER")
	Options.AutomaticEnv()

	Options.SetDefault("interval", "200")
	Options.SetDefault("device_base_path", "/dev/input")
}

func LoadConfig(reader io.Reader) error {
	if err := Options.ReadConfig(reader); err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	if err := Options.Unmarshal(&Config); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return nil
}

func LoadConfigFromFile(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer fd.Close()

	return LoadConfig(fd)
}
