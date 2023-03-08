package audio

import (
	"encoding/binary"
	"os"

	wav "github.com/youpy/go-wav"
)

// Effect is a static sound effect
type Effect struct {
	Format *wav.WavFormat
	paBuf  []int32
}

// LoadEffect loads a wav into memory and prepare the buffer and source in OpenAL
func LoadEffect(filename string) (*Effect, error) {
	var e Effect
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

	samples := len(wav) / 4
	e.paBuf = make([]int32, samples)
	for i := 0; i < samples; i++ {
		p := 4 * i
		e.paBuf[i] = int32(binary.LittleEndian.Uint32(wav[p : p+4]))
	}

	return &e, nil
}

// PlayEffect plays a sound effect
func PlayEffect(e *Effect) {
	if len(e.paBuf) > 0 && len(e.paBuf) < maxSeLen {
		paSeBuf = [maxSeLen]int32{}
		copy(paSeBuf[:], e.paBuf)
		paSeLen = len(e.paBuf)
		paSePtr = 0
	}
}
