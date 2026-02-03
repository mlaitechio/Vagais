# merv.one Backend

A comprehensive AI Agent Marketplace platform backend built with Go, featuring microservices architecture, graceful degradation, and enterprise-ready features.

## 🚀 Features

### Core Features
- **User Management**: Complete user registration, authentication, and profile management
- **Organization Support**: Multi-tenant organization management with role-based access
- **AI Agent Management**: Create, deploy, and manage AI agents with various configurations
- **Marketplace**: Agent discovery, reviews, ratings, and search functionality
- **Runtime Execution**: Real-time agent execution with streaming capabilities
- **Analytics**: Comprehensive usage tracking and business intelligence

### Enterprise Features
- **Licensing System**: Enterprise licensing with offline validation
- **Payment Processing**: Multi-gateway payment support (Stripe, PayPal, UPI)
- **Billing & Subscriptions**: Flexible subscription management
- **Admin Panel**: Complete administrative interface
- **Security**: JWT authentication, role-based access, rate limiting

### Technical Features
- **Graceful Degradation**: System works even when optional services are unavailable
- **Microservices Ready**: Modular architecture for easy scaling
- **Database Support**: PostgreSQL (production) and SQLite (development)
- **Caching**: Redis integration for performance
- **Search**: Elasticsearch for advanced search capabilities
- **File Storage**: MinIO for scalable file storage
- **Message Queue**: RabbitMQ for async processing

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   Load Balancer │
│   (React)       │◄──►│   (Gin)         │◄──►│   (Nginx)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Core Services                           │
├─────────────────┬─────────────────┬───────────────────────────┤
│  Auth Service   │  User Service   │  Agent Service            │
├─────────────────┼─────────────────┼───────────────────────────┤
│Marketplace Svc  │ Runtime Service │ Integration Service       │
├─────────────────┼─────────────────┼───────────────────────────┤
│Analytics Service│ Billing Service │ Notification Service      │
├─────────────────┼─────────────────┼───────────────────────────┤
│ License Service │ Payment Service │ Admin Service             │
└─────────────────┴─────────────────┴───────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Data Layer                                │
├─────────────────┬─────────────────┬───────────────────────────┤
│   PostgreSQL    │     Redis       │    Elasticsearch          │
├─────────────────┼─────────────────┼───────────────────────────┤
│     MinIO       │   RabbitMQ      │    File Storage           │
└─────────────────┴─────────────────┴───────────────────────────┘
```

## 🛠️ Technology Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL (production), SQLite (development)
- **ORM**: GORM
- **Authentication**: JWT
- **Caching**: Redis
- **Search**: Elasticsearch
- **File Storage**: MinIO
- **Message Queue**: RabbitMQ
- **Monitoring**: Prometheus
- **Documentation**: Swagger/OpenAPI

## 📦 Installation

### Prerequisites

- Go 1.21 or higher
- Git
- SQLite (for development) or PostgreSQL (for production)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/mlaitechio/merv.one.git
   cd agais.ai/backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

5. **Test the backend**
   ```bash
   go run test_main.go
   ```

### Environment Configuration

Create a `.env` file with the following variables:

```env
# Core Configuration
ENVIRONMENT=development
PORT=8080
DATABASE_TYPE=sqlite

# Database (SQLite for development)
DATABASE_URL=agais.db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRATION_HOURS=24

# Optional Services (can be disabled for graceful degradation)
REDIS_URL=redis://localhost:6379
ELASTICSEARCH_URL=http://localhost:9200
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Security
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://agais.ai
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# Payment (optional)
STRIPE_SECRET_KEY=sk_test_...
PAYPAL_CLIENT_ID=your-paypal-client-id
PAYPAL_CLIENT_SECRET=your-paypal-secret

# Email (optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

## 🚀 Running the Application

### Development Mode
```bash
# Run with hot reload (requires air)
air

# Or run directly
go run main.go
```

### Production Mode
```bash
# Build the application
go build -o agais-backend main.go

# Run the binary
./agais-backend
```

### Docker Deployment
```bash
# Build Docker image
docker build -t agais-backend .

# Run with Docker Compose
docker-compose up -d
```

## 📚 API Documentation

Once the server is running, you can access:

- **API Documentation**: `http://localhost:8080/swagger/index.html`
- **Health Check**: `http://localhost:8080/health`
- **Metrics**: `http://localhost:8080/metrics`

### API Endpoints

The API is organized into the following groups:

- **Authentication**: `/api/v1/auth/*`
- **Users**: `/api/v1/users/*`
- **Organizations**: `/api/v1/organizations/*`
- **Agents**: `/api/v1/agents/*`
- **Marketplace**: `/api/v1/marketplace/*`
- **Runtime**: `/api/v1/runtime/*`
- **Analytics**: `/api/v1/analytics/*`
- **Billing**: `/api/v1/billing/*`
- **Payments**: `/api/v1/payments/*`
- **Licenses**: `/api/v1/licenses/*`
- **Notifications**: `/api/v1/notifications/*`
- **Admin**: `/api/v1/admin/*`

## 🔧 Development

### Project Structure

```
backend/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── .env                    # Environment variables
├── env.example            # Environment template
├── internal/              # Internal application code
│   ├── config/           # Configuration management
│   ├── database/         # Database initialization
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Database models
│   ├── routes/           # Route definitions
│   └── services/         # Business logic services
├── test_main.go          # Test application
└── README.md             # This file
```

### Adding New Features

1. **Create a new service** in `internal/services/`
2. **Add models** in `internal/models/`
3. **Create handlers** in `internal/handlers/`
4. **Define routes** in `internal/routes/`
5. **Add middleware** if needed in `internal/middleware/`

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/services
```

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Role-Based Access Control**: Granular permissions system
- **Rate Limiting**: Protection against abuse
- **CORS Configuration**: Cross-origin request handling
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: GORM with parameterized queries
- **XSS Protection**: Security headers and input sanitization

## 📊 Monitoring & Observability

- **Health Checks**: `/health` endpoint
- **Metrics**: Prometheus metrics at `/metrics`
- **Logging**: Structured logging with different levels
- **Tracing**: Request ID tracking
- **Performance**: Response time monitoring

## 🚀 Deployment

### Docker Deployment

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o agais-backend main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/agais-backend .
EXPOSE 8080
CMD ["./agais-backend"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: agais-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: agais-backend
  template:
    metadata:
      labels:
        app: agais-backend
    spec:
      containers:
      - name: agais-backend
        image: agais-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENVIRONMENT
          value: "production"
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: [API Docs](http://localhost:8080/swagger/index.html)
- **Issues**: [GitHub Issues](https://github.com/mlaitechio/merv.one/issues)
- **Discussions**: [GitHub Discussions](https://github.com/mlaitechio/merv.one/discussions)

## 🎯 Roadmap

- [ ] WebSocket support for real-time features
- [ ] GraphQL API
- [ ] Advanced analytics dashboard
- [ ] Multi-region deployment
- [ ] Advanced caching strategies
- [ ] Machine learning model management
- [ ] Advanced security features
- [ ] Performance optimizations

---

**Built with ❤️ for the AI community**