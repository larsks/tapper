package environ

import (
	"fmt"
	"os"
	"strconv"
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

func (env *Environ) GetInt(key string, fallback int64) int64 {
	val := env.Get(key, "")

	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fallback
	}

	return intVal
}
