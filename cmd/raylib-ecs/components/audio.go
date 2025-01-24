/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package components

import (
	"gomp_game/pkgs/gomp/ecs"
	"io"
	"math"
	"time"
)

type AudioState uint

const (
	Stoped AudioState = iota
	Paused
	Playing
	Looping
)

type AudioStream interface {
	Read(samples []int16) (int, error)
	Close()
}

type Audible struct {
	State     AudioState
	ChannelId ecs.EntityID
	Volume    float64
	Pan       float64
	AudioStream
}

func (a *Audible) Read(samples []int16) (int, error) {
	if a.State == Stoped || a.State == Paused {
		return 0, io.EOF
	}

	a.State = Playing

	return a.AudioStream.Read(samples)
}

func (a *Audible) Close() {
	if a.AudioStream != nil {
		a.AudioStream.Close()
	}

	a.State = Stoped
}

type SineWaveGenerator struct {
	SampleFactor float64
	SampleRate   uint32
	Phase        float64
	Duration     time.Duration
	playedTime   time.Duration
}

func (s *SineWaveGenerator) Read(samples []int16) (int, error) {
	if s.Duration != 0 && s.playedTime >= s.Duration {
		return 0, io.EOF
	}

	var samplesCount int64 = int64(len(samples)) / 2 // 2 channels
	var samplesDuration time.Duration = time.Duration((samplesCount / int64(s.SampleRate))) * time.Second

	if s.Duration == 0 {
		durationLeft := s.Duration - s.playedTime
		if durationLeft < samplesDuration {
			samplesCount = int64(durationLeft.Seconds() * float64(s.SampleRate))
		}
	}

	for i := 0; i < int(samplesCount*2); i += 2 {
		v := math.Sin(s.Phase * 2.0 * math.Pi)

		value := int16(v * math.MaxInt16)

		samples[i] = value
		samples[i+1] = samples[i]

		_, s.Phase = math.Modf(s.Phase + s.SampleFactor)
	}

	s.playedTime += samplesDuration

	return len(samples), nil
}

func (s *SineWaveGenerator) Close() {}

// Channel is stereo pcm `io.ReadCloser` that supports paning and volume control
type AudioChannel struct {
	IsClosed bool
	// Pan is in range [-1.0, 1.0]. Panics otherwise
	Pan float64
	// Volume is in range [0.0, 1.0]. Panics otherwise
	Volume       float64
	Sources      []AudioStream
	mixingBuffer []int16
	AudioStream
}

func (c *AudioChannel) Read(samples []int16) (n int, err error) {
	if c.IsClosed {
		return 0, io.EOF
	}

	bufferLength := len(samples)

	if c.mixingBuffer == nil || bufferLength != len(c.mixingBuffer) {
		c.mixingBuffer = make([]int16, bufferLength)
	} else {
		clear(c.mixingBuffer)
	}

	if c.Pan < -1.0 || c.Pan > 1.0 {
		panic("Pan must be in range [-1.0, 1.0]")
	}

	if c.Volume < 0.0 || c.Volume > 1.0 {
		panic("Volume must be in range [0.0, 1.0]")
	}

	for _, source := range c.Sources {
		clear(c.mixingBuffer)
		n, err := source.Read(c.mixingBuffer)

		if err != nil {
			return n, err
		}

		for i := 0; i < bufferLength; i += 2 {
			// mix left and right channels
			samples[i] += c.mixingBuffer[i]
			samples[i+1] += c.mixingBuffer[i+1]
		}

	}

	for i := 0; i < len(samples); i += 2 {
		// convert to float to apply volume and pan
		leftFloat := float64(samples[i])
		rightFloat := float64(samples[i+1])

		// apply volume
		leftFloat *= c.Volume
		rightFloat *= c.Volume

		// apply pan
		if c.Pan < 0.0 {
			rightFloat *= 1.0 - -c.Pan
		} else if c.Pan > 0.0 {
			leftFloat *= 1.0 - c.Pan
		}

		// update p
		samples[i] = int16(leftFloat)
		samples[i+1] = int16(rightFloat)
	}

	return len(samples), nil
}

func (c *AudioChannel) Close() {
	c.IsClosed = true
}
