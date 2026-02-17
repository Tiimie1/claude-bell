package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ensureSound(event, presetName string) (string, error) {
	dir := soundsDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s_%s.wav", event, sanitize(presetName))
	path := filepath.Join(dir, filename)

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// Check built-in presets
	if presets, ok := EventPresets[event]; ok {
		for _, p := range presets {
			if p.Name == presetName {
				if err := generateWAV(path, p.Tones); err != nil {
					return "", err
				}
				return path, nil
			}
		}
	}

	// Fall back to custom sounds
	cs, err := findCustomSound(presetName)
	if err != nil {
		return "", err
	}
	if cs != nil {
		if err := generateWAV(path, cs.Tones); err != nil {
			return "", err
		}
		return path, nil
	}

	return "", fmt.Errorf("unknown preset %q for event %q", presetName, event)
}

// playSound plays a WAV file using afplay. It blocks until playback finishes.
func playSound(path string) error {
	cmd := exec.Command("afplay", path)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// sanitize replaces spaces with underscores and lowercases for filenames.
func sanitize(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' {
			out = append(out, '_')
		} else if c >= 'A' && c <= 'Z' {
			out = append(out, c+32)
		} else {
			out = append(out, c)
		}
	}
	return string(out)
}
