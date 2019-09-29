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
func LoadEffect(filename string) (*Effect, error) {
	var e Effect

	al.OpenDevice() // Move this to an init

	e.source = al.GenSources(1)[0]
	e.buffer = al.GenBuffers(1)[0]

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := wav.NewReader(file)

	e.Format, err = reader.Format()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = reader.Read(e.Data[:])
	if err != nil {
		return nil, err
	}

	e.buffer.BufferData(al.FormatStereo16, e.Data[:], int32(e.Format.SampleRate))
	e.source.QueueBuffers(e.buffer)

	return &e, nil
}

// PlayEffect plays a sound effect
func PlayEffect(e *Effect) {
	al.PlaySources(e.source)
	for e.source.State() == al.Playing {
	}
}
