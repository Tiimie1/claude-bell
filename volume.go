package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func cmdVolume() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		fmt.Printf("Current volume: %s\n", formatVolume(cfg.Volume))
		fmt.Println("Set a new level with: claude-bell volume <value>")
		fmt.Println("Accepted formats: 0-1, 0-100, or percent (e.g. 0.65, 65, 65%).")
		return
	}

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: claude-bell volume [value]")
		os.Exit(1)
	}

	vol, err := parseVolumeArg(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintln(os.Stderr, "usage: claude-bell volume [value]")
		os.Exit(1)
	}

	cfg.Volume = vol
	if err := saveConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Volume set to %s\n", formatVolume(vol))
}

func parseVolumeArg(input string) (float64, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return 0, fmt.Errorf("volume cannot be empty")
	}

	if strings.HasSuffix(s, "%") {
		pct, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(s, "%")), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid percent volume %q", input)
		}
		if pct < 0 || pct > 100 {
			return 0, fmt.Errorf("percent volume must be between 0 and 100")
		}
		return pct / 100, nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid volume %q", input)
	}
	if v < 0 {
		return 0, fmt.Errorf("volume cannot be negative")
	}
	if v <= 1 {
		return v, nil
	}
	if v <= 100 {
		return v / 100, nil
	}

	return 0, fmt.Errorf("volume must be between 0-1 or 0-100")
}

func formatVolume(v float64) string {
	return fmt.Sprintf("%.0f%% (%.2f)", clampVolume(v)*100, clampVolume(v))
}
