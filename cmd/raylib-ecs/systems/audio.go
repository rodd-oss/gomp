/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package systems

import (
	"encoding/binary"
	"math"
	"slices"
	"time"

	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SampleRate = 48000
	BufferSize = 64 * time.Millisecond
)

var (
	SFXChannelId ecs.EntityID
)

type AudioController struct {
	stream rl.AudioStream
	// Supported values for sampleRate are 8000, 11025, 16000, 22050, 32000, 44100, 48000
	// Default: 48000
	SampleRate uint32
	// Supported values for BufferSize are 64, 128, 256, 512, 1024, 2048, 4096
	// Default: 64ms
	BufferSize time.Duration
	channels   *ecs.ComponentManager[components.AudioChannel]

	MasterChannelId     ecs.EntityID
	masterChannelBuffer []int16
	tempBuffer          []byte
}

func (c *AudioController) Init(world *ecs.World) {
	if c.SampleRate == 0 {
		c.SampleRate = SampleRate
	}

	if c.BufferSize == 0 {
		c.BufferSize = BufferSize
	}

	// TODO: make better validation
	// validate sampleRate
	if c.SampleRate != 8000 && c.SampleRate != 11025 && c.SampleRate != 16000 && c.SampleRate != 22050 && c.SampleRate != 32000 && c.SampleRate != 44100 && c.SampleRate != 48000 {
		panic("Unsupported sample rate")
	}

	// TODO: make better validation
	// validate bufferSize
	if c.BufferSize != 64*time.Millisecond && c.BufferSize != 128*time.Millisecond && c.BufferSize != 256*time.Millisecond && c.BufferSize != 512*time.Millisecond && c.BufferSize != 1024*time.Millisecond && c.BufferSize != 2048*time.Millisecond && c.BufferSize != 4096*time.Millisecond {
		panic("Unsupported buffer size")
	}

	// create master channel entity from where audio is played
	c.MasterChannelId = world.CreateEntity("Master audio channel")

	c.channels = components.AudioChannelService.GetManager(world)

	// create master channel
	masterChannel := c.channels.Create(c.MasterChannelId, components.AudioChannel{
		Volume:  0.1,
		Sources: make([]components.AudioStream, 0),
	})

	c.InitAudioChannels(world)

	c.channels.All(func(entityID ecs.EntityID, channel *components.AudioChannel) bool {
		if entityID == c.MasterChannelId {
			return true
		}

		masterChannel.Sources = append(masterChannel.Sources, channel)

		return true
	})

	c.tempBuffer = make([]byte, 4)

	rl.InitAudioDevice()

	c.stream = rl.LoadAudioStream(c.SampleRate, 16, 2)
	rl.SetAudioStreamCallback(c.stream, c.audioStreamCallback)
	rl.PlayAudioStream(c.stream)
}

func (c *AudioController) Update(world *ecs.World) {
}

func (c *AudioController) audioStreamCallback(data []float32, samplesCount int) {
	masterChannel := c.channels.Get(c.MasterChannelId)

	bufferLength := len(data) * 2

	if c.masterChannelBuffer == nil || bufferLength != len(c.masterChannelBuffer) {
		c.masterChannelBuffer = make([]int16, bufferLength)
	} else {
		clear(c.masterChannelBuffer)
	}

	_, err := masterChannel.Read(c.masterChannelBuffer)

	if err != nil {
		return
	}

	var dataIndex int
	for i := 0; i < len(c.masterChannelBuffer); i += 2 {
		binary.LittleEndian.PutUint16(c.tempBuffer, uint16(c.masterChannelBuffer[i]))
		binary.LittleEndian.PutUint16(c.tempBuffer[2:], uint16(c.masterChannelBuffer[i+1]))

		// end value of data[i] contains 4 bytes. First two bytes for left channel and second two bytes for right channel
		// value for each sample is int16 (-32768, 32767) that represents the amplitude of the sample
		// go bindings for raylib uses float32 and original raylib uses char array
		// since we are using 16bit 2 channel audio we have to store 2*sizeof(int16) bytes in float32 format
		data[dataIndex] = math.Float32frombits(binary.LittleEndian.Uint32(c.tempBuffer))
		dataIndex++
	}
}

func (c *AudioController) FixedUpdate(world *ecs.World) {
	audibles := components.AudibleService.GetManager(world)

	audibles.AllData(func(audible *components.Audible) bool {
		if audible.AudioStream != nil {
			channel := c.channels.Get(audible.ChannelId)
			streamIndex := slices.IndexFunc(channel.Sources, func(source components.AudioStream) bool { return source == audible })

			if streamIndex == -1 {
				channel.Sources = append(channel.Sources, audible)
			}
		}
		return true
	})
}

func (c *AudioController) Destroy(world *ecs.World) {
	rl.StopAudioStream(c.stream)
	rl.UnloadAudioStream(c.stream)

	// close master channel first
	c.CloseChannel(world, c.MasterChannelId)

	// close all channels
	c.channels.All(func(entityID ecs.EntityID, channel *components.AudioChannel) bool {
		c.CloseChannel(world, entityID)
		return true
	})

	rl.CloseAudioDevice()
}

func (c *AudioController) CloseChannel(world *ecs.World, entityID ecs.EntityID) {
	channels := components.AudioChannelService.GetManager(world)
	channel := channels.Get(entityID)

	if channel != nil {
		channel.Close()
	}

	channels.Remove(entityID)
}

func (c *AudioController) InitAudioChannels(world *ecs.World) {
	SFXChannelId = world.CreateEntity("SFX audio channel")
	c.channels.Create(SFXChannelId, components.AudioChannel{
		Volume:  1.0,
		Pan:     0.0,
		Sources: make([]components.AudioStream, 0),
	})
}
