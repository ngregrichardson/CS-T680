FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /polls-api

FROM alpine:latest AS run

WORKDIR /

COPY --from=build /polls-api /polls-api

EXPOSE 1080

CMD ["/polls-api"]