package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds all user configuration
type Config struct {
	Display DisplayConfig `toml:"display"`
}

// DisplayConfig holds display-related settings
type DisplayConfig struct {
	// WrapWidth is the maximum width for text wrapping (0 = use terminal width)
	WrapWidth int `toml:"wrap_width"`

	// ShowLineNumbers shows line numbers in the margin
	ShowLineNumbers bool `toml:"show_line_numbers"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Display: DisplayConfig{
			WrapWidth:       100, // Default to 100 characters
			ShowLineNumbers: false,
		},
	}
}

// ConfigPath returns the path to the configuration file
func ConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "gemini-client", "config.toml"), nil
}

// Load loads the configuration from the default location
// If the file doesn't exist, returns the default configuration
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist, return defaults
		return DefaultConfig(), nil
	}

	// Load the config file
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves the configuration to the default location
func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create the file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Encode the config
	encoder := toml.NewEncoder(f)
	return encoder.Encode(c)
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig() error {
	cfg := DefaultConfig()
	return cfg.Save()
}
