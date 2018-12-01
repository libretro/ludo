# ludo [![Build Status](https://travis-ci.org/libretro/ludo.svg?branch=master)](https://travis-ci.org/libretro/ludo) [![GoDoc](https://godoc.org/github.com/libretro/ludo?status.svg)](https://godoc.org/github.com/libretro/ludo)

Ludo is a work in progress libretro frontend written in go.

<img src="https://raw.githubusercontent.com/kivutar/ludo-assets/master/illustration.png" />

It is able to launch most non GL libretro cores.

It works on OSX, Linux and Windows. You can download releases [here](https://github.com/libretro/ludo/releases)

## Dependencies

 * GLFW 3.2
 * OpenGL 4.1
 * OpenAL

#### On OSX

You can execute the following command and follow the instructions about exporting PKG_CONFIG

    brew install glfw openal-soft

#### On Debian or Ubuntu

    sudo apt-get install libglfw3-dev libopenal-dev xorg-dev

#### On Raspbian

You need to enable the experimental VC4 OpenGL support (Full KMS) in raspi-config.

    sudo apt-get install libglfw3-dev libopenal-dev xorg-dev

#### On Windows

Setup openal headers and dll in mingw-w64 `include` and `lib` folders.

## Building

    go get github.com/libretro/ludo
    cd $GOPATH/src/github.com/libretro/ludo
    go build

## Running

    ./ludo

If on a RaspberryPi:

    ./ludo -glver=21
