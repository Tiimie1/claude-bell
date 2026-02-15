# claude-bell

Notification sounds for [Claude Code](https://docs.anthropic.com/en/docs/claude-code). Get an audible chime when Claude finishes a task, needs your attention, or hits a context limit.

No dependencies. Pure Go. Generates sine-wave WAV files and plays them with macOS `afplay`.

## Install

### Homebrew (macOS)

```bash
brew install Tiimie1/tap/claude-bell
```

### Go install

```bash
go install github.com/Tiimie1/claude-bell@latest
```

### Install script

```bash
curl -sSL https://raw.githubusercontent.com/Tiimie1/claude-bell/main/install.sh | bash
```

### Manual

Download the binary for your architecture from [GitHub Releases](https://github.com/Tiimie1/claude-bell/releases), then move it to somewhere in your PATH:

```bash
chmod +x claude-bell
sudo mv claude-bell /usr/local/bin/
```

## Quick start

```bash
# 1. Pick a sound for each event
claude-bell setup

# 2. Install hooks into Claude Code
claude-bell install

# Done! You'll hear sounds when Claude Code triggers events.
```

## Usage

```
claude-bell setup       Pick a sound for each event
claude-bell test        Play all configured sounds
claude-bell install     Add hooks to ~/.claude/settings.json
claude-bell uninstall   Remove hooks from ~/.claude/settings.json
claude-bell play <event>  Play sound for an event (used by hooks)
```

## How it works

1. `claude-bell setup` lets you pick from preset sounds for three events
2. `claude-bell install` writes async [hooks](https://docs.anthropic.com/en/docs/claude-code/hooks) into `~/.claude/settings.json`
3. When Claude Code triggers an event, it runs `claude-bell play <event>`, which generates a WAV file (cached) and plays it via `afplay`

Config is stored in `~/.config/claude-bell/config.json`. Generated WAV files are cached in `~/.config/claude-bell/sounds/`.

## Available sounds

| Event | Preset | Description |
|-------|--------|-------------|
| **stop** | Major Chime | Bright C-E-G rising triad |
| | Octave Chime | Simple C4 to C5 jump |
| | Resolve | G major arpeggio (G-B-D-G) |
| **notification** | Doorbell | Classic E-C two-tone |
| | Attention | Double tap on A5 |
| | Question | Rising C to E interval |
| **limit** | Descending Warning | G-D-A falling pattern |
| | Low Buzz | Triple pulse on A3 |
| | Slide Down | E5 to E3 octave drop |

## Building from source

```bash
git clone https://github.com/Tiimie1/claude-bell.git
cd claude-bell
make build
```

## Requirements

- macOS (uses `afplay` for audio playback)

## License

[MIT](LICENSE)
