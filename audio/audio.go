// Package audio uses OpenAL to play game audio by exposing the two audio
// callbacks Sample and SampleBatch for the libretro implementation.
package audio

import (
	"log"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"golang.org/x/mobile/exp/audio/al"
)

const bufSize = 1024 * 8

var audio struct {
	source     al.Source
	buffers    []al.Buffer
	rate       int32
	numBuffers int32
	tmpBuf     [bufSize]byte
	tmpBufPtr  int32
	resPtr     int32
}

// Effects are sound effects
var Effects map[string]*Effect

// SetVolume sets the audio volume
func SetVolume(vol float32) {
	audio.source.SetGain(vol)
}

// Init initializes the audio device
func Init() {
	err := al.OpenDevice()
	if err != nil {
		log.Println(err)
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
func Reconfigure(rate int32) {
	audio.rate = rate
	audio.numBuffers = 4

	log.Printf("[OpenAL]: Using %v buffers of %v bytes.\n", audio.numBuffers, bufSize)

	audio.source = al.GenSources(1)[0]
	audio.buffers = al.GenBuffers(int(audio.numBuffers))
	audio.resPtr = audio.numBuffers
	audio.tmpBufPtr = 0
	audio.tmpBuf = [bufSize]byte{}

	audio.source.SetGain(settings.Current.AudioVolume)
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func alUnqueueBuffers() bool {
	val := audio.source.BuffersProcessed()

	if val <= 0 {
		return false
	}

	audio.source.UnqueueBuffers(audio.buffers[audio.resPtr:val]...)
	audio.resPtr += val
	return true
}

func alGetBuffer() al.Buffer {
	if audio.resPtr == 0 {
		for {
			if alUnqueueBuffers() {
				break
			}
			time.Sleep(time.Millisecond)
		}
	}

	audio.resPtr--
	return audio.buffers[audio.resPtr]
}

func fillInternalBuf(buf []byte) int32 {
	readSize := min(bufSize-audio.tmpBufPtr, int32(len(buf)))
	copy(audio.tmpBuf[audio.tmpBufPtr:], buf[:readSize])
	audio.tmpBufPtr += readSize
	return readSize
}

func write(buf []byte, size int32) int32 {
	written := int32(0)

	if state.Global.FastForward {
		return size
	}

	for size > 0 {

		rc := fillInternalBuf(buf[written:])

		written += rc
		size -= rc

		if audio.tmpBufPtr != bufSize {
			break
		}

		buffer := alGetBuffer()

		buffer.BufferData(al.FormatStereo16, audio.tmpBuf[:], audio.rate)
		audio.tmpBufPtr = 0
		audio.source.QueueBuffers(buffer)

		if audio.source.State() != al.Playing {
			al.PlaySources(audio.source)
		}
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
