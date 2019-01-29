// Package audio uses OpenAL to play game audio by exposing the two audio
// callbacks Sample and SampleBatch for the libretro implementation.
package audio

import (
	"log"

	"golang.org/x/mobile/exp/audio/al"
)

const bufSize = 1024 * 4

var audio struct {
	rate       int32
	numBuffers int32
	tmpBuf     [bufSize]byte
	tmpBufPtr  int32
	bufPtr     int32
	resPtr     int32
}

// SetVolume sets the audio volume
func SetVolume(vol float32) {
}

// Init initializes the audio package. It opens the AL devices, sets the number of buffers, the
// volume and the source.
func Init(rate int32) {
	err := al.OpenDevice()
	if err != nil {
		log.Println(err)
	}

	audio.rate = rate
	audio.numBuffers = 4

	log.Printf("[OpenAL]: Using %v buffers of %v bytes.\n", audio.numBuffers, bufSize)

	audio.resPtr = audio.numBuffers
	audio.tmpBufPtr = 0
	audio.tmpBuf = [bufSize]byte{}
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func write(buf []byte, size int32) int32 {
	written := int32(0)
	return written
}

// Sample renders a single audio frame.
// It is passed as a callback to the libretro implementation.
func Sample(left int16, right int16) {
	buf := []byte{byte(left), byte(right)}
	write(buf, 4)
}

// SampleBatch renders multiple audio frames in one go
// It is passed as a callback to the libretro implementation.
func SampleBatch(buf []byte, size int32) int32 {
	return write(buf, size*4)
}
