 # merv.one - AI Agent Marketplace Platform

A comprehensive SaaS platform for hosting, discovering, and managing AI agents with a futuristic 2026 design.

## 🚀 Features

- **AI Agent Marketplace**: Discover, deploy, and manage intelligent AI agents
- **Multi-tenant Architecture**: Support for organizations and individual users
- **Graceful Feature Fallbacks**: Works seamlessly even when optional services are disabled
- **Futuristic UI**: 2026-inspired design with animations and modern UX
- **Enterprise Ready**: Security, licensing, and payment integration
- **Scalable Backend**: Microservices architecture with Go
- **Real-time Features**: WebSocket support and live updates

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL (production) / SQLite (development)
- **Cache**: Redis (optional)
- **Search**: Elasticsearch (optional)
- **Storage**: MinIO (optional)
- **Message Queue**: RabbitMQ (optional)

### Frontend
- **Framework**: React 18+ with TypeScript
- **UI Library**: Material-UI with custom theme
- **Animations**: Framer Motion
- **State Management**: Zustand
- **Routing**: React Router DOM

## 📦 Installation

### Prerequisites
- Go 1.21+
- Node.js 18+
- SQLite (for local development)
- Git

### Backend Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/mlaitechio/vagais.git
   cd agais
   ```

2. **Navigate to backend directory**
   ```bash
   cd backend
   ```

3. **Install Go dependencies**
   ```bash
   go mod tidy
   ```

4. **Create environment file**
   ```bash
   cp .env.example .env
   ```

5. **Configure environment variables**
   ```env
   ENVIRONMENT=development
   PORT=8080
   DATABASE_TYPE=sqlite
   JWT_SECRET_KEY=your-secret-key-change-in-production
   ```

6. **Run the backend**
   ```bash
   go run main.go
   ```

The backend will start on `http://localhost:8080` with SQLite database.

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Start development server**
   ```bash
   npm run dev
   ```

The frontend will start on `http://localhost:3000`

## 🔧 Configuration

### Backend Configuration

The application supports graceful fallbacks for optional services:

- **Redis**: Caching and session management (optional)
- **Elasticsearch**: Search functionality (optional)
- **MinIO**: File storage (optional)
- **RabbitMQ**: Message queuing (optional)
- **License Management**: On-premise licensing (optional)
- **Payment Processing**: Stripe/PayPal integration (optional)

### Environment Variables

```env
# Core Configuration
ENVIRONMENT=development
PORT=8080
DATABASE_TYPE=sqlite

# Database (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=agais
DB_SSLMODE=disable

# JWT
JWT_SECRET_KEY=your-secret-key
JWT_EXPIRATION_HOURS=24

# Optional Services
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

ELASTICSEARCH_URL=http://localhost:9200
ELASTICSEARCH_USERNAME=
ELASTICSEARCH_PASSWORD=

MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY_ID=minioadmin
MINIO_SECRET_ACCESS_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=agais

RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Security
ALLOWED_DOMAINS=*
BLOCKED_DOMAINS=
RATE_LIMIT=100
MAX_FILE_SIZE=10485760

# Payment (Optional)
STRIPE_SECRET_KEY=
STRIPE_PUBLISHABLE_KEY=
PAYPAL_CLIENT_ID=
PAYPAL_SECRET=

# Email (Optional)
SMTP_HOST=
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
FROM_EMAIL=noreply@agais.ai
```

## 🏗️ Architecture

### Backend Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                            Load Balancer                            │
├─────────────────────────────────────────────────────────────────────┤
│                            API Gateway                              │
├─────────────────────────────┬───────────────────────────────────────┤
│        Web Frontend         │           Mobile Apps                 │
├─────────────────────────────┴───────────────────────────────────────┤
│                         Microservices Layer                         │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┤
│ User        │ Agent       │ Marketplace │ Runtime     │ Integration │
│ Service     │ Service     │ Service     │ Service     │ Service     │
├─────────────┼─────────────┼─────────────┼─────────────┼─────────────┤
│                         Message Bus (Kafka)                         │
├─────────────────────────────────────────────────────────────────────┤
│                         Data Layer                                  │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┤
│ PostgreSQL  │ Redis       │ Elasticsearch│ MinIO       │ External    │
│ (Primary)   │ (Cache)     │ (Search)    │ (Storage)   │ APIs/LLMs   │
└─────────────┴─────────────┴─────────────┴─────────────┴─────────────┘
```

### Frontend Architecture

- **React 18+**: Modern React with hooks and functional components
- **Material-UI**: Custom themed components with futuristic design
- **Framer Motion**: Smooth animations and transitions
- **React Router**: Client-side routing
- **Zustand**: Lightweight state management
- **Axios**: HTTP client for API communication

## 🎨 Design System

### Color Palette
- **Primary**: `#98177E` (Purple)
- **Secondary**: `#00D4FF` (Cyan)
- **Success**: `#00FF88` (Green)
- **Warning**: `#FFB800` (Yellow)
- **Error**: `#FF4757` (Red)
- **Background**: `#0A0A0A` (Dark)
- **Surface**: `#1A1A1A` (Dark Gray)

### Typography
- **Font Family**: Inter, Roboto, Helvetica, Arial
- **Gradient Text**: Primary to secondary color gradients
- **Responsive**: Mobile-first design approach

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Role-based Access Control**: Admin, staff, maintainer, user roles
- **Rate Limiting**: API rate limiting to prevent abuse
- **CORS Protection**: Cross-origin resource sharing configuration
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: Parameterized queries
- **XSS Protection**: Content Security Policy headers

## 🚀 Deployment

### Local Development
```bash
# Backend
cd backend
go run main.go

# Frontend
cd frontend
npm run dev
```

### Production Deployment

1. **Backend**: Deploy to cloud (AWS, Azure, GCP) with Docker
2. **Frontend**: Build and deploy to CDN
3. **Database**: Use managed PostgreSQL service
4. **Optional Services**: Deploy based on requirements

## 📚 API Documentation

The API documentation is available at:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI Spec**: `http://localhost:8080/swagger/doc.json`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- **Email**: support@agais.ai
- **Documentation**: [docs.agais.ai](https://docs.agais.ai)
- **Issues**: [GitHub Issues](https://github.com/mlaitechio/vagais/issues)

## 🔮 Roadmap

- [ ] Multi-language support
- [ ] Mobile app development
- [ ] Advanced analytics dashboard
- [ ] AI agent marketplace features
- [ ] Enterprise SSO integration
- [ ] Advanced security features
- [ ] Performance optimizations
- [ ] Community features

---

**Built with ❤️ for the future of AI agents**