package video

import (
	"reflect"
	"testing"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/kivutar/glfont"
	"github.com/libretro/ludo/libretro"
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

type WindowMock struct{}

func (m WindowMock) GetFramebufferSize() (width, height int)     { return 320, 240 }
func (m WindowMock) Destroy()                                    {}
func (m WindowMock) MakeContextCurrent()                         {}
func (m WindowMock) SetSizeLimits(minw, minh, maxw, maxh int)    {}
func (m WindowMock) SetInputMode(mode glfw.InputMode, value int) {}
func (m WindowMock) GetKey(key glfw.Key) glfw.Action             { return 0 }
func (m WindowMock) SetShouldClose(bool)                         {}
func (m WindowMock) ShouldClose() bool                           { return false }
func (m WindowMock) SetTitle(string)                             {}
func (m WindowMock) SwapBuffers()                                {}

func TestVideo_vertexArray(t *testing.T) {

	var myWindowMock WindowMock

	type fields struct {
		Window         WindowInterface
		Geom           libretro.GameGeometry
		Font           *glfont.Font
		program        uint32
		roundedProgram uint32
		borderProgram  uint32
		circleProgram  uint32
		demulProgram   uint32
		vao            uint32
		vbo            uint32
		texID          uint32
		white          uint32
		pitch          int32
		pixFmt         uint32
		pixType        uint32
		bpp            int32
	}
	type args struct {
		x     float32
		y     float32
		w     float32
		h     float32
		scale float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []float32
	}{
		{
			name: "Works",
			fields: fields{
				Window: myWindowMock,
			},
			args: args{
				x:     10,
				y:     11,
				w:     300,
				h:     400,
				scale: 2,
			},
			want: []float32{
				-0.9375, -5.758333, 0, 1,
				-0.9375, 0.9083333, 0, 0,
				2.8125, -5.758333, 1, 1,
				2.8125, 0.9083333, 1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			video := &Video{
				Window:         tt.fields.Window,
				Geom:           tt.fields.Geom,
				Font:           tt.fields.Font,
				program:        tt.fields.program,
				roundedProgram: tt.fields.roundedProgram,
				borderProgram:  tt.fields.borderProgram,
				circleProgram:  tt.fields.circleProgram,
				demulProgram:   tt.fields.demulProgram,
				vao:            tt.fields.vao,
				vbo:            tt.fields.vbo,
				texID:          tt.fields.texID,
				white:          tt.fields.white,
				pitch:          tt.fields.pitch,
				pixFmt:         tt.fields.pixFmt,
				pixType:        tt.fields.pixType,
				bpp:            tt.fields.bpp,
			}
			if got := video.vertexArray(tt.args.x, tt.args.y, tt.args.w, tt.args.h, tt.args.scale); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Video.vertexArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
