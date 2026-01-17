FROM golang:1.25.5-alpine AS builder

WORKDIR /app
COPY go.mod go.sum

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /bin/service ./cmd/app/main.go

FROM alpine:latest

RUN apk --no-cache add tzdata

WORKDIR /app

COPY --from=builder /bin/service .

ENTRYPOINT [ "./service" ]


