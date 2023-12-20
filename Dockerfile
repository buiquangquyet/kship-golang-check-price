# Builder
FROM docker.citigo.com.vn/kship/templates/golang-21:builder AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=1 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /app

RUN go build -v -o /app/main .

# Runtime
FROM docker.citigo.com.vn/templates/golang:runtime

WORKDIR /app

COPY --from=builder /app/main /app/main

ENTRYPOINT ["/app/main"]
