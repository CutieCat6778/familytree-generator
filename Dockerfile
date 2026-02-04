FROM golang:1.25.6-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/familytree-server ./cmd/server

FROM alpine:3.19

WORKDIR /app
COPY --from=build /out/familytree-server ./familytree-server
COPY --from=build /src/data /data

EXPOSE 3005 
ENTRYPOINT ["./familytree-server"]
CMD ["-port", "3005", "-data", "/data"]
