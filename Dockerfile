FROM golang:1.22.1-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN  go mod download
COPY . ./
RUN go build -o server ./cmd/server/main.go

FROM golang:1.22.1-alpine AS server
COPY --from=build /app/server ./
CMD ["./server"]
