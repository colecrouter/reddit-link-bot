package config

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	DiscordToken string `yaml:"discord_token"`
}

func (c *Config) Load() error {
	s, _ := os.Stat("./config.yml")
	if s == nil {
		f, err := os.Create("./config.yml")
		if err != nil {
			return err
		}
		defer f.Close()

		def := &Config{
			DiscordToken: "token",
		}

		b, err := yaml.Marshal(def)
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		if err != nil {
			return err
		}
	}

	b, err := os.ReadFile("./config.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, c)
	if err != nil {
		return err
	}

	return nil
}
