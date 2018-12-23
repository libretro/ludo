FROM dockcross/linux-armv7

RUN curl -fsSL "https://golang.org/dl/go1.11.4.linux-armv6l.tar.gz" | tar -xzC /usr/local
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7

RUN apt-get update && apt-get install libglfw3-dev libopenal-dev xorg-dev -y --allow-unauthenticated && apt-get clean
