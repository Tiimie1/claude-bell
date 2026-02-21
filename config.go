package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Stop         string  `json:"stop,omitempty"`
	Notification string  `json:"notification,omitempty"`
	Limit        string  `json:"limit,omitempty"`
	Volume       float64 `json:"volume"`
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "claude-bell")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

func soundsDir() string {
	return filepath.Join(configDir(), "sounds")
}

func claudeSettingsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude", "settings.json")
}

func loadConfig() (Config, error) {
	cfg := Config{Volume: 1.0}
	data, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	var disk struct {
		Stop         string   `json:"stop,omitempty"`
		Notification string   `json:"notification,omitempty"`
		Limit        string   `json:"limit,omitempty"`
		Volume       *float64 `json:"volume"`
	}
	if err := json.Unmarshal(data, &disk); err != nil {
		return cfg, err
	}

	cfg.Stop = disk.Stop
	cfg.Notification = disk.Notification
	cfg.Limit = disk.Limit
	if disk.Volume != nil {
		cfg.Volume = clampVolume(*disk.Volume)
	}

	return cfg, nil
}

func getConfigField(cfg Config, event string) string {
	switch event {
	case "stop":
		return cfg.Stop
	case "notification":
		return cfg.Notification
	case "limit":
		return cfg.Limit
	}
	return ""
}

func setConfigField(cfg Config, event, value string) Config {
	switch event {
	case "stop":
		cfg.Stop = value
	case "notification":
		cfg.Notification = value
	case "limit":
		cfg.Limit = value
	}
	return cfg
}

func saveConfig(cfg Config) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	cfg.Volume = clampVolume(cfg.Volume)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}

func clampVolume(v float64) float64 {
	switch {
	case v < 0:
		return 0
	case v > 1:
		return 1
	default:
		return v
	}
}
