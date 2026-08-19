// Package audio uses OpenAL to play game audio by exposing the two audio
// callbacks Sample and SampleBatch for the libretro implementation.
package audio

import (
	"encoding/binary"
	"log"
	"math"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"golang.org/x/mobile/exp/audio/al"
)

const (
	bufSize             = 1024
	bytesPerFrame       = 4
	framesPerBuffer     = bufSize / bytesPerFrame
	outputRate          = 48000
	targetLatencyMs     = 96
	minNumBuffers       = 2
	audioMaxTimingSkew  = 0.05
	maxPlaybackRateSkew = 0.005
)

var (
	source        al.Source
	buffers       []al.Buffer
	freeBufs      []al.Buffer
	rate          int32
	inputRate     float64
	videoSyncRate float64
	numBuffers    int32
	queueBytes    int32
	inputBuf      []int16
	readPhase     float64
	tmpBuf        [bufSize]byte
	srcRatioOrig  float64
	srcRatioCurr  float64
)

// Effects are sound effects
var Effects map[string]*Effect

// SetVolume sets the audio volume
func SetVolume(vol float32) {
	if source == 0 {
		return
	}
	source.SetGain(vol)
}

// Stop clears queued game audio without tearing down the OpenAL device.
func Stop() {
	clearQueue()
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
func Reconfigure(r int32) {
	releaseSource()

	rate = outputRate
	inputRate = float64(r)
	numBuffers = queueBufferCount(rate)
	queueBytes = numBuffers * bufSize
	updateInputRate()

	log.Printf(
		"[OpenAL]: Using %v buffers of %v bytes (~%v ms queued).\n",
		numBuffers,
		bufSize,
		queueLatencyMs(rate, queueBytes),
	)

	source = al.GenSources(1)[0]
	buffers = al.GenBuffers(int(numBuffers))
	freeBufs = append(freeBufs[:0], buffers...)
	inputBuf = inputBuf[:0]
	readPhase = 0
	tmpBuf = [bufSize]byte{}

	source.SetGain(settings.Current.AudioVolume)
}

func releaseSource() {
	if source == 0 {
		return
	}

	clearQueue()
	al.DeleteSources(source)
	source = 0

	if len(buffers) > 0 {
		al.DeleteBuffers(buffers...)
	}

	buffers = nil
	freeBufs = nil
	queueBytes = 0
	inputBuf = nil
	readPhase = 0
	tmpBuf = [bufSize]byte{}
	inputRate = 0
	videoSyncRate = 0
	srcRatioOrig = 0
	srcRatioCurr = 0
}

func clearQueue() {
	if source == 0 {
		return
	}

	al.StopSources(source)

	queued := source.BuffersQueued()
	if queued > 0 {
		unqueued := make([]al.Buffer, queued)
		source.UnqueueBuffers(unqueued...)
	}

	freeBufs = append(freeBufs[:0], buffers...)
	inputBuf = inputBuf[:0]
	readPhase = 0
	tmpBuf = [bufSize]byte{}
	srcRatioCurr = srcRatioOrig
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func clampFloat64(v, low, high float64) float64 {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func queueBufferCount(rate int32) int32 {
	totalBytes := int32(math.Ceil(float64(rate*bytesPerFrame*targetLatencyMs) / 1000.0))
	buffers := int32(math.Ceil(float64(totalBytes) / float64(bufSize)))

	if buffers < minNumBuffers {
		return minNumBuffers
	}
	return buffers
}

func adjustedInputRate(inputRate, inputFPS, targetVideoSyncRate float64) float64 {
	if inputRate <= 0 {
		return 0
	}
	if inputFPS <= 0 || targetVideoSyncRate <= 0 {
		return inputRate
	}

	timingSkew := math.Abs(1.0 - inputFPS/targetVideoSyncRate)
	if timingSkew <= audioMaxTimingSkew {
		return inputRate * targetVideoSyncRate / inputFPS
	}

	return inputRate
}

func updateInputRate() {
	if inputRate <= 0 || rate <= 0 {
		srcRatioOrig = 0
		srcRatioCurr = 0
		return
	}

	syncInputRate := adjustedInputRate(inputRate, state.CoreFPS, videoSyncRate)
	if syncInputRate <= 0 {
		syncInputRate = inputRate
	}

	srcRatioOrig = float64(rate) / syncInputRate
	srcRatioCurr = srcRatioOrig
}

// SetVideoTiming updates the current video timing used for audio rate control.
func SetVideoTiming(refreshHz float64, swapInterval int) {
	syncRate := refreshHz
	if swapInterval > 1 {
		syncRate /= float64(swapInterval)
	}

	if math.Abs(syncRate-videoSyncRate) < 0.01 {
		return
	}

	videoSyncRate = syncRate
	updateInputRate()
}

func queueLatencyMs(rate, bytes int32) int32 {
	if rate <= 0 {
		return 0
	}
	return int32(math.Round(float64(bytes*1000) / float64(rate*bytesPerFrame)))
}

func pendingQueueBytes(queued, processed, offset int32) int32 {
	active := queued - processed
	if active <= 0 {
		return 0
	}

	bytes := active * bufSize
	if offset > 0 {
		bytes -= min(offset, bufSize)
	}
	if bytes < 0 {
		return 0
	}
	return bytes
}

func queuedWholeBufferBytes() int32 {
	if source == 0 {
		return 0
	}

	queued := source.BuffersQueued()
	processed := source.BuffersProcessed()
	return pendingQueueBytes(queued, processed, 0)
}

func sourceWriteAvail() int32 {
	alUnqueueBuffers()

	avail := queueBytes - queuedWholeBufferBytes()
	if avail < 0 {
		return 0
	}
	if avail > queueBytes {
		return queueBytes
	}
	return avail
}

func currentResampleStep(writeAvail int32) float64 {
	if rate <= 0 || queueBytes <= 0 || srcRatioOrig <= 0 {
		return 1.0
	}

	halfSize := float64(queueBytes) / 2
	direction := 0.0
	if halfSize > 0 {
		direction = (float64(writeAvail) - halfSize) / halfSize
	}

	adjust := 1.0 + maxPlaybackRateSkew*clampFloat64(direction, -1.0, 1.0)
	srcRatioCurr = srcRatioOrig * adjust
	return 1.0 / srcRatioCurr
}

func alUnqueueBuffers() int32 {
	if source == 0 {
		return 0
	}

	val := source.BuffersProcessed()

	if val <= 0 {
		return 0
	}

	unqueued := make([]al.Buffer, val)
	source.UnqueueBuffers(unqueued...)
	freeBufs = append(freeBufs, unqueued...)
	return val
}

func alGetBuffer() al.Buffer {
	for len(freeBufs) == 0 {
		if alUnqueueBuffers() > 0 {
			break
		}

		// OpenAL does not provide a blocking dequeue API, so keep the queue
		// short and sleep briefly while waiting for the next buffer to retire.
		time.Sleep(time.Millisecond)
	}

	last := len(freeBufs) - 1
	buffer := freeBufs[last]
	freeBufs = freeBufs[:last]
	return buffer
}

func appendInput(buf []byte) int32 {
	readSize := int32(len(buf))
	for i := 0; i+1 < len(buf); i += 2 {
		inputBuf = append(inputBuf, int16(binary.LittleEndian.Uint16(buf[i:])))
	}
	return readSize
}

func availableInputFrames() int {
	return len(inputBuf) / 2
}

func requiredInputFrames(step float64) int {
	return int(math.Ceil(readPhase+step*float64(framesPerBuffer-1))) + 2
}

func renderOutputBuffer(step float64) {
	for frame := 0; frame < framesPerBuffer; frame++ {
		base := int(readPhase)
		frac := readPhase - float64(base)

		left0 := float64(inputBuf[base*2])
		right0 := float64(inputBuf[base*2+1])
		left1 := float64(inputBuf[(base+1)*2])
		right1 := float64(inputBuf[(base+1)*2+1])

		left := int16(math.Round(left0 + (left1-left0)*frac))
		right := int16(math.Round(right0 + (right1-right0)*frac))

		offset := frame * bytesPerFrame
		binary.LittleEndian.PutUint16(tmpBuf[offset:], uint16(left))
		binary.LittleEndian.PutUint16(tmpBuf[offset+2:], uint16(right))

		readPhase += step
	}

	discardFrames := int(readPhase)
	if discardFrames <= 0 {
		return
	}

	discardSamples := discardFrames * 2
	copy(inputBuf, inputBuf[discardSamples:])
	inputBuf = inputBuf[:len(inputBuf)-discardSamples]
	readPhase -= float64(discardFrames)
}

func write(buf []byte, size int32) int32 {
	if state.FastForward {
		clearQueue()
		return size
	}

	if source == 0 || len(buffers) == 0 {
		return size
	}

	written := appendInput(buf[:size])

	for {
		writeAvail := sourceWriteAvail()
		step := currentResampleStep(writeAvail)
		if availableInputFrames() < requiredInputFrames(step) {
			break
		}
		buffer := alGetBuffer()
		renderOutputBuffer(step)
		buffer.BufferData(al.FormatStereo16, tmpBuf[:], rate)
		source.QueueBuffers(buffer)

		if source.State() != al.Playing {
			al.PlaySources(source)
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
