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
	updated := cfg

	fmt.Println("claude-bell setup")
	fmt.Println("=================")
	fmt.Println("Interactive setup for Claude Code notification sounds.")
	fmt.Printf("Playback volume: %s (change with 'claude-bell volume')\n", formatVolume(cfg.Volume))
	fmt.Println("Input: number=select, p<number>=preview, s=skip, Enter=keep current.")
	fmt.Println()

	for idx, event := range EventNames {
		current := getConfigField(updated, event)
		options := buildEventOptions(event, customSounds)

		fmt.Printf("[%d/%d] %s\n", idx+1, len(EventNames), event)
		fmt.Printf("  %s\n", EventDescriptions[event])
		if current == "" {
			fmt.Println("  Current: (none)")
		} else {
			fmt.Printf("  Current: %s\n", current)
		}
		fmt.Println()

		for i, opt := range options {
			label := opt.name
			if opt.custom {
				label += " (custom)"
			}
			if strings.EqualFold(opt.name, current) {
				label += " [current]"
			}
			fmt.Printf("  %d) %s\n", i+1, label)
		}
		fmt.Println("  s) Skip (no sound)")
		fmt.Println()

		choice, ok := promptEventChoice(reader, event, current, options, updated.Volume)
		if !ok {
			fmt.Fprintln(os.Stderr, "setup canceled")
			os.Exit(1)
		}
		updated = setConfigField(updated, event, choice)
		fmt.Println()
	}

	fmt.Println("Summary")
	fmt.Println("-------")
	for _, event := range EventNames {
		selected := getConfigField(updated, event)
		if selected == "" {
			selected = "(none)"
		}
		fmt.Printf("  %-14s %s\n", event+":", selected)
	}
	fmt.Printf("  %-14s %s\n", "volume:", formatVolume(updated.Volume))
	fmt.Println()

	if !promptYesNo(reader, "Save these settings? [Y/n]: ", true) {
		fmt.Println("No changes saved.")
		return
	}

	if err := saveConfig(updated); err != nil {
		fmt.Fprintf(os.Stderr, "error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config saved.")
	fmt.Println("Next step: run 'claude-bell install' to add hooks to Claude Code.")
}

func buildEventOptions(event string, customSounds []CustomSound) []menuOption {
	presets := EventPresets[event]
	options := make([]menuOption, 0, len(presets)+len(customSounds))

	for _, p := range presets {
		options = append(options, menuOption{name: p.Name})
	}
	for _, cs := range customSounds {
		options = append(options, menuOption{name: cs.Name, custom: true})
	}

	return options
}

func promptEventChoice(reader *bufio.Reader, event, current string, options []menuOption, volume float64) (string, bool) {
	for {
		fmt.Print("  Choice: ")
		input, err := reader.ReadString('\n')
		if err != nil && len(input) == 0 {
			return "", false
		}
		input = strings.TrimSpace(input)
		lower := strings.ToLower(input)

		if input == "" {
			if current == "" {
				fmt.Println("  Keeping: (none)")
			} else {
				fmt.Printf("  Keeping: %s\n", current)
			}
			return current, true
		}

		if lower == "s" || lower == "skip" || input == "0" {
			fmt.Println("  Selected: (none)")
			return "", true
		}

		if strings.HasPrefix(lower, "p") {
			num, err := strconv.Atoi(strings.TrimSpace(lower[1:]))
			if err == nil {
				idx := num - 1
				if idx >= 0 && idx < len(options) {
					fmt.Printf("  Previewing: %s\n", options[idx].name)
					path, err := ensureSound(event, options[idx].name)
					if err != nil {
						fmt.Fprintf(os.Stderr, "  error: %v\n", err)
						continue
					}
					if err := playSound(path, volume); err != nil {
						fmt.Fprintf(os.Stderr, "  playback error: %v\n", err)
					}
					continue
				}
			}
		}

		num, err := strconv.Atoi(input)
		if err == nil {
			idx := num - 1
			if idx >= 0 && idx < len(options) {
				fmt.Printf("  Selected: %s\n", options[idx].name)
				return options[idx].name, true
			}
		}

		fmt.Printf("  Invalid input. Enter 1-%d, p1-p%d, s, or Enter.\n", len(options), len(options))
	}
}

func promptYesNo(reader *bufio.Reader, prompt string, defaultYes bool) bool {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil && len(input) == 0 {
			return defaultYes
		}
		s := strings.ToLower(strings.TrimSpace(input))

		if s == "" {
			return defaultYes
		}
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}

		fmt.Println("Please enter y or n.")
	}
}
