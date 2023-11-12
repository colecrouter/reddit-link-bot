package config

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	DiscordToken string `yaml:"discord_token"`
}

func (c *Config) Load() error {
	// Read and write back the config file
	// This is to ensure any new fields are added to the config file

	f, err := os.OpenFile("./config.yml", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	// Create a default config
	def := &Config{
		DiscordToken: "token",
	}

	// Read the config file
	b, err := os.ReadFile("./config.yml")
	if err != nil {
		return err
	}

	// Unmarshal the config file
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return err
	}

	// Marshal the default config
	b, err = yaml.Marshal(def)
	if err != nil {
		return err
	}

	// Write the default config back to the file
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
