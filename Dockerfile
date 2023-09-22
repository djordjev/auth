FROM golang:1.20-alpine as builder

LABEL maintainer="Gyorgy Technologies"

WORKDIR "/src"

ADD . .

RUN ["go", "mod", "tidy"]

RUN ["go", "build", "-o", "./build/auth", "./cmd/auth/main.go"]

FROM alpine:latest

WORKDIR /src

COPY --from=builder /src/build/auth auth

RUN chmod -R 777 auth

CMD ./auth

EXPOSE 13010
