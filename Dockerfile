FROM alpine AS voices
WORKDIR /voices
RUN apk add --no-cache curl \
        && \
    curl -sL https://huggingface.co/rhasspy/piper-voices/resolve/main/voices.json?download=true -o voices.json

FROM golang:alpine AS builder
WORKDIR /app
COPY go.* .
RUN go mod download
COPY *.go .
RUN go build -o gin .

FROM alpine as piper
ARG TARGETARCH
ARG TARGETVARIANT
ARG BINARY_PIPER_VERSION='1.2.0'
RUN apk add --no-cache curl tar gzip \
    && mkdir /piper \
    && curl -L -s \
        "https://github.com/rhasspy/piper/releases/download/v${BINARY_PIPER_VERSION}/piper_${TARGETARCH}${TARGETVARIANT}.tar.gz" \
        | tar -zxvf - -C /piper


FROM debian:latest
ENV VOICES_PATH="/voices"
ENV VOICES_JSON_PATH="/app/voices.json"

WORKDIR /app
VOLUME ["$VOICES_PATH"]

ENV GIN_MODE=release

RUN apt update && apt install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=piper /piper /usr/share
COPY --from=voices /voices/voices.json $VOICES_JSON_PATH
COPY --from=builder /app/gin .

CMD ["/app/gin"]