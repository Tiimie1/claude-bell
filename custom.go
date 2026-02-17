package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CustomSound represents a user-created sound with its encoded representation.
type CustomSound struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Tones []Tone `json:"tones"`
}

func customSoundsPath() string {
	return filepath.Join(configDir(), "custom-sounds.json")
}

func loadCustomSounds() ([]CustomSound, error) {
	data, err := os.ReadFile(customSoundsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var sounds []CustomSound
	if err := json.Unmarshal(data, &sounds); err != nil {
		return nil, err
	}
	return sounds, nil
}

func saveCustomSounds(sounds []CustomSound) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(sounds, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(customSoundsPath(), data, 0644)
}

func findCustomSound(name string) (*CustomSound, error) {
	sounds, err := loadCustomSounds()
	if err != nil {
		return nil, err
	}
	for i := range sounds {
		if strings.EqualFold(sounds[i].Name, name) {
			return &sounds[i], nil
		}
	}
	return nil, nil
}

func addCustomSound(name, code string, tones []Tone) error {
	sounds, err := loadCustomSounds()
	if err != nil {
		return err
	}

	// Check for duplicate name
	for _, s := range sounds {
		if strings.EqualFold(s.Name, name) {
			return fmt.Errorf("custom sound %q already exists (use delete first to replace)", name)
		}
	}

	sounds = append(sounds, CustomSound{
		Name:  name,
		Code:  code,
		Tones: tones,
	})
	return saveCustomSounds(sounds)
}

func deleteCustomSound(name string) error {
	sounds, err := loadCustomSounds()
	if err != nil {
		return err
	}

	for i, s := range sounds {
		if strings.EqualFold(s.Name, name) {
			sounds = append(sounds[:i], sounds[i+1:]...)
			return saveCustomSounds(sounds)
		}
	}
	return fmt.Errorf("custom sound %q not found", name)
}
