FROM golang:1.16 as dependencies
WORKDIR /src
COPY go.mod .
COPY go.sum .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go mod download
RUN go get -u github.com/psampaz/go-mod-outdated
RUN go list -u -m -json all | go-mod-outdated -direct -update

FROM dependencies AS builder
WORKDIR /src
COPY . .
WORKDIR /src/cmd/apid
RUN go build -o /out/cmd

FROM debian:stretch-slim
ENV DEBIAN_FRONTEND noninteractive
RUN adduser --disabled-password --no-create-home --gecos '' appuser
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl net-tools ffmpeg \
    && apt-get clean -y \
    && apt-get autoremove -y \
    && rm -rf /tmp/* /var/tmp/* \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /app && chown appuser:appuser /app
RUN mkdir /var/videos && chown appuser:appuser /var/videos
WORKDIR /app
USER appuser
COPY --from=builder --chown=appuser /out/cmd .

ENV GIN_MODE=release
EXPOSE 8085
ENTRYPOINT ["/app/cmd"]