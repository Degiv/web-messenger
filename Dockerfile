FROM golang:1.22.2 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 go build -installsuffix 'static' -o messenger ./cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/messenger .

EXPOSE 3000

CMD ["./messenger"]