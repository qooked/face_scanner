FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go install -v github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=build /app/main .
COPY --from=build /app/config ./config
COPY --from=build /app/files ./files
COPY --from=build /go/bin/migrate /usr/local/bin/migrate



EXPOSE 8080
CMD ["./main"]
