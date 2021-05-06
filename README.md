# ludo ![Build Status](https://github.com/libretro/ludo/workflows/CI/badge.svg) [![GoDoc](https://godoc.org/github.com/libretro/ludo?status.svg)](https://godoc.org/github.com/libretro/ludo)

Ludo is a work in progress libretro frontend written in go.

<img src="https://raw.githubusercontent.com/kivutar/ludo-assets/master/illustration.png" />

It is able to launch most non GL libretro cores.

It works on OSX, Linux, Linux ARM and Windows. You can download releases [here](https://github.com/libretro/ludo/releases)

## Dependencies

- GLFW 3.3
- OpenGL >= 2.1
- OpenAL

#### On OSX

You can execute the following command and follow the instructions about exporting PKG_CONFIG

    brew install openal-soft

#### On Debian or Ubuntu

    sudo apt-get install libopenal-dev xorg-dev golang

#### On Raspbian

You need to enable the experimental VC4 OpenGL support (Full KMS) in raspi-config.

    sudo apt-get install libopenal-dev xorg-dev

#### On Alpine / postmarketOS

    sudo apk add musl-dev gcc openal-soft-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev mesa-dev

#### On Windows

Setup openal headers and dll in mingw-w64 `include` and `lib` folders.

## Building

    git clone --recursive https://github.com/libretro/ludo.git
    cd ludo
    go build

For more detailed build steps, please refer to [our continuous delivery config](https://github.com/libretro/ludo/blob/master/.github/workflows/cd.yml).

## Running

    ./ludo
