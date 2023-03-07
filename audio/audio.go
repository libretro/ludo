// Package audio uses OpenAL to play game audio by exposing the two audio
// callbacks Sample and SampleBatch for the libretro implementation.
package audio

import (
	"encoding/binary"
	"log"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/gordonklaus/portaudio"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

const bufSize = 1024 * 4
const maxSeLen = 44100 * 8

var (
	paBuf      [bufSize]int32
	paSeBuf    [maxSeLen]int32
	paRate     float64
	paPtr      int64
	paPlayPtr  int64
	paSePtr    int
	paSeLen    int
	paStream   *portaudio.Stream
	paSeStream *portaudio.Stream
	paUp       = false
)

// Effects are sound effects
var Effects map[string]*Effect

// PortAudio Callback
func paCallback(out []int32) {
	for i := range out {
		if paPlayPtr <= paPtr {
			out[i] = int32(settings.Current.AudioVolume * float32(paBuf[paPlayPtr-(paPlayPtr/bufSize)*bufSize]))
			paPlayPtr++
		} else {
			out[i] = 0
		}

	}
}

// Create PortAudio parameters
func NewParameters(out *portaudio.DeviceInfo) (p portaudio.StreamParameters) {
	if out != nil {
		p := &p.Output
		p.Device = out
		p.Channels = 2
		if out.MaxOutputChannels < 2 {
			p.Channels = out.MaxOutputChannels
		}
		p.Latency = out.DefaultLowOutputLatency
	}
	p.SampleRate = paRate / 2
	p.FramesPerBuffer = portaudio.FramesPerBufferUnspecified
	return p
}

// Init initializes the audio device
func Init() {
	if paUp {
		return
	}
	err1 := portaudio.Initialize()
	if err1 != nil {
		log.Println(err1)
	}

	paRate = 44100
	paPtr = 0
	paPlayPtr = 0
	paSePtr = 0
	paSeLen = 0

	h, err := portaudio.DefaultOutputDevice()
	if err != nil {
		log.Fatalln(err)
	}
	paStream, err = portaudio.OpenStream(NewParameters(h), paCallback)
	if err != nil {
		log.Fatalln(err)
	}
	if err = paStream.Start(); err != nil {
		log.Fatalln(err)
	}
	paSeStream, err = portaudio.OpenStream(NewParameters(h), func(out []int32) {
		for i := range out {
			if paSePtr < paSeLen {
				out[i] = int32(settings.Current.MenuAudioVolume * float32(paSeBuf[paSePtr]))
				paSePtr++
			} else {
				out[i] = 0
			}

		}
	})
	if err != nil {
		log.Fatalln(err)
	}
	if err = paSeStream.Start(); err != nil {
		log.Fatalln(err)
	}
	Effects = map[string]*Effect{}

	assets := settings.Current.AssetsDirectory
	paths, _ := filepath.Glob(assets + "/sounds/*.wav")
	for _, path := range paths {
		path := path
		filename := utils.FileName(path)
		Effects[filename], _ = LoadEffect(path)
	}
}

// Reconfigure initializes the audio package. It sets the number of buffers, the
// volume and the source for the games.
func Reconfigure(r int32) {
	paRate = float64(r)
	paBuf = [bufSize]int32{}
	paPtr = 0
	paPlayPtr = 0
	if paStream != nil {
		if err := paStream.Close(); err != nil {
			log.Fatalln(err)
		}
	}
	h, err := portaudio.DefaultOutputDevice()
	if err != nil {
		log.Fatalln(err)
	}
	paStream, err = portaudio.OpenStream(NewParameters(h), paCallback)
	if err != nil {
		log.Fatalln(err)
	}
	if err = paStream.Start(); err != nil {
		log.Fatalln(err)
	}
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func write(buf []byte, size int32) int32 {
	written := int32(0)

	if state.FastForward {
		return size
	}

	time.Sleep(time.Millisecond * time.Duration((paPtr-paPlayPtr)/(bufSize/4)*8))

	mm := int(min(size/4, bufSize))
	for i := 0; i < mm; i++ {
		p := 4 * (int32(i))
		paBuf[paPtr-(paPtr/bufSize)*bufSize] = int32(binary.LittleEndian.Uint32(buf[p : p+4]))
		paPtr++
		written += 4
	}

	return written
}

// Sample renders a single audio frame.
// It is passed as a callback to the libretro implementation.
func Sample(left int16, right int16) {
	// simulate the kind of raw byte array that would be provided from C via SampleBatch.
	// (effectively typecasting int16 array to byte array)
	buf := []int16{left, right}
	pi := (*[4]byte)(unsafe.Pointer(&buf[0]))
	write((*pi)[:], 4)
}

// SampleBatch renders multiple audio frames in one go
// It is passed as a callback to the libretro implementation.
func SampleBatch(buf []byte, size int32) int32 {
	return write(buf, size*4)
}
