package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Address    string `yaml:"server_address"`
		WebAddress string `yaml:"web_address"`

		AccessToken string `yaml:"access_token"`
	} `yaml:"server"`
	ToggleableEntities []struct {
		EntityID string `yaml:"entity_id"`
		Domain   string `yaml:"domain"`
	} `yaml:"entities"`
}

// LoadConfig returns a new decoded Config struct
func LoadConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	// Default url
	if config.Server.WebAddress == "" {
		config.Server.WebAddress = "https://" + config.Server.Address
	}

	return config, nil
}
