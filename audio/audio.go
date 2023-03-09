// Package audio uses OpenAL to play game audio by exposing the two audio
// callbacks Sample and SampleBatch for the libretro implementation.
package audio

import (
	"log"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/gordonklaus/portaudio"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
)

const bufSize = 256 * 32
const bufThreshold = 256 * 24
const bufThreshold2 = 256 * 8
const bufBlock = 256 * 8 * 1000
const maxSeLen = 44100 * 8

var (
	paBuf      [bufSize]int32
	paSeBuf    [maxSeLen]int32
	paRate     int32
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

func paInterleave(in uint32) uint32 {
	var bin uint32 = 0b11111111111111111111111111111111
	return in & bin
}

// PortAudio Callback
// func paCallback(out [][]int16) {
func paCallback(out []int32) {
	// log.Println("X", len(out), len(out[0]))
	for i := range out {
		if !state.MenuActive {
			if paPlayPtr < paPtr {
				out[i] = paBuf[paPlayPtr-(paPlayPtr/bufSize)*bufSize] // I'll do volume stuff later
				// out[i] = int32(paInterleave(uint32(settings.Current.AudioVolume * float32(paBuf[paPlayPtr-(paPlayPtr/bufSize)*bufSize]))))
				// out[0][i] = int16(paInterleave(uint32(settings.Current.AudioVolume*float32(paBuf[paPlayPtr-(paPlayPtr/bufSize)*bufSize]))) << 16)
				// out[1][i] = int16(paInterleave(uint32(settings.Current.AudioVolume*float32(paBuf[paPlayPtr-(paPlayPtr/bufSize)*bufSize]))) & 0xFFFF)
				paPlayPtr++
				if paPtr-paPlayPtr < bufThreshold2 {
					// We have no choice but block pa here (can we speed up the core for a little?)
					time.Sleep(time.Millisecond * time.Duration(int64(bufBlock/int(paRate))))
				}
			} else {
				out[i] = 0
				// out[0][i] = 0
				// out[1][i] = 0
			}
		} else {
			out[i] = 0
			// out[0][i] = 0
			// out[1][i] = 0
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
			log.Println("[PA]Output is mono")
			p.Channels = out.MaxOutputChannels
		}
		p.Latency = out.DefaultLowOutputLatency
	}
	// p.SampleRate = float64(paRate)
	p.SampleRate = float64(paRate / int32(p.Output.Channels))
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
				out[i] = int32(settings.Current.MenuAudioVolume * float32(int32(paSeBuf[paSePtr])))
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
	paRate = r
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func write(buf []byte, size int32) int32 {
	written := int32(0)

	bufOff := paPtr - paPlayPtr
	// bufDiv := 1
	if state.FastForward {
		bufOff = 0
		// bufDiv = 4
	} else {
		blk := (paPtr - paPlayPtr) / bufThreshold
		if blk > 0 {
			log.Println("[PA]Core goes too fast, slowing down")
			time.Sleep(time.Millisecond * time.Duration(blk*int64(bufBlock/int(paRate))))
		}
	}

	mm := min(int(size/4), int(bufSize-bufOff))
	for i := 0; i < mm; i++ {
		p := 4 * (int32(i))
		// paBuf[paPtr-(paPtr/bufSize)*bufSize] = binary.LittleEndian.Uint32(buf[p:p+4]) / uint32(bufDiv)
		// paBuf[paPtr-(paPtr/bufSize)*bufSize] = int32(binary.LittleEndian.Uint32(buf[p:p+4])) / int32(bufDiv)
		// I wonder how the stereo audio data is organized... Should I change to this?
		// ps := (*int32)(unsafe.Pointer(&buf[p])) // / int32(bufDiv) // we have to change vol in another way
		// paBuf[paPtr-(paPtr/bufSize)*bufSize] = *ps // / int32(bufDiv) // we have to change vol in another way
		paBuf[paPtr-(paPtr/bufSize)*bufSize] = *(*int32)(unsafe.Pointer(&buf[p])) // / int32(bufDiv) // we have to change vol in another way
		written += 4
	}

	// Reset ptr into single range
	for {
		if paPtr-paPlayPtr < bufSize*2 {
			break
		}
		log.Println("[PA]Buffer bit the tail! Are we fast-forwarding?")
		paPtr -= int64(bufSize)
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
