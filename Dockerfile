FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/6Matt/se390-internal

ENV USER matt
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET XuzU1jUxpthIg6l5

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://matt@localhost:5432/se390-internal?sslmode=disable

WORKDIR /go/src/github.com/6Matt/se390-internal

RUN godep go build

EXPOSE 8888
CMD ./se390-internal