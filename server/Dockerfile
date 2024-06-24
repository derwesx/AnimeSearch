FROM golang:1.22-bullseye as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o app ./cmd/AnimeSearch/main.go

FROM ubuntu:22.04
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app .

EXPOSE 8080

CMD ["./app"]