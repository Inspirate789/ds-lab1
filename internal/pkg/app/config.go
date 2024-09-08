package app

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Logging struct {
		Level int `koanf:"level"`
	} `koanf:"logging"`
	Web WebConfig `koanf:"web"`
	DB  struct {
		DriverName       string `koanf:"driver_name"`
		ConnectionString string `koanf:"connection_string"`
	} `koanf:"db"`
}

func ReadLocalConfig(configPath string) (Config, error) {
	config := koanf.New(".")

	err := config.Load(file.Provider(configPath), yaml.Parser())
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	err = config.Unmarshal("", &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
