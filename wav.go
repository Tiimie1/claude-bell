package main

import (
	"encoding/binary"
	"math"
	"os"
)

const (
	sampleRate = 44100
	bitDepth   = 16
	numChans   = 1
)

const fadeDurationMs = 5

func calcFadeSamples() int {
	return fadeDurationMs * sampleRate / 1000
}

// generateWAV creates a WAV file from a sequence of tones.
func generateWAV(path string, tones []Tone) error {
	samples := renderTones(tones)
	return writeWAV(path, samples)
}

// renderTones generates PCM samples for a sequence of tones.
func renderTones(tones []Tone) []int16 {
	var samples []int16
	for _, t := range tones {
		numSamples := int(t.Duration * sampleRate)
		fade := calcFadeSamples()
		for i := 0; i < numSamples; i++ {
			var sample float64
			if t.Freq > 0 {
				sample = math.Sin(2 * math.Pi * t.Freq * float64(i) / sampleRate)
				// Apply fade-in envelope
				if i < fade {
					sample *= float64(i) / float64(fade)
				}
				// Apply fade-out envelope
				if i >= numSamples-fade {
					sample *= float64(numSamples-1-i) / float64(fade)
				}
				sample *= 0.5 // master volume
			}
			samples = append(samples, int16(sample*math.MaxInt16))
		}
	}
	return samples
}

// writeWAV writes PCM samples as a 16-bit mono WAV file.
func writeWAV(path string, samples []int16) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dataSize := uint32(len(samples) * 2) // 2 bytes per sample
	fileSize := 36 + dataSize             // total - 8 for RIFF header

	byteRate := uint32(sampleRate * numChans * bitDepth / 8)
	blockAlign := uint16(numChans * bitDepth / 8)

	// RIFF header
	f.Write([]byte("RIFF"))
	binary.Write(f, binary.LittleEndian, fileSize)
	f.Write([]byte("WAVE"))

	// fmt subchunk
	f.Write([]byte("fmt "))
	binary.Write(f, binary.LittleEndian, uint32(16)) // subchunk size
	binary.Write(f, binary.LittleEndian, uint16(1))  // PCM format
	binary.Write(f, binary.LittleEndian, uint16(numChans))
	binary.Write(f, binary.LittleEndian, uint32(sampleRate))
	binary.Write(f, binary.LittleEndian, byteRate)
	binary.Write(f, binary.LittleEndian, blockAlign)
	binary.Write(f, binary.LittleEndian, uint16(bitDepth))

	// data subchunk
	f.Write([]byte("data"))
	binary.Write(f, binary.LittleEndian, dataSize)
	binary.Write(f, binary.LittleEndian, samples)

	return nil
}
