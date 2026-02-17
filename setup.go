package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type menuOption struct {
	name   string
	custom bool
}

func cmdSetup() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	customSounds, _ := loadCustomSounds()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("claude-bell setup")
	fmt.Println("=================")
	fmt.Println()
	fmt.Println("For each event, pick a sound or enter 0 to skip.")
	fmt.Println("Enter p1, p2, ... to preview a sound before selecting.")
	fmt.Println()

	for _, event := range EventNames {
		presets := EventPresets[event]
		desc := EventDescriptions[event]
		current := getConfigField(cfg, event)

		// Build combined options list
		var options []menuOption
		for _, p := range presets {
			options = append(options, menuOption{name: p.Name})
		}
		for _, cs := range customSounds {
			options = append(options, menuOption{name: cs.Name, custom: true})
		}

		fmt.Printf("--- %s ---\n", event)
		fmt.Printf("  %s\n", desc)
		if current != "" {
			fmt.Printf("  Current: %s\n", current)
		}
		fmt.Println()

		for i, opt := range options {
			if opt.custom {
				fmt.Printf("  %d) %s (custom)\n", i+1, opt.name)
			} else {
				fmt.Printf("  %d) %s\n", i+1, opt.name)
			}
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

			// Preview: p1, p2, p12, etc.
			if strings.HasPrefix(input, "p") {
				num, err := strconv.Atoi(input[1:])
				if err == nil {
					idx := num - 1
					if idx >= 0 && idx < len(options) {
						fmt.Printf("  Previewing: %s\n", options[idx].name)
						path, err := ensureSound(event, options[idx].name)
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
			}

			// Selection: 1, 2, 12, etc.
			num, err := strconv.Atoi(input)
			if err == nil {
				idx := num - 1
				if idx >= 0 && idx < len(options) {
					cfg = setConfigField(cfg, event, options[idx].name)
					fmt.Printf("  Selected: %s\n\n", options[idx].name)
					break
				}
			}

			fmt.Printf("  Invalid input. Enter 1-%d to select, p1-p%d to preview, or 0 to skip.\n", len(options), len(options))
		}
	}

	if err := saveConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config saved! Run 'claude-bell install' to add hooks to Claude Code.")
}
