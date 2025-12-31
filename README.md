# Demo Service - Production-Ready Go CRUD Application

Production-ready Go приложение с CRUD операциями для управления продуктами, включающее JWT аутентификацию, Prometheus метрики, Swagger документацию, rate limiting и структурированное логирование.

## Возможности

- ✅ REST API для CRUD операций с продуктами
- ✅ JWT аутентификация
- ✅ Prometheus метрики
- ✅ Swagger/OpenAPI документация
- ✅ Rate limiting
- ✅ Структурированное логирование
- ✅ Embedded SQLite база данных
- ✅ Graceful shutdown
- ✅ Health check endpoints

## Требования

- Go 1.21 или выше
- Linux/Windows/macOS

## Установка и запуск

### 1. Клонирование и установка зависимостей

```bash
go mod download
```

### 2. Настройка конфигурации

Скопируйте `.env.example` в `.env` и настройте параметры:

```bash
cp .env.example .env
```

Отредактируйте `.env` файл, особенно `JWT_SECRET` для продакшена.

### 3. Запуск приложения

```bash
go run cmd/server/main.go
```

Или скомпилируйте и запустите:

```bash
go build -o demo-service cmd/server/main.go
./demo-service
```

Приложение будет доступно на `http://localhost:8080` (или порт, указанный в `SERVER_PORT`).

## API Endpoints

### Аутентификация

- `POST /api/v1/auth/register` - Регистрация нового пользователя
- `POST /api/v1/auth/login` - Вход в систему (возвращает JWT токен)

### Продукты (требуют JWT токен)

- `POST /api/v1/products` - Создание продукта
- `GET /api/v1/products` - Список продуктов (с пагинацией)
- `GET /api/v1/products/:id` - Получение продукта по ID
- `PUT /api/v1/products/:id` - Обновление продукта
- `DELETE /api/v1/products/:id` - Удаление продукта

### Системные

- `GET /health` - Health check
- `GET /metrics` - Prometheus метрики
- `GET /swagger/index.html` - Swagger UI

## Использование API

### 1. Регистрация пользователя

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }'
```

### 2. Вход и получение токена

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }'
```

Ответ содержит `token` - используйте его в заголовке `Authorization: Bearer <token>`.

### 3. Создание продукта

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "name": "Test Product",
    "description": "Product description",
    "price": 99.99,
    "stock": 100
  }'
```

### 4. Получение списка продуктов

```bash
curl -X GET "http://localhost:8080/api/v1/products?page=1&limit=10" \
  -H "Authorization: Bearer <your-token>"
```

## Деплой на AWS EC2

### 1. Подготовка EC2 инстанса

```bash
# Обновление системы (Ubuntu/Debian)
sudo apt update && sudo apt upgrade -y

# Установка Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Создание директории для приложения
mkdir -p /opt/demo-service
```

### 2. Копирование файлов на сервер

```bash
# С вашего локального компьютера
scp -r . user@your-ec2-ip:/opt/demo-service/
```

### 3. Компиляция на сервере

```bash
cd /opt/demo-service
go mod download
go build -o demo-service cmd/server/main.go
```

### 4. Создание systemd service

Создайте файл `/etc/systemd/system/demo-service.service`:

```ini
[Unit]
Description=Demo Service
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/opt/demo-service
EnvironmentFile=/opt/demo-service/.env
ExecStart=/opt/demo-service/demo-service
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### 5. Запуск сервиса

```bash
sudo systemctl daemon-reload
sudo systemctl enable demo-service
sudo systemctl start demo-service
sudo systemctl status demo-service
```

### 6. Настройка Nginx (опционально)

Создайте `/etc/nginx/sites-available/demo-service`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Активируйте конфигурацию:

```bash
sudo ln -s /etc/nginx/sites-available/demo-service /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 7. Настройка Security Group

Убедитесь, что в Security Group EC2 открыты порты:
- 22 (SSH)
- 80 (HTTP, если используете Nginx)
- 443 (HTTPS, если используете SSL)
- 8080 (если обращаетесь напрямую к приложению)

## Мониторинг

### Prometheus метрики

Метрики доступны на `/metrics` endpoint. Вы можете настроить Prometheus для сбора метрик:

```yaml
scrape_configs:
  - job_name: 'demo-service'
    static_configs:
      - targets: ['localhost:8080']
```

### Логи

Логи приложения можно просмотреть через systemd:

```bash
sudo journalctl -u demo-service -f
```

Или если запускаете напрямую, логи выводятся в stdout/stderr.

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `SERVER_PORT` | Порт HTTP сервера | 8080 |
| `DB_PATH` | Путь к SQLite файлу | ./data/demo.db |
| `JWT_SECRET` | Секретный ключ для JWT | (обязательно) |
| `JWT_EXPIRY` | Время жизни JWT токена | 24h |
| `RATE_LIMIT_RPS` | Запросов в секунду | 10 |
| `LOG_LEVEL` | Уровень логирования (debug, info, warn, error) | info |

## Структура проекта

```
demo-service/
├── cmd/server/main.go          # Точка входа
├── internal/
│   ├── config/                 # Конфигурация
│   ├── handler/                # HTTP handlers
│   ├── service/                # Бизнес-логика
│   ├── repository/             # Работа с БД
│   ├── model/                  # Модели данных
│   ├── middleware/             # Middleware
│   ├── metrics/                # Prometheus метрики
│   └── database/               # Инициализация БД
├── pkg/jwt/                    # JWT утилиты
├── migrations/                 # SQL миграции
└── docs/                       # Swagger документация
```

## Разработка

### Генерация Swagger документации

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go
```

### Запуск тестов

```bash
go test ./...
```

## Лицензия

MIT

