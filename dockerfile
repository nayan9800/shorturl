FROM golang:1.17.3-alpine3.14 AS builder
RUN mkdir build
ADD . /build/
WORKDIR /build/
RUN go build -o app


FROM alpine:latest 
RUN mkdir app
WORKDIR /app/
COPY --from=builder /build/. .
EXPOSE 8080
ENTRYPOINT [ "./app" ]

