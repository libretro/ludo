FROM arm32v7/golang:1.11-stretch

RUN apt-get update && apt-get install libglfw3-dev libopenal-dev xorg-dev -y --allow-unauthenticated && apt-get clean
