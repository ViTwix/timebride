FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/timebride ./cmd/app

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/timebride .
COPY .env .
COPY web/ web/
COPY config/ config/
COPY migrations/ migrations/

EXPOSE 3000
CMD ["./timebride"]
