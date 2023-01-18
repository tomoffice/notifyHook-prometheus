FROM golang:alpine as builder
WORKDIR /app
COPY ./notifyHook .
RUN go mod download
RUN go build -o exec

FROM alpine
WORKDIR /usr/bin
COPY --from=builder /app/exec .
RUN apk update && apk add tzdata
