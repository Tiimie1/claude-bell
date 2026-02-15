package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cmdSetup() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("claude-bell setup")
	fmt.Println("=================")
	fmt.Println()
	fmt.Println("For each event, pick a sound or enter 0 to skip.")
	fmt.Println("Enter p1, p2, p3 to preview a sound before selecting.")
	fmt.Println()

	for _, event := range EventNames {
		presets := EventPresets[event]
		desc := EventDescriptions[event]
		current := getConfigField(cfg, event)

		fmt.Printf("--- %s ---\n", event)
		fmt.Printf("  %s\n", desc)
		if current != "" {
			fmt.Printf("  Current: %s\n", current)
		}
		fmt.Println()

		for i, p := range presets {
			fmt.Printf("  %d) %s\n", i+1, p.Name)
		}
		fmt.Println("  0) Skip (no sound)")
		fmt.Println()

		for {
			fmt.Print("  Choice: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "0" {
				cfg = setConfigField(cfg, event, "")
				fmt.Println("  Skipped.")
				fmt.Println()
				break
			}

			if len(input) == 2 && input[0] == 'p' {
				idx := int(input[1] - '1')
				if idx >= 0 && idx < len(presets) {
					fmt.Printf("  Previewing: %s\n", presets[idx].Name)
					path, err := ensureSound(event, presets[idx].Name)
					if err != nil {
						fmt.Fprintf(os.Stderr, "  error: %v\n", err)
						continue
					}
					if err := playSound(path); err != nil {
						fmt.Fprintf(os.Stderr, "  playback error: %v\n", err)
					}
					continue
				}
			}

			if len(input) == 1 && input[0] >= '1' && input[0] <= '9' {
				idx := int(input[0] - '1')
				if idx >= 0 && idx < len(presets) {
					cfg = setConfigField(cfg, event, presets[idx].Name)
					fmt.Printf("  Selected: %s\n\n", presets[idx].Name)
					break
				}
			}

			fmt.Println("  Invalid input. Enter 1-3 to select, p1-p3 to preview, or 0 to skip.")
		}
	}

	if err := saveConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config saved! Run 'claude-bell install' to add hooks to Claude Code.")
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
