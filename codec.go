package main

import (
	"encoding/base64"
	"fmt"
	"math"
)

// midiToFreq converts a MIDI note number to frequency in Hz.
// MIDI 0 = silence, MIDI 1-127 = notes.
func midiToFreq(midi byte) float64 {
	if midi == 0 {
		return 0
	}
	return 440.0 * math.Pow(2.0, (float64(midi)-69.0)/12.0)
}

// decodeTones decodes a base64url string back into a slice of Tones.
func decodeTones(code string) ([]Tone, error) {
	data, err := base64.RawURLEncoding.DecodeString(code)
	if err != nil {
		return nil, fmt.Errorf("invalid code: %w", err)
	}
	if len(data) == 0 || len(data)%2 != 0 {
		return nil, fmt.Errorf("invalid code: must contain an even number of bytes")
	}
	tones := make([]Tone, len(data)/2)
	for i := range tones {
		midi := data[i*2]
		ticks := data[i*2+1]
		tones[i] = Tone{
			Freq:     midiToFreq(midi),
			Duration: float64(ticks) / 100.0,
		}
	}
	return tones, nil
}
