FROM golang:alpine3.10 as builder
RUN mkdir /build /deps
RUN apk update && apk add git
WORKDIR /deps
ADD go.mod /deps
RUN go mod download
ADD / /build
WORKDIR /build
RUN go build -o snowedin .
FROM alpine:3.10.3
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/snowedin /app/
COPY --from=builder /build/tests/config.yaml /app/
WORKDIR /app
ENTRYPOINT ["./snowedin"]
