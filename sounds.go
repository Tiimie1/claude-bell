package main

// Tone represents a single sine-wave tone with frequency and duration.
type Tone struct {
	Freq     float64 `json:"freq"`     // Hz
	Duration float64 `json:"duration"` // seconds
}

// SoundPreset defines a named sequence of tones.
type SoundPreset struct {
	Name  string
	Tones []Tone
}

// EventPresets maps event names to their available sound presets.
var EventPresets = map[string][]SoundPreset{
	"stop": {
		{
			Name: "Major Chime",
			Tones: []Tone{
				{Freq: 523.25, Duration: 0.15}, // C5
				{Freq: 659.25, Duration: 0.15}, // E5
				{Freq: 783.99, Duration: 0.25}, // G5
			},
		},
		{
			Name: "Octave Chime",
			Tones: []Tone{
				{Freq: 261.63, Duration: 0.2}, // C4
				{Freq: 523.25, Duration: 0.3}, // C5
			},
		},
		{
			Name: "Resolve",
			Tones: []Tone{
				{Freq: 392.00, Duration: 0.12}, // G4
				{Freq: 493.88, Duration: 0.12}, // B4
				{Freq: 587.33, Duration: 0.12}, // D5
				{Freq: 783.99, Duration: 0.25}, // G5
			},
		},
	},
	"notification": {
		{
			Name: "Doorbell",
			Tones: []Tone{
				{Freq: 659.25, Duration: 0.2}, // E5
				{Freq: 523.25, Duration: 0.3}, // C5
			},
		},
		{
			Name: "Attention",
			Tones: []Tone{
				{Freq: 880.00, Duration: 0.1},  // A5
				{Freq: 0, Duration: 0.08},       // gap
				{Freq: 880.00, Duration: 0.12},  // A5
			},
		},
		{
			Name: "Question",
			Tones: []Tone{
				{Freq: 523.25, Duration: 0.15}, // C5
				{Freq: 659.25, Duration: 0.25}, // E5
			},
		},
	},
	"limit": {
		{
			Name: "Descending Warning",
			Tones: []Tone{
				{Freq: 392.00, Duration: 0.15}, // G4
				{Freq: 293.66, Duration: 0.15}, // D4
				{Freq: 220.00, Duration: 0.25}, // A3
			},
		},
		{
			Name: "Low Buzz",
			Tones: []Tone{
				{Freq: 220.00, Duration: 0.1},  // A3
				{Freq: 0, Duration: 0.06},       // gap
				{Freq: 220.00, Duration: 0.1},   // A3
				{Freq: 0, Duration: 0.06},       // gap
				{Freq: 220.00, Duration: 0.12},  // A3
			},
		},
		{
			Name: "Slide Down",
			Tones: []Tone{
				{Freq: 659.25, Duration: 0.12}, // E5
				{Freq: 493.88, Duration: 0.12}, // B4
				{Freq: 164.81, Duration: 0.25}, // E3
			},
		},
	},
}

// EventNames defines the display order for events.
var EventNames = []string{"stop", "notification", "limit"}

// EventDescriptions provides a human-readable description for each event.
var EventDescriptions = map[string]string{
	"stop":         "Task complete - Claude finishes responding",
	"notification": "Needs attention - permission prompts, questions",
	"limit":        "Context limit - auto-compaction triggered",
}
