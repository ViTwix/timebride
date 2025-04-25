FROM golang:1.23-alpine AS builder

# Встановлюємо необхідні залежності для збірки
RUN apk add --no-cache gcc musl-dev git make

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
COPY internal/ internal/

# Створюємо директорії для статичних файлів
RUN mkdir -p /app/web/public/css /app/web/public/js /app/web/public/img /app/web/static/uploads /app/web/static/temp

# Встановлюємо правильні права доступу
RUN chmod -R 755 /app/web
RUN chmod +x /app/timebride

# Встановлюємо змінні середовища для шляхів
ENV TEMPLATE_DIR=./web/templates
ENV STATIC_DIR=./web/static
ENV PUBLIC_DIR=./web/public
ENV CONTROLLERS_DIR=./web/controllers

# Додаємо логування для діагностики
RUN ls -la /app
RUN ls -la /app/timebride || echo "Executable not found"
RUN file /app/timebride || echo "Cannot determine file type"

EXPOSE 3000
CMD ["/bin/sh", "-c", "/app/timebride"]
