# Dockerfile used to build the application

# Build container containing our pre-pulled libraries
FROM golang:latest as build

# Static compile
ENV CGO_ENABLED=0
ENV GOOS=linux

# Ensure we have the libraries - docker will cache these between builds
RUN go get -v \
      log \
      github.com/go-stomp/stomp \
      github.com/streadway/amqp

WORKDIR /src

ADD src /src/

RUN go build \
      -v \
      -x \
      -o /usr/local/bin/bridge \
      .

# Runtime container containing just the single binary
FROM scratch
COPY --from=build /usr/local/bin/bridge /usr/local/bin/bridge
CMD ["bridge"]
