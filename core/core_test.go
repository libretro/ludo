package core

import (
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/libretro/go-playthemall/libretro"
	"github.com/libretro/go-playthemall/options"
	"github.com/libretro/go-playthemall/state"
	"github.com/libretro/go-playthemall/utils"
	"github.com/libretro/go-playthemall/video"
)

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

func Test_coreLoad(t *testing.T) {
	state.Global.Verbose = true

	exts := map[string]string{
		"darwin":  ".dylib",
		"linux":   ".so",
		"windows": ".dll",
	}

	ext := exts[runtime.GOOS]

	Init(&video.Video{Window: &WindowMock{}}, &options.Options{})

	out := utils.CaptureOutput(func() { Load("testdata/uzem_libretro" + ext) })

	t.Run("The core is loaded", func(t *testing.T) {
		if state.Global.Core == (libretro.Core{}) {
			t.Errorf("got = %v, want not libretro.Core{}", state.Global.Core)
		}
	})

	t.Run("Logs information about the loaded core", func(t *testing.T) {
		got := out
		want := `[Core]: Name: Uzem
[Core]: Version: v2.0
[Core]: Valid extensions: uze
[Core]: Need fullpath: false
[Core]: Block extract: false
[Core]: Core loaded: Uzem
`
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	state.Global.Core.UnloadGame()
	state.Global.Core.Deinit()
	state.Global.GamePath = ""
	state.Global.Verbose = false
}

func Test_getGameInfo(t *testing.T) {
	type args struct {
		filename     string
		blockExtract bool
	}
	tests := []struct {
		name    string
		args    args
		want    libretro.GameInfo
		wantErr bool
	}{
		{
			name: "Returns the right path and size for an unzipped ROM",
			args: args{filename: "testdata/ZoomingSecretary.uze", blockExtract: false},
			want: libretro.GameInfo{
				Path: "testdata/ZoomingSecretary.uze",
				Size: 61286,
			},
			wantErr: false,
		},
		{
			name: "Returns the right path and size for a zipped ROM",
			args: args{filename: "testdata/ZoomingSecretary.zip", blockExtract: false},
			want: libretro.GameInfo{
				Path: os.TempDir() + "/ZoomingSecretary.uze",
				Size: 61286,
			},
			wantErr: false,
		},
		{
			name: "Returns the right path and size for a zipped ROM with blockExtract",
			args: args{filename: "testdata/ZoomingSecretary.zip", blockExtract: true},
			want: libretro.GameInfo{
				Path: "testdata/ZoomingSecretary.zip",
				Size: 25599,
			},
			wantErr: false,
		},
		{
			name: "Returns the right path and size for a zipped ROM with blockExtract",
			args: args{filename: "testdata/ZoomingSecretary.zip", blockExtract: true},
			want: libretro.GameInfo{
				Path: "testdata/ZoomingSecretary.zip",
				Size: 25599,
			},
			wantErr: false,
		},
		{
			name:    "Returns an error when a file doesn't exists",
			args:    args{filename: "testdata/ZoomingSecretary2.zip", blockExtract: true},
			want:    libretro.GameInfo{},
			wantErr: true,
		},
		{
			name: "Doesn't attempt to unzip a file that has no .zip extension",
			args: args{filename: "testdata/ZoomingSecretary.uze", blockExtract: true},
			want: libretro.GameInfo{
				Path: "testdata/ZoomingSecretary.uze",
				Size: 61286,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getGameInfo(tt.args.filename, tt.args.blockExtract)
			if (err != nil) != tt.wantErr {
				t.Errorf("getGameInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getGameInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unzipGame(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int64
		wantErr bool
	}{
		{
			name:    "Should unzip to the right path",
			args:    args{filename: "testdata/ZoomingSecretary.zip"},
			want:    os.TempDir() + "/ZoomingSecretary.uze",
			want1:   61286,
			wantErr: false,
		},
		{
			name:    "Returns an error if the file is not a zip",
			args:    args{filename: "testdata/ZoomingSecretary.uze"},
			want:    "",
			want1:   0,
			wantErr: true,
		},
		{
			name:    "Returns an error if the file doesn't exists",
			args:    args{filename: "testdata/ZoomingSecretary2.zip"},
			want:    "",
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := unzipGame(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("unzipGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unzipGame() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unzipGame() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
