package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
)

type AudioMixer struct {
	mu             sync.Mutex
	playing        map[int]*playback
	nextID         int
	channels       int
	mixAmp         float64
	framesPerWrite int
	sampleRate     int
	bytesPerSample int
	paused         bool
}

type playback struct {
	id            int
	file          *os.File
	done          chan struct{}
	buffer        []byte
	valid         int
	volume        float64
	totalBytes    int64
	bytesRead     int64
	mu            sync.RWMutex
	stopRequested bool
}

type PlaybackHandle struct {
	pb *playback
	am *AudioMixer
}

func (ph *PlaybackHandle) Stop() {
	ph.pb.mu.Lock()
	ph.pb.stopRequested = true
	ph.pb.mu.Unlock()

	<-ph.pb.done
}

func (ph *PlaybackHandle) Progress() float64 {
	ph.pb.mu.RLock()
	defer ph.pb.mu.RUnlock()

	if ph.pb.totalBytes == 0 {
		return 0.0
	}

	progress := float64(ph.pb.bytesRead) / float64(ph.pb.totalBytes)
	if progress > 1.0 {
		progress = 1.0
	}
	return progress
}

func (ph *PlaybackHandle) IsPlaying() bool {
	ph.am.mu.Lock()
	defer ph.am.mu.Unlock()

	_, exists := ph.am.playing[ph.pb.id]
	return exists
}

func (ph *PlaybackHandle) SetVolume(volume float64) {
	ph.pb.mu.Lock()
	defer ph.pb.mu.Unlock()
	ph.pb.volume = volume
}

func (ph *PlaybackHandle) GetVolume() float64 {
	ph.pb.mu.RLock()
	defer ph.pb.mu.RUnlock()
	return ph.pb.volume
}

func NewAudioMixer(channels int, mixAmp float64, framesPerWrite int, sampleRate int, bytesPerSample int) *AudioMixer {
	return &AudioMixer{
		playing:        make(map[int]*playback),
		channels:       channels,
		mixAmp:         mixAmp,
		framesPerWrite: framesPerWrite,
		sampleRate:     sampleRate,
		bytesPerSample: bytesPerSample,
		paused:         false,
	}
}

func (am *AudioMixer) BufferSize() int {
	return am.framesPerWrite * am.channels * am.bytesPerSample
}

func (am *AudioMixer) Period() time.Duration {
	return time.Duration(am.framesPerWrite) * time.Second / time.Duration(am.sampleRate)
}

func (am *AudioMixer) Pause() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.paused = true
}

func (am *AudioMixer) Resume() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.paused = false
}

func (am *AudioMixer) IsPaused() bool {
	am.mu.Lock()
	defer am.mu.Unlock()
	return am.paused
}

func (am *AudioMixer) TogglePause() bool {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.paused = !am.paused
	return am.paused
}

func (am *AudioMixer) Play(filePath string, volume float64) (*PlaybackHandle, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}

	fileInfo, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	pb := &playback{
		file:       f,
		done:       make(chan struct{}),
		buffer:     make([]byte, am.BufferSize()),
		volume:     volume,
		totalBytes: fileInfo.Size(),
		bytesRead:  0,
	}

	am.mu.Lock()
	pb.id = am.nextID
	am.nextID++
	am.playing[pb.id] = pb
	am.mu.Unlock()

	handle := &PlaybackHandle{
		pb: pb,
		am: am,
	}

	go func() {
		<-pb.done
		am.mu.Lock()
		delete(am.playing, pb.id)
		am.mu.Unlock()
		f.Close()
	}()

	return handle, nil
}

func (am *AudioMixer) FillBuffers() {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.paused {
		return
	}

	for id, pb := range am.playing {
		pb.mu.RLock()
		stopRequested := pb.stopRequested
		pb.mu.RUnlock()

		if stopRequested {
			close(pb.done)
			delete(am.playing, id)
			continue
		}

		pb.valid = 0
		for pb.valid < len(pb.buffer) {
			n, err := pb.file.Read(pb.buffer[pb.valid:])
			if n > 0 {
				pb.valid += n
				pb.mu.Lock()
				pb.bytesRead += int64(n)
				pb.mu.Unlock()
			}
			if err != nil {
				if err == io.EOF {
					close(pb.done)
					delete(am.playing, id)
					break
				}
				log.Errorf("error reading audio file: %v", err)
				close(pb.done)
				delete(am.playing, id)
				break
			}
		}
	}
}

func (am *AudioMixer) MixInto(buffer []byte) {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.paused || len(am.playing) == 0 {
		return
	}

	maxInt16 := float64(1<<15 - 1) // 32767

	for _, pb := range am.playing {
		if pb.valid == 0 {
			continue
		}

		framesToMix := min(pb.valid/(am.channels*am.bytesPerSample), am.framesPerWrite)

		pb.mu.RLock()
		effectiveVolume := am.mixAmp * pb.volume
		pb.mu.RUnlock()

		for frame := range framesToMix {
			frameOffset := frame * am.channels * am.bytesPerSample

			// Decode and mix left channel
			mixL := int16(binary.LittleEndian.Uint16(pb.buffer[frameOffset:]))
			mixLf := (float64(mixL) / maxInt16) * effectiveVolume

			origL := int16(binary.LittleEndian.Uint16(buffer[frameOffset:]))
			origLF := float64(origL) / maxInt16
			mixedLF := origLF + mixLf
			if mixedLF > 1.0 {
				mixedLF = 1.0
			} else if mixedLF < -1.0 {
				mixedLF = -1.0
			}
			mixedL := int16(mixedLF * maxInt16)
			binary.LittleEndian.PutUint16(buffer[frameOffset:], uint16(mixedL))

			// Decode and mix right channel
			mixR := int16(binary.LittleEndian.Uint16(pb.buffer[frameOffset+2:]))
			mixRf := (float64(mixR) / maxInt16) * effectiveVolume

			origR := int16(binary.LittleEndian.Uint16(buffer[frameOffset+2:]))
			origRF := float64(origR) / maxInt16
			mixedRF := origRF + mixRf
			if mixedRF > 1.0 {
				mixedRF = 1.0
			} else if mixedRF < -1.0 {
				mixedRF = -1.0
			}
			mixedR := int16(mixedRF * maxInt16)
			binary.LittleEndian.PutUint16(buffer[frameOffset+2:], uint16(mixedR))
		}
	}
}

func sendAudio(s ssh.Session, mixer *AudioMixer) {
	buffer := make([]byte, mixer.BufferSize())
	nextTick := time.Now()

	for {
		nextTick = nextTick.Add(mixer.Period())
		time.Sleep(time.Until(nextTick))

		clear(buffer)
		mixer.FillBuffers()
		mixer.MixInto(buffer)

		_, err := s.Write(buffer)
		if err != nil {
			log.Infof("stopped sending audio: %v", err)
			return
		}
	}
}
