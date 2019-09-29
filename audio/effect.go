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
}

var menu struct {
	source  al.Source
	buffers []al.Buffer
}

func LoadEffect(filename string) *Effect {
	al.OpenDevice()
	menu.source = al.GenSources(1)[0]
	menu.buffers = al.GenBuffers(1)

	file, _ := os.Open(filename)
	reader := wav.NewReader(file)

	var e Effect
	e.Format, _ = reader.Format()
	defer file.Close()

	reader.Read(e.Data[:])

	return &e
}

func PlayEffect(e *Effect) {
	menu.buffers[0].BufferData(al.FormatStereo16, e.Data[:], int32(e.Format.SampleRate))
	menu.source.QueueBuffers(menu.buffers[0])
	al.PlaySources(menu.source)
	for menu.source.State() == al.Playing {

	}
	menu.source.UnqueueBuffers(menu.buffers[0])
}
