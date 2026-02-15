package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "setup":
		cmdSetup()
	case "test":
		cmdTest()
	case "install":
		cmdInstall()
	case "uninstall":
		cmdUninstall()
	case "play":
		cmdPlay()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`claude-bell - Notification sounds for Claude Code

Usage:
  claude-bell <command>

Commands:
  setup       Interactive setup: pick a sound for each event
  test        Play all configured sounds
  install     Add hooks to ~/.claude/settings.json
  uninstall   Remove hooks from ~/.claude/settings.json
  play <event>  Play sound for an event (used by hooks)
`)
}

func cmdPlay() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: claude-bell play <event>")
		os.Exit(1)
	}
	event := os.Args[2]

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var presetName string
	switch event {
	case "stop":
		presetName = cfg.Stop
	case "notification":
		presetName = cfg.Notification
	case "limit":
		presetName = cfg.Limit
	default:
		fmt.Fprintf(os.Stderr, "unknown event: %s\n", event)
		os.Exit(1)
	}

	if presetName == "" {
		return // no sound configured, exit silently
	}

	path, err := ensureSound(event, presetName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := playSound(path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cmdTest() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	events := []struct {
		name   string
		preset string
	}{
		{"stop", cfg.Stop},
		{"notification", cfg.Notification},
		{"limit", cfg.Limit},
	}

	any := false
	for _, e := range events {
		if e.preset == "" {
			continue
		}
		any = true
		fmt.Printf("Playing %s: %s\n", e.name, e.preset)
		path, err := ensureSound(e.name, e.preset)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  error: %v\n", err)
			continue
		}
		if err := playSound(path); err != nil {
			fmt.Fprintf(os.Stderr, "  playback error: %v\n", err)
		}
	}

	if !any {
		fmt.Println("No sounds configured. Run 'claude-bell setup' first.")
	}
}

func executablePath() (string, bool) {
	exe, err := os.Executable()
	if err != nil {
		return "", false
	}
	if strings.Contains(exe, "go-build") || strings.Contains(exe, "/tmp/") {
		return exe, true // path is valid but is a temp binary
	}
	return exe, false
}
