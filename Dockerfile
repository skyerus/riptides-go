# Compile stage
FROM golang:1.12.5-alpine3.9 AS build-env
ENV CGO_ENABLED 0

RUN apk add --no-cache git
RUN go get github.com/derekparker/delve/cmd/dlv

ADD . /go/src/github.com/skyerus/riptides-go
RUN go build -gcflags "all=-N -l" -o /riptides-go github.com/skyerus/riptides-go

# Final stage
FROM alpine:3.9
EXPOSE 8080 40000

RUN apk add --no-cache libc6-compat

WORKDIR /
COPY --from=build-env /riptides-go /
COPY --from=build-env /go/bin/dlv /
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/riptides-go"]