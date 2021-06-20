FROM golang:1.15

#ENV GO111MODULE=on \
#    CGO_ENABLED=0 \
#    GOOS=linux \
#    GOARCH=amd64

WORKDIR /build
COPY . .
#RUN go mod download

RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .

EXPOSE 8080

CMD ["/dist/main"]