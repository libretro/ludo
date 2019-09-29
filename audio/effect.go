package audio

import (
	"os"

	wav "github.com/youpy/go-wav"
	"golang.org/x/mobile/exp/audio/al"
)

// Effect is a static sound effect
type Effect struct {
	Data   [4096]byte
	Format *wav.WavFormat
	source al.Source
	buffer al.Buffer
}

// LoadEffect loads a wav into memory and prepare the buffer and source in OpenAL
func LoadEffect(filename string) *Effect {
	var e Effect

	al.OpenDevice()
	e.source = al.GenSources(1)[0]
	e.buffer = al.GenBuffers(1)[0]

	file, _ := os.Open(filename)
	reader := wav.NewReader(file)

	e.Format, _ = reader.Format()
	defer file.Close()

	reader.Read(e.Data[:])

	e.buffer.BufferData(al.FormatStereo16, e.Data[:], int32(e.Format.SampleRate))
	e.source.QueueBuffers(e.buffer)

	return &e
}

// PlayEffect plays a sound effect
func PlayEffect(e *Effect) {
	al.PlaySources(e.source)
	for e.source.State() == al.Playing {
	}
}
