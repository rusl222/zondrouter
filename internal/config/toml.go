package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

var version string = "0.2.2"

type Conf struct {
	Version  string
	Modbus_m []struct {
		Description string
		Master      Direction
		Slave       []Direction
	}
}

func Load(path string) error {
	var c Conf

	tomlData, err := os.ReadFile(path)
	if err != nil {
		return os.ErrNotExist
	}

	_, err = toml.Decode(string(tomlData), &c)
	if err != nil {
		return err
	}

	if c.Version != version {
		return errors.New("Version toml-file not supported")
	}

	for _, l := range c.Modbus_m {
		lines = append(lines, Line{
			Description: l.Description,
			Master:      l.Master,
			Slave:       l.Slave,
		})
	}

	return nil
}
