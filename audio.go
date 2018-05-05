package main

import (
	"fmt"
	"time"

	"golang.org/x/mobile/exp/audio/al"
)

const bufSize = 1024 * 4

var audio struct {
	source     al.Source
	buffers    []al.Buffer
	rate       int32
	numBuffers int32
	tmpBuf     [bufSize]byte
	tmpBufPtr  int32
	bufPtr     int32
	resPtr     int32
}

func audioInit(rate int32) {
	err := al.OpenDevice()
	if err != nil {
		fmt.Println(err)
	}

	audio.rate = rate
	audio.numBuffers = 4

	fmt.Printf("[OpenAL]: Using %v buffers of %v bytes.\n", audio.numBuffers, bufSize)

	audio.source = al.GenSources(1)[0]
	audio.buffers = al.GenBuffers(int(audio.numBuffers))
	audio.resPtr = audio.numBuffers

	audio.source.SetGain(settings.AudioVolume)
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

func fillInternalBuf(buf []byte, size int32) int32 {
	readSize := min(bufSize-audio.tmpBufPtr, size)
	copy(audio.tmpBuf[audio.tmpBufPtr:], buf[audio.bufPtr:audio.bufPtr+readSize])
	audio.tmpBufPtr += readSize
	return readSize
}

func audioWrite(buf []byte, size int32) int32 {
	written := int32(0)

	for size > 0 {

		rc := fillInternalBuf(buf, size)

		written += rc
		audio.bufPtr += rc
		size -= rc

		if audio.tmpBufPtr != bufSize {
			break
		}

		buffer := alGetBuffer()

		buffer.BufferData(al.FormatStereo16, audio.tmpBuf[:], int32(audio.rate))
		audio.tmpBufPtr = 0
		audio.source.QueueBuffers(buffer)

		if audio.source.State() != al.Playing {
			al.PlaySources(audio.source)
		}
	}

	audio.bufPtr = 0

	return written
}

func audioSample(left int16, right int16) {
	buf := []byte{byte(left), byte(right)}
	audioWrite(buf, 4)
}

func audioSampleBatch(buf []byte, size int32) int32 {
	return audioWrite(buf, size*4)
}
