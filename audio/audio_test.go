package audio

import (
	"math"
	"testing"
)

func TestPendingQueueBytes(t *testing.T) {
	tests := []struct {
		name      string
		queued    int32
		processed int32
		offset    int32
		want      int32
	}{
		{
			name:      "Nothing queued",
			queued:    0,
			processed: 0,
			offset:    0,
			want:      0,
		},
		{
			name:      "Processed buffers do not count as pending",
			queued:    4,
			processed: 2,
			offset:    0,
			want:      2 * bufSize,
		},
		{
			name:      "Subtract current playback offset",
			queued:    3,
			processed: 1,
			offset:    128,
			want:      2*bufSize - 128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pendingQueueBytes(tt.queued, tt.processed, tt.offset)
			if got != tt.want {
				t.Errorf("pendingQueueBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueueBufferCount(t *testing.T) {
	tests := []struct {
		rate int32
		want int32
	}{
		{rate: 32040, want: 13},
		{rate: 48000, want: 18},
		{rate: 96000, want: 36},
	}

	for _, tt := range tests {
		if got := queueBufferCount(tt.rate); got != tt.want {
			t.Errorf("queueBufferCount(%v) = %v, want %v", tt.rate, got, tt.want)
		}
	}
}

func TestCurrentPlaybackRate(t *testing.T) {
	rate = outputRate
	inputRate = 32768
	videoSyncRate = 60.0
	srcRatioOrig = 0
	srcRatioCurr = 0
	updateInputRate()
	queueBytes = 5 * bufSize
	nominalStep := 1.0 / srcRatioOrig

	t.Run("Low write availability speeds playback up", func(t *testing.T) {
		got := currentResampleStep(0)
		if got <= nominalStep {
			t.Errorf("currentResampleStep() = %v, want > %v", got, nominalStep)
		}
	})

	t.Run("High write availability slows playback down", func(t *testing.T) {
		got := currentResampleStep(queueBytes)
		if got >= nominalStep {
			t.Errorf("currentResampleStep() = %v, want < %v", got, nominalStep)
		}
	})

	t.Run("Half-full queue stays at nominal rate", func(t *testing.T) {
		got := currentResampleStep(queueBytes / 2)
		if math.Abs(got-nominalStep) > 1e-9 {
			t.Errorf("currentResampleStep() = %v, want %v", got, nominalStep)
		}
	})
}

func TestAdjustedInputRate(t *testing.T) {
	got := adjustedInputRate(32768, 59.7275, 60.0)
	want := 32768 * 60.0 / 59.7275
	if math.Abs(got-want) > 1e-6 {
		t.Fatalf("adjustedInputRate() = %v, want %v", got, want)
	}

	got = adjustedInputRate(32768, 50.0, 60.0)
	if got != 32768 {
		t.Fatalf("adjustedInputRate() = %v, want %v outside timing skew window", got, 32768)
	}
}

func TestRequiredInputFrames(t *testing.T) {
	readPhase = 0.25
	got := requiredInputFrames(1.0)
	if got != framesPerBuffer+2 {
		t.Errorf("requiredInputFrames() = %v, want %v", got, framesPerBuffer+2)
	}
}
