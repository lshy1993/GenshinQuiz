# Genshin Quiz Go Backend

A modern Go backend API for the Genshin Impact Quiz application, built with Go-Chi, OpenAPI, Go-Jet, PostgreSQL, and Goose migrations.

## 🏗️ Architecture

- **Web Framework**: Go-Chi for HTTP routing and middleware
- **API Documentation**: OpenAPI 3.0 with automatic code generation
- **Database**: PostgreSQL with Go-Jet as query builder
- **Migrations**: Goose for database schema management
- **Containerization**: Docker and Docker Compose
- **Code Generation**: Automatic API and model generation

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later
- PostgreSQL (or Docker for containerized setup)
- Git

### 1. Setup Development Environment

```bash
# Run the setup script
./scripts/setup.sh
```

This script will:
- Check Go installation
- Create `.env` file from template
- Install dependencies and dev tools
- Make scripts executable

### 2. Configure Environment

Edit the `.env` file with your settings:

```bash
cp .env.example .env
# Edit .env with your database credentials and preferences
```

### 3. Start Database

**Option A: Using Docker (Recommended)**
```bash
docker-compose up postgres -d
```

**Option B: Local PostgreSQL**
- Ensure PostgreSQL is running
- Create database: `createdb genshin_quiz`

### 4. Run Migrations

```bash
# Apply all migrations
./scripts/migrate.sh up

# Check migration status
./scripts/migrate.sh status
```

### 5. Generate Code

```bash
# Generate Go-Jet models from database
./scripts/generate_models.sh

# Generate OpenAPI code
./scripts/generate_api.sh
```

### 6. Start the Server

```bash
# Development mode
go run main.go

# Or with hot reloading (if air is installed)
air
```

The API will be available at: `http://localhost:8080`

## 📁 Project Structure

```
go-backend/
├── api/                    # OpenAPI specifications
│   └── openapi.yaml
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── database/         # Database connection and utilities
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/          # Data models and DTOs
│   ├── repository/      # Data access layer
│   ├── services/        # Business logic layer
│   └── generated/       # Generated code (API, models)
├── migrations/           # Database migrations
├── scripts/             # Utility scripts
├── docker-compose.yml   # Docker services
├── Dockerfile          # Container definition
├── go.mod              # Go module definition
└── main.go            # Application entry point
```

## 🛠️ Available Scripts

### Database Operations

```bash
# Apply migrations
./scripts/migrate.sh up

# Rollback last migration  
./scripts/migrate.sh down

# Check migration status
./scripts/migrate.sh status

# Create new migration
./scripts/migrate.sh create add_new_feature

# Reset database (⚠️ destructive)
./scripts/migrate.sh reset
```

### Code Generation

```bash
# Generate Go-Jet models from database
./scripts/generate_models.sh

# Generate OpenAPI server code
./scripts/generate_api.sh
```

### Docker Operations

```bash
# Start all services
docker-compose up

# Start only database
docker-compose up postgres -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down
```

## 🔌 API Endpoints

### Health Check
- `GET /health` - Service health status

### Users
- `GET /api/v1/users` - List users (with pagination)
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Quizzes
- `GET /api/v1/quizzes` - List quizzes (with filtering)
- `POST /api/v1/quizzes` - Create quiz
- `GET /api/v1/quizzes/{id}` - Get quiz by ID
- `PUT /api/v1/quizzes/{id}` - Update quiz
- `DELETE /api/v1/quizzes/{id}` - Delete quiz

## 🗄️ Database Schema

The application uses PostgreSQL with the following main tables:

- **users**: User accounts and profiles
- **quizzes**: Quiz definitions with metadata
- **questions**: Individual quiz questions
- **quiz_attempts**: User quiz attempt records
- **user_answers**: Individual question answers

### Enums
- `quiz_category`: characters, weapons, artifacts, lore, gameplay
- `quiz_difficulty`: easy, medium, hard
- `question_type`: multiple_choice, true_false, fill_in_blank

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Test API endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users
```

## 🚢 Deployment

### Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose up --build

# Production deployment
docker-compose -f docker-compose.prod.yml up -d
```

### Manual Deployment

1. Build the binary:
```bash
CGO_ENABLED=0 GOOS=linux go build -o main .
```

2. Run migrations:
```bash
./scripts/migrate.sh up
```

3. Start the server:
```bash
./main
```

## 🛡️ Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:password@localhost/genshin_quiz?sslmode=disable` |
| `PORT` | Server port | `8080` |
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |
| `LOG_LEVEL` | Logging level | `info` |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make changes and test thoroughly
4. Run code generation: `./scripts/generate_api.sh && ./scripts/generate_models.sh`
5. Commit changes: `git commit -am 'Add new feature'`
6. Push to branch: `git push origin feature/new-feature`
7. Submit a Pull Request

## 📝 Development Workflow

1. **Make database changes**: Create migration with `./scripts/migrate.sh create migration_name`
2. **Apply migrations**: Run `./scripts/migrate.sh up`
3. **Update models**: Run `./scripts/generate_models.sh` to regenerate Go-Jet models
4. **Update API**: Modify `api/openapi.yaml` if needed
5. **Regenerate API code**: Run `./scripts/generate_api.sh`
6. **Implement business logic**: Update services, handlers, etc.
7. **Test**: Verify endpoints work correctly

## 🔧 Tech Stack

- **Language**: Go 1.21
- **Web Framework**: Chi v5
- **Database**: PostgreSQL 15
- **Query Builder**: Go-Jet v2
- **Migrations**: Goose v3
- **API Spec**: OpenAPI 3.0
- **Code Generation**: oapi-codegen
- **Containerization**: Docker & Docker Compose
- **Logging**: Logrus
- **Environment**: godotenv

## 📚 Additional Resources

- [Go-Chi Documentation](https://go-chi.io/)
- [Go-Jet Documentation](https://github.com/go-jet/jet)
- [Goose Migration Tool](https://github.com/pressly/goose)
- [OpenAPI Specification](https://swagger.io/specification/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.