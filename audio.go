package main

const bufSize = 1024 * 4

var audio struct {
	rate       int32
	numBuffers int32
	tmpBuf     [bufSize]byte
	tmpBufPtr  int32
	bufPtr     int32
	resPtr     int32
}

func audioSetVolume(vol float32) {
}

func audioInit(rate int32) {
}

func audioWrite(buf []byte, size int32) int32 {
	return 0
}

func audioSample(left int16, right int16) {
}

func audioSampleBatch(buf []byte, size int32) int32 {
	return 0
}
