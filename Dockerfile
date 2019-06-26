# Compile stage
FROM golang:1.12.5-alpine3.9 AS build-env
ENV CGO_ENABLED 0

RUN apk add --no-cache git
RUN apk update; \
    apk add --no-cache \
    curl \
    curl -L -s https://github.com/golang/dep/releases/download/v0.5.3/dep-linux-amd64 -o /bin/dep; \
        chmod +x /bin/dep; \
        rm -rf /var/cache/apk/*; \
        rm -rf /tmp/*;

ADD . /go/src/github.com/skyerus/riptides-go
WORKDIR /go/src/github.com/skyerus/riptides-go
RUN dep ensure
WORKDIR /

RUN go build -gcflags "all=-N -l" -o /riptides-go github.com/skyerus/riptides-go

# Final stage
FROM alpine:3.9

RUN apk add --no-cache libc6-compat

WORKDIR /
COPY --from=build-env /riptides-go /
CMD ["/riptides-go"]