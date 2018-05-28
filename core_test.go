package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/kivutar/go-playthemall/libretro"
)

func Test_coreGetGameInfo(t *testing.T) {
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
			got, err := coreGetGameInfo(tt.args.filename, tt.args.blockExtract)
			if (err != nil) != tt.wantErr {
				t.Errorf("coreGetGameInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("coreGetGameInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_coreUnzipGame(t *testing.T) {
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
			got, got1, err := coreUnzipGame(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("coreUnzipGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("coreUnzipGame() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("coreUnzipGame() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
