FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /voters-api

FROM alpine:latest AS run

WORKDIR /

COPY --from=build /voters-api /voters-api

EXPOSE 1081

CMD ["/voters-api", "-p", "1081"]