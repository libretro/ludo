FROM resin/raspberrypi3-golang

ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7

RUN apt-get update && apt-get install libglfw3-dev libopenal-dev xorg-dev -y --allow-unauthenticated && apt-get clean
