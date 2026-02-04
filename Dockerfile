FROM golang:1.25.6-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/familytree-server ./cmd/server

FROM alpine:3.19

WORKDIR /app
COPY --from=build /out/familytree-server ./familytree-server

EXPOSE 8080
ENTRYPOINT ["./familytree-server"]
