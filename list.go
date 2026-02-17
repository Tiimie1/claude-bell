package main

import (
	"fmt"
	"os"
)

func cmdList() {
	sounds, err := loadCustomSounds()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(sounds) == 0 {
		fmt.Println("No custom sounds. Use 'claude-bell create <name> <code>' to add one.")
		return
	}

	fmt.Println("Custom sounds:")
	for _, s := range sounds {
		fmt.Printf("  %s (%d tones) - code: %s\n", s.Name, len(s.Tones), s.Code)
	}
}

func cmdDelete() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: claude-bell delete <name>")
		os.Exit(1)
	}

	name := os.Args[2]
	if err := deleteCustomSound(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted custom sound %q\n", name)
}
