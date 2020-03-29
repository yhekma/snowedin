FROM golang:alpine3.10 as builder
RUN mkdir /build
WORKDIR /build
RUN apk update && apk add git
ADD / /build
RUN go get -v -d .
RUN go build -o snowedin .
FROM alpine:3.10.3
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/snowedin /app/
COPY --from=builder /build/tests/config.yaml /app/
WORKDIR /app
ENTRYPOINT ["./snowedin"]
