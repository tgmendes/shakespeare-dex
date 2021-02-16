
FROM  golang:1.15.8-alpine AS build
ENV CGO_ENABLED=0
RUN mkdir -p /app

COPY go.* /app/
WORKDIR /app
RUN go mod download

COPY . .

RUN go build -o shakespeare-dex

FROM alpine:latest
COPY --from=build /app/shakespeare-dex /app/shakespeare-dex
WORKDIR /app
CMD ["./shakespeare-dex"]
