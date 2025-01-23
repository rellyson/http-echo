FROM golang:1.23.4-alpine AS builder

WORKDIR /tmp/build

# download dependencies
COPY go.mod go.sum ./
RUN go mod download

# build binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /tmp/build/http-echo ./cmd/server/main.go

# --------------------------------------------------
FROM scratch

WORKDIR /app

#  copy binary
COPY --from=builder /tmp/build/http-echo /app/http-echo
COPY VERSION /app/VERSION
COPY LICENSE /usr/share/doc/http-echo/LICENSE.txt

#  expose port
EXPOSE 3000/tcp

#  start http server
CMD ["/app/http-echo"]
