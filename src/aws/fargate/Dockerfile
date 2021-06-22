## STAGE 1
FROM golang:1.15 as builder

RUN apt update && apt install -y --no-install-recommends ca-certificates

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY . .
RUN go mod download

RUN go build

## STAGE 2
FROM scratch
COPY --from=builder /app/nlpipe /app/
COPY --from=builder /app/html /app/html
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8001

WORKDIR /app
ENTRYPOINT ["./nlpipe"]