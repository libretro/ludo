package main

import (
	"C"
	"fmt"

	"golang.org/x/mobile/exp/audio/al"

	"time"
	"unsafe"
)

/*
#include <stdlib.h>
*/
import "C"

const bufSize = 1024 * 4

var audio struct {
	source     al.Source
	buffers    []al.Buffer
	rate       int32
	numBuffers int32
	tmpBuf     [bufSize]byte
	tmpBufPtr  C.size_t
	bufPtr     C.size_t
	resPtr     int32
}

func audioInit(rate C.double) {
	err := al.OpenDevice()
	if err != nil {
		fmt.Println(err)
	}

	audio.rate = int32(rate)
	audio.numBuffers = 4

	fmt.Printf("[OpenAL]: Using %v buffers of %v bytes.\n", audio.numBuffers, bufSize)

	audio.source = al.GenSources(1)[0]
	audio.buffers = al.GenBuffers(int(audio.numBuffers))
	audio.resPtr = audio.numBuffers
}

func min(a, b C.size_t) C.size_t {
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

func alGetBuffer() (al.Buffer, error) {
	if audio.resPtr == 0 {
		for {
			if alUnqueueBuffers() {
				break
			}

			// if audio.nonblock
			//   return nil, true

			/* Must sleep as there is no proper blocking method. */
			time.Sleep(time.Millisecond)
		}
	}

	audio.resPtr--
	return audio.buffers[audio.resPtr], nil
}

func fillInternalBuf(buf unsafe.Pointer, size C.size_t) C.size_t {
	readSize := min(bufSize-audio.tmpBufPtr, size)
	copy(audio.tmpBuf[audio.tmpBufPtr:], C.GoBytes(buf, bufSize)[audio.bufPtr:audio.bufPtr+readSize])
	audio.tmpBufPtr += readSize
	return readSize
}

func audioWrite(buf unsafe.Pointer, size C.size_t) C.size_t {
	written := C.size_t(0)

	for size > 0 {

		rc := fillInternalBuf(buf, size)

		written += rc
		audio.bufPtr += rc
		size -= rc

		if audio.tmpBufPtr != bufSize {
			break
		}

		buffer, err := alGetBuffer()
		if err != nil {
			break
		}

		buffer.BufferData(al.FormatStereo16, audio.tmpBuf[:], audio.rate)
		audio.tmpBufPtr = 0
		audio.source.QueueBuffers(buffer)

		if audio.source.State() != al.Playing {
			al.PlaySources(audio.source)
		}
	}

	audio.bufPtr = 0

	return written
}

//export coreAudioSample
func coreAudioSample(left C.int16_t, right C.int16_t) {
	buf := []C.int16_t{left, right}
	audioWrite(unsafe.Pointer(&buf), 4)
}

//export coreAudioSampleBatch
func coreAudioSampleBatch(data unsafe.Pointer, frames C.size_t) C.size_t {
	return audioWrite(data, frames*4)
}
