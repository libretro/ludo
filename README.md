# go-playthemall [![Build Status](https://travis-ci.org/libretro/go-playthemall.svg?branch=master)](https://travis-ci.org/libretro/go-playthemall) [![Build status](https://ci.appveyor.com/api/projects/status/dndcl5m1pepnhbdk?svg=true)](https://ci.appveyor.com/project/libretro/go-playthemall-tfah4)

go-playthemall is an attempt to write a libretro frontend written in go.

<img src="assets/illustration.png" height="373" />

It is able to launch most non GL libretro cores.

It is tested on OSX, Linux and Windows.

## Dependencies

 * GLFW 3.2
 * OpenGL 4.1
 * OpenAL

On OSX you can execute the following command and follow the instructions about exporting PKG_CONFIG

    brew install glfw openal-soft

On Debian or Ubuntu:

    sudo apt-get install libglfw3-dev libopenal-dev xorg-dev

On Windows, setup openal headers and dll in mingw-w64 `include` and `lib` folders

## Building

    go get github.com/libretro/go-playthemall
    cd $GOPATH/src/github.com/libretro/go-playthemall
    go build

## Running

    ./go-playthemall -L nestopia_libretro.so mario3.nes
