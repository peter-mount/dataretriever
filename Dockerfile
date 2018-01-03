# Dockerfile used to build the application

# Build container containing our pre-pulled libraries
FROM golang:latest as build

# Static compile
ENV CGO_ENABLED=0
ENV GOOS=linux

# Ensure we have the libraries - docker will cache these between builds
RUN go get -v \
      flag \
      log \
      io/ioutil \
      path/filepath \
      github.com/go-stomp/stomp \
      github.com/streadway/amqp \
      gopkg.in/yaml.v2

# Import the source and compile
WORKDIR /src

ADD src /src/

RUN go build \
      -v \
      -x \
      -o /usr/local/bin/bridge \
      .

# Finally build the final runtime container containing just the single static binary
FROM scratch
COPY --from=build /usr/local/bin/bridge /usr/local/bin/bridge
CMD ["bridge"]
