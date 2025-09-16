// Package sound provides audio playback functionality for alarm sounds
package sound

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/go-audio/wav"
)

// Player manages audio playback for alarm sounds
type Player struct {
	context     *oto.Context
	soundFiles  map[string]string // name -> filepath mapping
	currentPlay chan struct{}     // channel to stop current playback
	mu          sync.Mutex
}

// NewPlayer creates a new audio player instance
func NewPlayer() (*Player, error) {
	// Initialize oto context with proper settings
	op := &oto.NewContextOptions{
		SampleRate:   44100, // Standard sample rate
		ChannelCount: 2,     // Stereo
		Format:       oto.FormatSignedInt16LE, // 16-bit signed little endian
	}
	
	ctx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio context: %w", err)
	}

	// Wait for the audio context to be ready
	<-readyChan

	player := &Player{
		context:    ctx,
		soundFiles: make(map[string]string),
	}

	return player, nil
}

// LoadSoundsFromDirectory loads all WAV files from the specified directory
func (p *Player) LoadSoundsFromDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".wav") {
			// Create a friendly name from filename
			name := strings.TrimSuffix(file.Name(), ".wav")
			name = strings.ReplaceAll(name, "mixkit-", "")
			name = strings.ReplaceAll(name, "-", " ")
			name = strings.Title(name)
			
			fullPath := filepath.Join(dirPath, file.Name())
			p.soundFiles[name] = fullPath
		}
	}

	return nil
}

// GetAvailableSounds returns a slice of available sound names
func (p *Player) GetAvailableSounds() []string {
	p.mu.Lock()
	defer p.mu.Unlock()

	sounds := make([]string, 0, len(p.soundFiles))
	for name := range p.soundFiles {
		sounds = append(sounds, name)
	}
	return sounds
}

// PlaySound plays the specified sound file in a loop for the given duration
func (p *Player) PlaySound(soundName string, duration time.Duration) error {
	p.mu.Lock()
	filePath, exists := p.soundFiles[soundName]
	if !exists {
		p.mu.Unlock()
		return fmt.Errorf("sound '%s' not found", soundName)
	}

	// Stop any currently playing sound
	if p.currentPlay != nil {
		close(p.currentPlay)
	}
	p.currentPlay = make(chan struct{})
	stopChan := p.currentPlay
	p.mu.Unlock()

	// Load and decode WAV file
	audioData, sampleRate, err := p.loadWAVFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to load audio file %s: %w", filePath, err)
	}

	// Play in a separate goroutine
	go p.playLoop(audioData, sampleRate, duration, stopChan)

	return nil
}

// StopSound stops any currently playing sound
func (p *Player) StopSound() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.currentPlay != nil {
		close(p.currentPlay)
		p.currentPlay = nil
	}
}

// loadWAVFile loads and decodes a WAV file
func (p *Player) loadWAVFile(filePath string) ([]byte, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		return nil, 0, fmt.Errorf("invalid WAV file: %s", filePath)
	}

	// Read audio data
	audioBuffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, 0, err
	}

	// Get format information
	format := decoder.Format()
	sampleRate := int(format.SampleRate)
	channels := int(format.NumChannels)
	
	// Convert to the format expected by oto (16-bit stereo at 44100 Hz)
	var audioData []byte
	
	if channels == 1 {
		// Mono to stereo - duplicate each sample
		audioData = make([]byte, len(audioBuffer.Data)*4) // 2 channels * 2 bytes per sample
		for i, sample := range audioBuffer.Data {
			// Left channel
			audioData[i*4] = byte(sample)
			audioData[i*4+1] = byte(sample >> 8)
			// Right channel (duplicate)
			audioData[i*4+2] = byte(sample)
			audioData[i*4+3] = byte(sample >> 8)
		}
	} else {
		// Assume stereo, convert to bytes
		audioData = make([]byte, len(audioBuffer.Data)*2)
		for i, sample := range audioBuffer.Data {
			audioData[i*2] = byte(sample)
			audioData[i*2+1] = byte(sample >> 8)
		}
	}

	return audioData, sampleRate, nil
}

// playLoop plays audio data in a loop for the specified duration
func (p *Player) playLoop(audioData []byte, sampleRate int, duration time.Duration, stopChan chan struct{}) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Calculate how long one loop of the audio takes
	// audioData contains stereo 16-bit samples (4 bytes per sample frame)
	sampleFrames := len(audioData) / 4 // 2 channels * 2 bytes per sample
	audioDuration := time.Duration(sampleFrames) * time.Second / time.Duration(sampleRate)
	
	// If the calculated duration is too short, set a minimum
	if audioDuration < 100*time.Millisecond {
		audioDuration = 100 * time.Millisecond
	}
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-stopChan:
			return
		default:
			// Create a new player for this playback using bytes.Reader
			reader := bytes.NewReader(audioData)
			player := p.context.NewPlayer(reader)
			player.Play()

			// Wait for the audio to finish or stop signal
			select {
			case <-time.After(audioDuration):
				// Audio finished, continue loop
				player.Close()
			case <-ctx.Done():
				player.Close()
				return
			case <-stopChan:
				player.Close()
				return
			}
		}
	}
}

// Close closes the audio player and releases resources
func (p *Player) Close() error {
	p.StopSound()
	// Note: oto v3 Context doesn't have a Close method
	// Resources are automatically cleaned up by the garbage collector
	return nil
}
