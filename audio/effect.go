package audio

import (
	"os"

	wav "github.com/youpy/go-wav"
	"golang.org/x/mobile/exp/audio/al"
)

// Effect is a static sound effect
type Effect struct {
	Format *wav.WavFormat
	source al.Source
	buffer al.Buffer
}

// LoadEffect loads a wav into memory and prepare the buffer and source in OpenAL
func LoadEffect(filename string) (*Effect, error) {
	var e Effect
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

	wav := []byte{}
	for {
		var data [4096]byte
		n, _ := reader.Read(data[:])
		if n == 0 {
			break
		}
		wav = append(wav, data[:]...)
	}

	e.buffer.BufferData(al.FormatMono16, wav, int32(e.Format.SampleRate))
	e.source.QueueBuffers(e.buffer)

	return &e, nil
}

// PlayEffect plays a sound effect, blocking
func PlayEffect(e *Effect) {
	al.PlaySources(e.source)
	for e.source.State() == al.Playing {
	}
}
