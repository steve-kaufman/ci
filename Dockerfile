FROM golang:alpine as builder

WORKDIR /build

COPY go.mod .
COPY *.go .

RUN go get -d .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ci


FROM scratch

COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=builder /build/ci /bin/

CMD [ "ci" ]
