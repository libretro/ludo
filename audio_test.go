package main

import "testing"

func Test_fillInternalBuf(t *testing.T) {
	audioInit(48000)
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
				buf:  make([]byte, 4096),
				size: 3000,
			},
			want: 3000,
		},
		{
			name: "Fill the buffer fully",
			args: args{
				buf:  make([]byte, 4096),
				size: 3000,
			},
			want: 1096,
		},
		{
			name: "Fill the buffer fully",
			args: args{
				buf:  make([]byte, 4096),
				size: 6000,
			},
			want: 6000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fillInternalBuf(tt.args.buf, tt.args.size); got != tt.want {
				t.Errorf("fillInternalBuf() = %v, want %v", got, tt.want)
			}
		})
	}
}
