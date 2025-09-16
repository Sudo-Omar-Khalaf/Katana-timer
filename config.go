package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	NotificationThresholdHours float64 `json:"notification_threshold_hours"`
	AutoExportEnabled         bool    `json:"auto_export_enabled"`
	ExportPath               string  `json:"export_path"`
	TimerUpdateIntervalMs    int     `json:"timer_update_interval_ms"`
	MaxActivityLength        int     `json:"max_activity_length"`
	MaxTagLength            int     `json:"max_tag_length"`
	MaxTags                 int     `json:"max_tags"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		NotificationThresholdHours: 2.0,
		AutoExportEnabled:         false,
		ExportPath:               ".",
		TimerUpdateIntervalMs:    500,
		MaxActivityLength:        100,
		MaxTagLength:            20,
		MaxTags:                 5,
	}
}

// LoadConfig loads configuration from file or creates default
func LoadConfig() (*Config, error) {
	configPath := filepath.Join("data", "config.json")
	
	// Create data directory if it doesn't exist
	os.MkdirAll("data", 0755)
	
	// If config file doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		return config, config.Save()
	}
	
	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), err
	}
	
	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return DefaultConfig(), err
	}
	
	return config, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath := filepath.Join("data", "config.json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
