# go-nanoarch

go-nanoarch is an attempt to write a minimal libretro frontend in go. It is a port of https://github.com/heuripedes/nanoarch

As the C version:

> nanoarch is a small libretro frontend (nanoarch.c has less than 1000 lines of code) created for educational purposes. It only provides the required (video, audio and basic input) features to run most non-libretro-gl cores and there's no UI or configuration support.

It is tested on OSX and Linux.

## Dependencies

 * GLFW 3.2
 * OpenGL 4.1
 * OpenAL

On OSX you can execute the following command and follow the instructions about exporting PKG_CONFIG

    brew install glfw openal-soft

On Debian or Ubuntu:

    sudo apt-get install libglfw3-dev libopenal-dev xorg-dev

## Building

    go-get github.com/kivutar/go-nanoarch
    go-build github.com/kivutar/go-nanoarch

## Running

    go-nanoarch -L nestopia_libretro.so -G mario3.nes
