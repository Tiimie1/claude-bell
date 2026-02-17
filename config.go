package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Stop         string `json:"stop,omitempty"`
	Notification string `json:"notification,omitempty"`
	Limit        string `json:"limit,omitempty"`
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
	var cfg Config
	data, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
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
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}
