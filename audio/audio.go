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
	// "golang.org/x/mobile/exp/audio/al"
)

const bufSize = 1024 * 4
const maxSeLen = 44100 * 8

var (
	// source     al.Source
	// buffers    []al.Buffer
	// rate       int32
	// numBuffers int32
	// tmpBuf     [bufSize]byte
	// tmpBufPtr  int32
	// resPtr     int32
	paBuf     [bufSize]int32
	paSeBuf   [maxSeLen]int32
	paRate    float64
	paPtr     int64
	paPlayPtr int64
	paSePtr   int
	paSeLen   int
	// paHost     *portaudio.HostApiInfo
	paStream   *portaudio.Stream
	paSeStream *portaudio.Stream
)

// Effects are sound effects
var Effects map[string]*Effect

// SetVolume sets the audio volume
// func SetVolume(vol float32) {
// 	source.SetGain(vol)
// }

// PortAudio Callback
func paCallback(out []int32) {
	for i := range out {
		// if paMark[i] {
		// 	out[i] = paBuf[i]
		// 	paMark[i] = false
		if paPlayPtr <= paPtr {
			out[i] = int32(settings.Current.AudioVolume * float32(paBuf[paPlayPtr-(paPlayPtr/bufSize)*bufSize]))
			paPlayPtr++
			// if paPtr-paPlayPtr > bufSize/4 {
			// fmt.Println((paPtr - paPlayPtr) / (bufSize / 4))
			// }
			// fmt.Println(paStream.Info().SampleRate)
		} else {
			// out[i] = int32(rand.Int31n(100000000))
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
	// err := al.OpenDevice()
	// if err != nil {
	// 	log.Println(err)
	// }

	err1 := portaudio.Initialize()
	if err1 != nil {
		log.Println(err1)
	}

	paRate = 44100
	paPtr = 0
	paPlayPtr = 0
	paSePtr = 0
	paSeLen = 0

	// paHost, _ := portaudio.DefaultHostApi()
	// h, _ := portaudio.DefaultHostApi()
	h, _ := portaudio.DefaultOutputDevice()
	// paStream, _ = portaudio.OpenStream(portaudio.LowLatencyParameters(nil, h.DefaultOutputDevice), func(out []int32) {
	// 	for i := range out {
	// 		// if paMark[i] {
	// 		// 	out[i] = paBuf[i]
	// 		// 	paMark[i] = false
	// 		if paPlayPtr < paPtr {
	// 			out[i] = int32(settings.Current.AudioVolume * float32(paBuf[paPlayPtr-(paPlayPtr/bufSize)*bufSize]))
	// 			paPlayPtr++
	// 		} else {
	// 			out[i] = int32(rand.Int31n(100000000))
	// 			// out[i] = 0
	// 		}

	// 	}
	// })
	// paStream, _ = portaudio.OpenStream(portaudio.LowLatencyParameters(nil, h.DefaultOutputDevice), paCallback)
	// paStream, _ = portaudio.OpenStream(NewParameters(paHost.DefaultOutputDevice), paCallback)
	paStream, _ = portaudio.OpenStream(NewParameters(h), paCallback)
	paStream.Start()
	// fmt.Println(portaudio.LowLatencyParameters(nil, h.DefaultOutputDevice).FramesPerBuffer)

	paSeStream, _ = portaudio.OpenStream(NewParameters(h), func(out []int32) {
		for i := range out {
			if paSePtr < paSeLen {
				out[i] = int32(settings.Current.MenuAudioVolume * float32(paSeBuf[paSePtr]))
				paSePtr++
			} else {
				out[i] = 0
			}

		}
	})
	paSeStream.Start()
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
	// rate = r
	// numBuffers = 4

	// log.Printf("[OpenAL]: Using %v buffers of %v bytes.\n", numBuffers, bufSize)

	// source = al.GenSources(1)[0]
	// buffers = al.GenBuffers(int(numBuffers))
	// resPtr = numBuffers
	// tmpBufPtr = 0
	// tmpBuf = [bufSize]byte{}

	// source.SetGain(settings.Current.AudioVolume)

	paRate = float64(r)
	paBuf = [bufSize]int32{}
	paPtr = 0
	paPlayPtr = 0
	paStream.Close()
	h, _ := portaudio.DefaultOutputDevice()
	// paStream, _ = portaudio.OpenStream(NewParameters(paHost.DefaultOutputDevice), paCallback)
	paStream, _ = portaudio.OpenStream(NewParameters(h), paCallback)
	paStream.Start()
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// func alUnqueueBuffers() bool {
// 	val := source.BuffersProcessed()

// 	if val <= 0 {
// 		return false
// 	}

// 	source.UnqueueBuffers(buffers[resPtr:val]...)
// 	resPtr += val
// 	return true
// }

// func alGetBuffer() al.Buffer {
// 	if resPtr == 0 {
// 		for {
// 			if alUnqueueBuffers() {
// 				break
// 			}
// 			time.Sleep(time.Millisecond)
// 		}
// 	}

// 	resPtr--
// 	return buffers[resPtr]
// }

// func fillInternalBuf(buf []byte) int32 {
// 	readSize := min(bufSize-tmpBufPtr, int32(len(buf)))
// 	copy(tmpBuf[tmpBufPtr:], buf[:readSize])
// 	tmpBufPtr += readSize
// 	return readSize
// }

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

	// for size > 0 {

	// 	rc := fillInternalBuf(buf[written:])

	// 	written += rc
	// 	size -= rc

	// 	if tmpBufPtr != bufSize {
	// 		break
	// 	}

	// 	buffer := alGetBuffer()

	// 	buffer.BufferData(al.FormatStereo16, tmpBuf[:], rate)
	// 	tmpBufPtr = 0
	// 	source.QueueBuffers(buffer)

	// 	if source.State() != al.Playing {
	// 		al.PlaySources(source)
	// 	}
	// }

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
