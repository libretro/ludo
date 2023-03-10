package audio

import (
	"os"
	"unsafe"

	wav "github.com/youpy/go-wav"
)

// Effect is a static sound effect
type Effect struct {
	Format *wav.WavFormat
	paBuf  []int32
}

// LoadEffect loads a wav into memory and prepare the buffer
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

	var step int = int(e.Format.NumChannels) * 2
	samples := len(wav) / step
	e.paBuf = make([]int32, samples)
	for i := 0; i < samples; i++ {
		if e.Format.NumChannels == 2 {
			e.paBuf[i] = *(*int32)(unsafe.Pointer(&wav[step*i]))
		} else {
			var s *[2]int16
			s[0] = *(*int16)(unsafe.Pointer(&wav[step*i]))
			s[1] = s[0]
			e.paBuf[i] = *(*int32)(unsafe.Pointer(&s[0]))
		}
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
