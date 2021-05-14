package video

import (
	"testing"
)

func TestXYWHTo4points(t *testing.T) {
	type args struct {
		x   float32
		y   float32
		w   float32
		h   float32
		fbh float32
	}
	tests := []struct {
		name   string
		args   args
		wantX1 float32
		wantY1 float32
		wantX2 float32
		wantY2 float32
		wantX3 float32
		wantY3 float32
		wantX4 float32
		wantY4 float32
	}{
		{
			name:   "Works",
			args:   args{x: 30, y: 40, w: 500, h: 600, fbh: 800},
			wantX1: 30,
			wantX2: 30,
			wantX3: 530,
			wantX4: 530,
			wantY1: 160,
			wantY2: 760,
			wantY3: 160,
			wantY4: 760,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX1, gotY1, gotX2, gotY2, gotX3, gotY3, gotX4, gotY4 := XYWHTo4points(tt.args.x, tt.args.y, tt.args.w, tt.args.h, tt.args.fbh)
			if gotX1 != tt.wantX1 {
				t.Errorf("XYWHTo4points() gotX1 = %v, want %v", gotX1, tt.wantX1)
			}
			if gotY1 != tt.wantY1 {
				t.Errorf("XYWHTo4points() gotY1 = %v, want %v", gotY1, tt.wantY1)
			}
			if gotX2 != tt.wantX2 {
				t.Errorf("XYWHTo4points() gotX2 = %v, want %v", gotX2, tt.wantX2)
			}
			if gotY2 != tt.wantY2 {
				t.Errorf("XYWHTo4points() gotY2 = %v, want %v", gotY2, tt.wantY2)
			}
			if gotX3 != tt.wantX3 {
				t.Errorf("XYWHTo4points() gotX3 = %v, want %v", gotX3, tt.wantX3)
			}
			if gotY3 != tt.wantY3 {
				t.Errorf("XYWHTo4points() gotY3 = %v, want %v", gotY3, tt.wantY3)
			}
			if gotX4 != tt.wantX4 {
				t.Errorf("XYWHTo4points() gotX4 = %v, want %v", gotX4, tt.wantX4)
			}
			if gotY4 != tt.wantY4 {
				t.Errorf("XYWHTo4points() gotY4 = %v, want %v", gotY4, tt.wantY4)
			}
		})
	}
}
