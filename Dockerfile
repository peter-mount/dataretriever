# Dockerfile used to build the application

# ============================================================
# The base golang environment with curl, git, uptodate tzdata
# and go-bindata installed
FROM golang:alpine as golang
RUN apk add --no-cache \
      curl \
      git \
      tzdata

# Build container containing our pre-pulled libraries
FROM golang as source

WORKDIR /work

# Download dependencies before copying any sources then we
# can use the docker cache to limit updates
ADD go.mod .
RUN go mod download

# Ensure we have the libraries - docker will cache these between builds
#RUN go get -v \
#      flag \
#      github.com/go-stomp/stomp \
#      github.com/peter-mount/golib/rabbitmq \
#      github.com/peter-mount/golib/statistics  \
#      github.com/streadway/amqp \
#      gopkg.in/robfig/cron.v2 \
#      gopkg.in/yaml.v2 \
#      io/ioutil \
#      log \
#      net/http \
#      path/filepath \
#      time

ADD src/ .

# ============================================================
# Now compile our binaries
FROM source as compiler
ARG arch
ARG goos
ARG goarch
ARG goarm

RUN CGO_ENABLED=0 \
    GOOS=${goos} \
    GOARCH=${goarch} \
    GOARM=${goarm} \
    go build -o /usr/local/bin/dataretriever \
      .

# ============================================================
# Finally build the final runtime container
FROM alpine

RUN apk add --no-cache \
      curl \
      tzdata

COPY --from=compiler /usr/local/bin/dataretriever /usr/local/bin/dataretriever

CMD ["dataretriever"]
