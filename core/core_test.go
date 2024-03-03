package core

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"github.com/libretro/ludo/video"
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
	state.Verbose = true

	ext := utils.CoreExt()

	Init(&video.Video{})

	out := utils.CaptureOutput(func() { Load("testdata/vecx_libretro" + ext) })

	t.Run("The core is loaded", func(t *testing.T) {
		if state.Core == nil {
			t.Errorf("got = %v, want not nil", state.Core)
		}
	})

	t.Run("Logs information about the loaded core", func(t *testing.T) {
		got := out
		want := `[Core]: Name: VecX
[Core]: Version: 1.2 42366f8
[Core]: Valid extensions: bin|vec
[Core]: Need fullpath: false
[Core]: Block extract: false
`
		if !strings.Contains(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	Unload()
	state.GamePath = ""
	state.Verbose = false
}

func Test_getGameInfo(t *testing.T) {
	type args struct {
		filename     string
		blockExtract bool
	}
	tests := []struct {
		name    string
		args    args
		want    *libretro.GameInfo
		wantErr bool
	}{
		{
			name: "Returns the right path and size for an unzipped ROM",
			args: args{filename: "testdata/Polar Rescue (USA).vec", blockExtract: false},
			want: &libretro.GameInfo{
				Path: "testdata/Polar Rescue (USA).vec",
				Size: 8192,
			},
			wantErr: false,
		},
		{
			name: "Returns the right path and size for a zipped ROM",
			args: args{filename: "testdata/Polar Rescue (USA).zip", blockExtract: false},
			want: &libretro.GameInfo{
				Path: filepath.Join(os.TempDir(), "Polar Rescue (USA).vec"),
				Size: 8192,
			},
			wantErr: false,
		},
		{
			name: "Returns the right path and size for a zipped ROM with blockExtract",
			args: args{filename: "testdata/Polar Rescue (USA).zip", blockExtract: true},
			want: &libretro.GameInfo{
				Path: "testdata/Polar Rescue (USA).zip",
				Size: 6829,
			},
			wantErr: false,
		},
		{
			name: "Returns the right path and size for a zipped ROM with blockExtract",
			args: args{filename: "testdata/Polar Rescue (USA).zip", blockExtract: true},
			want: &libretro.GameInfo{
				Path: "testdata/Polar Rescue (USA).zip",
				Size: 6829,
			},
			wantErr: false,
		},
		{
			name:    "Returns an error when a file doesn't exists",
			args:    args{filename: "testdata/Polar Rescue (USA)2.zip", blockExtract: true},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Doesn't attempt to unzip a file that has no .zip extension",
			args: args{filename: "testdata/Polar Rescue (USA).vec", blockExtract: true},
			want: &libretro.GameInfo{
				Path: "testdata/Polar Rescue (USA).vec",
				Size: 8192,
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

func Test_unarchiveGame(t *testing.T) {
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
			args:    args{filename: "testdata/Polar Rescue (USA).zip"},
			want:    filepath.Join(os.TempDir(), "Polar Rescue (USA).vec"),
			want1:   8192,
			wantErr: false,
		},
		{
			name:    "Returns an error if the file is not a zip",
			args:    args{filename: "testdata/Polar Rescue (USA).vec"},
			want:    "",
			want1:   0,
			wantErr: true,
		},
		{
			name:    "Returns an error if the file doesn't exists",
			args:    args{filename: "testdata/Polar Rescue (USA)2.zip"},
			want:    "",
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := unarchiveGame(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("unarchiveGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unarchiveGame() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unarchiveGame() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_coreLoadGame(t *testing.T) {
	state.Verbose = true

	ext := utils.CoreExt()

	Init(&video.Video{})

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw")
	}

	Load("testdata/vecx_libretro" + ext)

	got := utils.CaptureOutput(func() { LoadGame("testdata/Polar Rescue (USA).vec") })

	t.Run("Logs information about the loaded game", func(t *testing.T) {
		want := `[Core]: Game loaded: testdata/Polar Rescue (USA).vec`
		if !strings.Contains(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Global state should be set by Load", func(t *testing.T) {
		if state.Core == nil {
			t.Errorf("got = %v, want %v", nil, state.Core)
		}
		if state.GamePath != "testdata/Polar Rescue (USA).vec" {
			t.Errorf("got = %v, want %v", state.GamePath, "testdata/Polar Rescue (USA).vec")
		}
		if !state.CoreRunning {
			t.Errorf("got = %v, want %v", state.CoreRunning, true)
		}
	})

	UnloadGame()
	Unload()

	t.Run("Global state should be cleared by Unload", func(t *testing.T) {
		if state.Core != nil {
			t.Errorf("got = %v, want %v", state.Core, nil)
		}
		if state.GamePath != "" {
			t.Errorf("got = %v, want %v", state.GamePath, "")
		}
		if state.CoreRunning {
			t.Errorf("got = %v, want %v", state.CoreRunning, false)
		}
	})
}
