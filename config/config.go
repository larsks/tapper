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
	}

	ConfigFile struct {
		Device   Device
		Interval uint64
		Actions  []Action
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
	Options.SetDefault("interval", "200")
	Options.AutomaticEnv()
}

func LoadConfig(reader io.Reader) error {
	if err := Options.ReadConfig(reader); err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	Options.Unmarshal(&Config)

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
