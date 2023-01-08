package environ

import (
	"fmt"
	"os"
)

type (
	Environ struct {
		Prefix string
	}
)

func NewEnviron(prefix string) *Environ {
	return &Environ{
		Prefix: prefix,
	}
}

func (env *Environ) Get(key, fallback string) string {
	if env.Prefix != "" {
		key = fmt.Sprintf("%s_%s", env.Prefix, key)
	}

	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
