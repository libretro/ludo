// Package audio uses OpenAL to play game audio by exposing the two audio
// callbacks Sample and SampleBatch for the libretro implementation.
package audio

// SetVolume sets the audio volume
func SetVolume(vol float32) {
}

// Init initializes the audio package. It opens the AL devices, sets the number of buffers, the
// volume and the source.
func Init(rate int32) {
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
