package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// hookDef defines how an event maps to a Claude Code hook.
type hookDef struct {
	event    string // our event name
	hookType string // Claude Code hook type
	matcher  string // optional matcher value
}

var hookDefs = []hookDef{
	{event: "stop", hookType: "Stop", matcher: ""},
	{event: "notification", hookType: "Notification", matcher: ""},
	{event: "limit", hookType: "PreCompact", matcher: "auto"},
}

func cmdInstall() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	if cfg.Stop == "" && cfg.Notification == "" && cfg.Limit == "" {
		fmt.Println("No sounds configured. Run 'claude-bell setup' first.")
		return
	}

	exePath, isTemp := executablePath()
	if isTemp {
		fmt.Println("WARNING: claude-bell is running from a temporary path (go run).")
		fmt.Println("Hooks will not work after this process exits.")
		fmt.Println("Build and install first: go build -o claude-bell . && sudo mv claude-bell /usr/local/bin/")
		fmt.Println()
	}

	settingsPath := claudeSettingsPath()

	settings, err := loadSettingsJSON(settingsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading settings: %v\n", err)
		os.Exit(1)
	}

	if err := backupSettings(settingsPath); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not create backup: %v\n", err)
	}

	hooks, _ := settings["hooks"].(map[string]any)
	if hooks == nil {
		hooks = make(map[string]any)
	}

	for _, hd := range hookDefs {
		preset := getConfigField(cfg, hd.event)
		if preset == "" {
			continue
		}

		hookKey := hd.hookType
		existing, _ := hooks[hookKey].([]any)

		var filtered []any
		for _, entry := range existing {
			m, ok := entry.(map[string]any)
			if !ok {
				filtered = append(filtered, entry)
				continue
			}
			if _, isBell := m["_claude_bell"]; isBell {
				continue
			}
			filtered = append(filtered, entry)
		}

		newEntry := map[string]any{
			"_claude_bell": true,
			"matcher":      hd.matcher,
			"hooks": []any{
				map[string]any{
					"type":    "command",
					"command": fmt.Sprintf("%s play %s", exePath, hd.event),
					"async":   true,
				},
			},
		}

		filtered = append(filtered, newEntry)
		hooks[hookKey] = filtered
	}

	settings["hooks"] = hooks

	if err := writeSettingsJSON(settingsPath, settings); err != nil {
		fmt.Fprintf(os.Stderr, "error writing settings: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Hooks installed into ~/.claude/settings.json")
	fmt.Println()
	for _, hd := range hookDefs {
		preset := getConfigField(cfg, hd.event)
		if preset == "" {
			continue
		}
		fmt.Printf("  %s (%s): %s\n", hd.event, hd.hookType, preset)
	}
}

func cmdUninstall() {
	settingsPath := claudeSettingsPath()

	settings, err := loadSettingsJSON(settingsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading settings: %v\n", err)
		os.Exit(1)
	}

	hooks, _ := settings["hooks"].(map[string]any)
	if hooks == nil {
		fmt.Println("No hooks found in settings.")
		return
	}

	removed := 0
	for key, val := range hooks {
		entries, ok := val.([]any)
		if !ok {
			continue
		}
		var filtered []any
		for _, entry := range entries {
			m, ok := entry.(map[string]any)
			if !ok {
				filtered = append(filtered, entry)
				continue
			}
			// Check for _claude_bell marker (new nested format)
			if _, isBell := m["_claude_bell"]; isBell {
				removed++
				continue
			}
			filtered = append(filtered, entry)
		}
		if len(filtered) == 0 {
			delete(hooks, key)
		} else {
			hooks[key] = filtered
		}
	}

	if len(hooks) == 0 {
		delete(settings, "hooks")
	} else {
		settings["hooks"] = hooks
	}

	if removed == 0 {
		fmt.Println("No claude-bell hooks found.")
		return
	}

	if err := writeSettingsJSON(settingsPath, settings); err != nil {
		fmt.Fprintf(os.Stderr, "error writing settings: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Removed %d claude-bell hook(s) from ~/.claude/settings.json\n", removed)
}

func loadSettingsJSON(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return nil, err
			}
			return make(map[string]any), nil
		}
		return nil, err
	}

	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	return settings, nil
}

func writeSettingsJSON(path string, settings map[string]any) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".settings-*.json")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}

	return os.Rename(tmpName, path)
}

func backupSettings(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	backupPath := path + ".claude-bell-backup"
	return os.WriteFile(backupPath, data, 0644)
}
