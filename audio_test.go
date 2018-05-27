package main

import "testing"

func Test_fillInternalBuf(t *testing.T) {
	type args struct {
		buf  []byte
		size int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fillInternalBuf(tt.args.buf, tt.args.size); got != tt.want {
				t.Errorf("fillInternalBuf() = %v, want %v", got, tt.want)
			}
		})
	}
}
