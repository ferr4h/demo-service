# Demo Service - Production-Ready Go CRUD Application

Production-ready Go application with CRUD operations for product management, including JWT authentication, Prometheus metrics, Swagger documentation, rate limiting, and structured logging. Now uses PostgreSQL database.

## Features

- ✅ REST API for CRUD operations with products
- ✅ JWT authentication
- ✅ Prometheus metrics
- ✅ Swagger/OpenAPI documentation
- ✅ Rate limiting
- ✅ Structured logging
- ✅ PostgreSQL database
- ✅ Docker & Docker Compose support
- ✅ Kubernetes deployment
- ✅ Graceful shutdown
- ✅ Health check endpoints

## Requirements

- Go 1.21 or higher
- Docker & Docker Compose (for containerized deployment)
- Kubernetes cluster (for Kubernetes deployment)
- PostgreSQL (when running locally without Docker)

## Installation and Running

### Option 1: Docker Compose (Recommended for Development)

#### 1. Clone Repository

```bash
git clone <your-repo-url>
cd demo-service
```

#### 2. Start with Docker Compose

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Demo service on port 8080

#### 3. Access the Application

The application will be available at `http://localhost:8080`

### Option 2: Local Development

#### 1. Start PostgreSQL

You can use Docker to run PostgreSQL locally:

```bash
docker run -d \
  --name postgres-demo \
  -e POSTGRES_DB=demo \
  -e POSTGRES_USER=demo \
  -e POSTGRES_PASSWORD=demo \
  -p 5432:5432 \
  postgres:15-alpine
```

#### 2. Configure Environment

Create `.env` file:

```bash
# Server Configuration
SERVER_PORT=8080

# Database Configuration (PostgreSQL)
DATABASE_URL=postgres://demo:demo@localhost:5432/demo?sslmode=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY=24h

# Rate Limiting
RATE_LIMIT_RPS=10

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=text
```

#### 3. Install Dependencies and Run

```bash
go mod download
go run cmd/server/main.go
```

### Option 3: Build and Run Binary

```bash
go build -o demo-service cmd/server/main.go
./demo-service
```

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

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster (local or cloud)
- `kubectl` configured to access your cluster

### 1. Deploy to Kubernetes

```bash
# Apply all manifests (PostgreSQL + Application)
kubectl apply -f k8s-complete.yaml

# Or deploy separately
kubectl apply -f k8s-postgres.yaml
kubectl apply -f deployment.yaml
```

### 2. Check Deployment Status

```bash
# Check pods
kubectl get pods

# Check services
kubectl get services

# Check persistent volumes
kubectl get pvc
```

### 3. Access the Application

```bash
# Get service URL (if using LoadBalancer)
kubectl get service demo-service

# Port forward for local access
kubectl port-forward svc/demo-service 8080:80
```

Then access at `http://localhost:8080`

### 4. View Logs

```bash
# Application logs
kubectl logs -l app=demo-service

# PostgreSQL logs
kubectl logs -l app=postgres
```

### 5. Cleanup

```bash
kubectl delete -f k8s-complete.yaml
```

## AWS EC2 Deployment

### 1. Prepare EC2 Instance

```bash
# Update system (Ubuntu/Debian)
sudo apt update && sudo apt upgrade -y

# Install PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# Configure PostgreSQL
sudo -u postgres psql -c "CREATE DATABASE demo;"
sudo -u postgres psql -c "CREATE USER demo WITH PASSWORD 'demo';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE demo TO demo;"

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Create application directory
mkdir -p /opt/demo-service
```

### 2. Configure Database Connection

Update the `.env` file with PostgreSQL connection:

```bash
# Database Configuration (PostgreSQL)
DATABASE_URL=postgres://demo:demo@localhost:5432/demo?sslmode=disable
```

### 3. Copy Files to Server

```bash
# From your local computer
scp -r . user@your-ec2-ip:/opt/demo-service/
```

### 4. Compile on Server

```bash
cd /opt/demo-service
go mod download
go build -o demo-service cmd/server/main.go
```

### 5. Create systemd Service

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

### 6. Start Service

```bash
sudo systemctl daemon-reload
sudo systemctl enable demo-service
sudo systemctl start demo-service
sudo systemctl status demo-service
```

### 7. Configure Nginx (optional)

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

### 8. Configure Security Group

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
| `DATABASE_URL` | PostgreSQL connection URL | postgres://demo:demo@localhost:5432/demo?sslmode=disable |
| `JWT_SECRET` | Secret key for JWT | (required) |
| `JWT_EXPIRY` | JWT token lifetime | 24h |
| `RATE_LIMIT_RPS` | Requests per second | 10 |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | info |
| `LOG_FORMAT` | Log format (text, json) | text |

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
│   └── database/               # Database initialization (PostgreSQL)
├── pkg/jwt/                    # JWT utilities
├── docker-compose.yml          # Docker Compose setup
├── k8s-complete.yaml           # Kubernetes manifests
├── Dockerfile                  # Docker image
├── go.mod                      # Go modules
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




