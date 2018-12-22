FROM resin/raspberrypi3-golang

RUN apt-get update && apt-get install libglfw3-dev libopenal-dev xorg-dev -y --allow-unauthenticated && apt-get clean
