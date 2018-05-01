# go-playthemall

go-playthemall is an attempt to write a libretro frontend in go.

It is able to launch most non GL libretro cores.

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

    go-get github.com/kivutar/go-playthemall
    go-build github.com/kivutar/go-playthemall

## Running

    go-playthemall -L nestopia_libretro.so -G mario3.nes
