# Compile stage
FROM golang:1.12.5-alpine3.9 AS build-env
ENV CGO_ENABLED 0

RUN apk add --no-cache git
RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/skyerus/riptides-go
WORKDIR /go/src/github.com/skyerus/riptides-go
RUN dep ensure
WORKDIR /

RUN go build -gcflags "all=-N -l" -o /riptides-go github.com/skyerus/riptides-go

# Final stage
FROM golang:1.12.5-alpine3.9

ENV CGO_ENABLED 0
RUN apk add --no-cache libc6-compat
RUN apk add --no-cache bash

WORKDIR /
COPY --from=build-env /riptides-go /
COPY --from=build-env /go/src/github.com/skyerus/riptides-go/assets /go/src/github.com/skyerus/riptides-go/assets
CMD ["/riptides-go"]