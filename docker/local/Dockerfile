FROM golang:1.20.5-alpine as builder

WORKDIR /app

COPY . /app

RUN go mod download && go build -o bin/web-grpc-video-chat ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin /app

CMD [ "./web-grpc-video-chat" ]
