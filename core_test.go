package main

import (
	"os/user"
	"reflect"
	"testing"

	"github.com/kivutar/go-playthemall/libretro"
)

func Test_coreGetGameInfo(t *testing.T) {
	usr, _ := user.Current()

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
				Path: usr.HomeDir + "/ZoomingSecretary.uze",
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
			name: "Doesn't attemtp to unzip a file that has no .zip extension",
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
