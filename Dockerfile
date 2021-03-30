FROM golang:alpine as builder

WORKDIR /server

COPY go.mod .
COPY main.go .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server

FROM scratch

COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=builder /server/server /bin/

CMD [ "server" ]
