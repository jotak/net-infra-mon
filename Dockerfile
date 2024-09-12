ARG TARGETARCH
FROM docker.io/library/golang:1.22 as go-builder

ARG TARGETARCH=amd64
ARG GO_BUILD_OPTS

WORKDIR /opt/app-root

COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor/ vendor/
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build $GO_BUILD_OPTS -mod vendor -o nim-monitoring-agent cmd/agent/main.go

FROM --platform=linux/$TARGETARCH registry.access.redhat.com/ubi9/ubi-minimal:9.4

COPY --from=go-builder /opt/app-root/nim-monitoring-agent ./

ENTRYPOINT ["./nim-monitoring-agent"]
