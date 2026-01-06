# Demo Service - Production-Ready Go CRUD Application

Production-ready Go application with CRUD operations for product management, including JWT authentication, Prometheus metrics, Swagger documentation, rate limiting, and structured logging.

## Features

- ✅ REST API for CRUD operations with products
- ✅ JWT authentication
- ✅ Prometheus metrics
- ✅ Swagger/OpenAPI documentation
- ✅ Rate limiting
- ✅ Structured logging
- ✅ Embedded SQLite database
- ✅ Graceful shutdown
- ✅ Health check endpoints

## Requirements

- Go 1.21 or higher
- Linux/Windows/macOS

## Installation and Running

### 1. Clone and Install Dependencies

```bash
go mod download
```

### 2. Configure Settings

Copy `.env.example` to `.env` and configure parameters:

```bash
cp .env.example .env
```

Edit the `.env` file, especially `JWT_SECRET` for production.

### 3. Run the Application

```bash
go run cmd/server/main.go
```

Or compile and run:

```bash
go build -o demo-service cmd/server/main.go
./demo-service
```

The application will be available at `http://localhost:8080` (or the port specified in `SERVER_PORT`).

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login (returns JWT token)

### Products (require JWT token)

- `POST /api/v1/products` - Create a product
- `GET /api/v1/products` - List products (with pagination)
- `GET /api/v1/products/:id` - Get product by ID
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product

### System

- `GET /health` - Health check
- `GET /metrics` - Prometheus метрики
- `GET /swagger/index.html` - Swagger UI

## API Usage

### 1. User Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }'
```

### 2. Login and Get Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }'
```

Response contains `token` - use it in the header `Authorization: Bearer <token>`.

### 3. Create Product

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

### 4. Get Product List

```bash
curl -X GET "http://localhost:8080/api/v1/products?page=1&limit=10" \
  -H "Authorization: Bearer <your-token>"
```

## Deploy to AWS EC2

### 1. Prepare EC2 Instance

```bash
# Update system (Ubuntu/Debian)
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Create application directory
mkdir -p /opt/demo-service
```

### 2. Copy Files to Server

```bash
# From your local computer
scp -r . user@your-ec2-ip:/opt/demo-service/
```

### 3. Compile on Server

```bash
cd /opt/demo-service
go mod download
go build -o demo-service cmd/server/main.go
```

### 4. Create systemd Service

Create file `/etc/systemd/system/demo-service.service`:

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

### 5. Start Service

```bash
sudo systemctl daemon-reload
sudo systemctl enable demo-service
sudo systemctl start demo-service
sudo systemctl status demo-service
```

### 6. Configure Nginx (optional)

Create `/etc/nginx/sites-available/demo-service`:

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

Activate configuration:

```bash
sudo ln -s /etc/nginx/sites-available/demo-service /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 7. Configure Security Group

Make sure the following ports are open in the EC2 Security Group:
- 22 (SSH)
- 80 (HTTP, if using Nginx)
- 443 (HTTPS, if using SSL)
- 8080 (if accessing application directly)

## Monitoring

### Prometheus Metrics

Metrics are available at `/metrics` endpoint. You can configure Prometheus to collect metrics:

```yaml
scrape_configs:
  - job_name: 'demo-service'
    static_configs:
      - targets: ['localhost:8080']
```

### Logs

Application logs can be viewed through systemd:

```bash
sudo journalctl -u demo-service -f
```

Or if running directly, logs are output to stdout/stderr.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | 8080 |
| `DB_PATH` | Path to SQLite file | ./data/demo.db |
| `JWT_SECRET` | Secret key for JWT | (required) |
| `JWT_EXPIRY` | JWT token lifetime | 24h |
| `RATE_LIMIT_RPS` | Requests per second | 10 |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | info |

## Project Structure

```
demo-service/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/                 # Configuration
│   ├── handler/                # HTTP handlers
│   ├── service/                # Business logic
│   ├── repository/             # Database operations
│   ├── model/                  # Data models
│   ├── middleware/             # Middleware
│   ├── metrics/                # Prometheus metrics
│   └── database/               # Database initialization
├── pkg/jwt/                    # JWT utilities
├── migrations/                 # SQL migrations
└── docs/                       # Swagger documentation
```

## Development

### Generate Swagger Documentation

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go
```

### Run Tests

```bash
go test ./...
```

## License

MIT

