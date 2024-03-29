package config

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Discord struct {
		Token    string   `yaml:"token"`
		Channels []string `yaml:"channels"`
		Roles    []string `yaml:"roles"`
		NoNSFW   bool     `yaml:"no_nsfw"`
	}

	// Optional
	Reddit *struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
	} `yaml:"reddit"`
}

func (c *Config) Load() error {
	// Read and write back the config file
	// This is to ensure any new fields are added to the config file

	f, err := os.OpenFile("./config.yml", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

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

	// Marshal the (potentially updated) config file
	b, err = yaml.Marshal(c)
	if err != nil {
		return err
	}

	// Write the new config back to the file
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
