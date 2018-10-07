# ludo [![Build Status](https://travis-ci.org/libretro/ludo.svg?branch=master)](https://travis-ci.org/libretro/ludo) [![Build status](https://ci.appveyor.com/api/projects/status/o3alckoqebr5t4s4?svg=true)](https://ci.appveyor.com/project/kivutar/ludo) [![GoDoc](https://godoc.org/github.com/libretro/ludo?status.svg)](https://godoc.org/github.com/libretro/ludo)

Ludo is a work in progress libretro frontend written in go.

<img src="assets/illustration.png" />

It is able to launch most non GL libretro cores.

It works on OSX, Linux and Windows.

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

    go get github.com/libretro/ludo
    cd $GOPATH/src/github.com/libretro/ludo
    go build

## Running

    ./ludo -L nestopia_libretro.so mario3.nes
