package main

import "golang.org/x/mobile/exp/audio/al"

func audioSetVolume(vol float32) {
}

func audioInit(rate int32) {
	al.OpenDevice()
}

func audioSample(left int16, right int16) {
}

func audioSampleBatch(buf []byte, size int32) int32 {
	return 0
}
