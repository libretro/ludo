package audio

import (
	"testing"
)

func Test_alUnqueueBuffers(t *testing.T) {
	t.Run("Return false if no buffers were processed", func(t *testing.T) {
		got := alUnqueueBuffers()
		if got {
			t.Errorf("alUnqueueBuffers() = %v, want %v", got, false)
		}
	})
}

func Test_Sample(t *testing.T) {
	t.Run("Doesn't crash when called", func(t *testing.T) {
		Sample(-30000, -30000)
		Sample( 30000,  30000)
	})
}

func Test_fillInternalBuf(t *testing.T) {
	Reconfigure(48000)
	type args struct {
		buf  []byte
		size int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "Fill the buffer partially",
			args: args{
				buf:  make([]byte, bufSize),
				size: 6000,
			},
			want: 6000,
		},
		{
			name: "Fill the buffer fully",
			args: args{
				buf:  make([]byte, bufSize),
				size: 6000,
			},
			want: 2192,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fillInternalBuf(tt.args.buf[:tt.args.size]); got != tt.want {
				t.Errorf("fillInternalBuf() = %v, want %v", got, tt.want)
			}
		})
	}
}
