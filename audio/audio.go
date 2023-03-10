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
const ffMult = 0.25

var (
	paBuf      [bufSize]int32
	paSeBuf    [maxSeLen]int32
	paRate     int32
	paPtr      int64
	paPlayPtr  int64
	paSePtr    int
	paSeLen    int
	paCh       int
	paStream   *portaudio.Stream
	paSeStream *portaudio.Stream
	paUp       = false
)

// Effects are sound effects
var Effects map[string]*Effect

func st2mono(in [2]int16) [2]int16 {
	return [2]int16{in[0]/2 + in[1]/2, 0}
}

func paAudProc(in int32) [2]int16 {
	pi := (*[2]int16)(unsafe.Pointer(&in))
	var mt float32 = 1.0
	if state.FastForward {
		mt = ffMult
	}
	pi[0] = int16(float32(pi[0]) * settings.Current.AudioVolume * mt)
	pi[1] = int16(float32(pi[1]) * settings.Current.AudioVolume * mt)
	if paCh == 1 {
		return st2mono(*pi)
	}
	return *pi
}

func paSeProc(in int32) [2]int16 {
	pi := (*[2]int16)(unsafe.Pointer(&in))
	pi[0] = int16(float32(pi[0]) * settings.Current.MenuAudioVolume)
	pi[1] = int16(float32(pi[1]) * settings.Current.MenuAudioVolume)
	if paCh == 1 {
		return st2mono(*pi)
	}
	return *pi
}

// PortAudio Callback
func paCallback(out []int16) {
	for i := range out {
		if i%paCh == paCh-1 {
			if !state.MenuActive {
				if paPlayPtr < paPtr {
					s := paAudProc(paBuf[paPlayPtr-(paPlayPtr/bufSize)*bufSize])
					for j := 0; j < paCh; j++ {
						out[i-(paCh-1)+j] = s[j]
					}
					paPlayPtr++
					if paPtr-paPlayPtr < bufThreshold2 {
						// We have no choice but block pa here (can we speed up the core for a little?)
						time.Sleep(time.Millisecond * time.Duration(int64(bufBlock/int(paRate))))
					}
				} else {
					for j := 0; j < paCh; j++ {
						out[i-(paCh-1)+j] = 0
					}
				}
			} else {
				for j := 0; j < paCh; j++ {
					out[i-(paCh-1)+j] = 0
				}
			}
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
			log.Println("[PA]Output device is mono")
			p.Channels = out.MaxOutputChannels
		}
		paCh = p.Channels
		p.Latency = out.DefaultLowOutputLatency
	}
	p.SampleRate = float64(paRate)
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
	paSeStream, err = portaudio.OpenStream(NewParameters(h), func(out []int16) {
		for i := range out {
			if i%paCh == paCh-1 {
				if paSePtr < paSeLen {
					s := paSeProc(paSeBuf[paSePtr])
					for j := 0; j < paCh; j++ {
						out[i-(paCh-1)+j] = s[j]
					}
					paSePtr++
				} else {
					for j := 0; j < paCh; j++ {
						out[i-(paCh-1)+j] = 0
					}
				}
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
	if state.FastForward {
		bufOff = 0
	} else {
		blk := (paPtr - paPlayPtr) / bufThreshold
		if blk > 0 {
			log.Println("[PA]Core goes too fast, slowing down")
			time.Sleep(time.Millisecond * time.Duration(blk*int64(bufBlock/int(paRate))))
		}
	}

	mm := min(int(size/4), int(bufSize-bufOff))
	for i := 0; i < mm; i++ {
		p := 4 * i
		paBuf[paPtr-(paPtr/bufSize)*bufSize] = *(*int32)(unsafe.Pointer(&buf[p]))
		paPtr++
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
