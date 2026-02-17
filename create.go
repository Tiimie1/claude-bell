package main

import (
	"fmt"
	"os"
	"strings"
)

func cmdCreate() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: claude-bell create <name> <code>")
		os.Exit(1)
	}

	name := os.Args[2]
	code := os.Args[3]

	// Validate name doesn't conflict with built-in presets
	for _, presets := range EventPresets {
		for _, p := range presets {
			if strings.EqualFold(p.Name, name) {
				fmt.Fprintf(os.Stderr, "error: %q conflicts with a built-in preset name\n", name)
				os.Exit(1)
			}
		}
	}

	// Decode and validate
	tones, err := decodeTones(code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := addCustomSound(name, code, tones); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created custom sound %q (%d tones)\n", name, len(tones))
}
