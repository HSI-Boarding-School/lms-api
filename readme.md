# Shiners LMS API

REST API untuk sistem Learning Management System (LMS) Shiners yang dibangun dengan Go Fiber.

## Tech Stack

- **Framework**: [Go Fiber](https://gofiber.io/) v2
- **Database**: PostgreSQL dengan GORM
- **Cache**: Redis
- **Authentication**: JWT (JSON Web Token)
- **Documentation**: Swagger (Swaggo)

## Struktur Project

```
api-shiners/
├── api/
│   ├── handlers/          # HTTP handlers/controllers
│   │   ├── dto/           # Data Transfer Objects
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── feedback_handler.go
│   │   └── health_handler.go
│   └── routes/            # Route definitions
├── pkg/
│   ├── auth/              # Auth service & repository
│   ├── user/              # User service & repository
│   ├── feedback/          # Feedback service & repository
│   ├── entities/          # Database models
│   ├── middleware/        # Auth, Admin, Teacher middleware
│   ├── config/            # Database & Redis config
│   └── utils/             # Helper functions
├── docs/                  # Swagger documentation
├── main.go
├── go.mod
└── go.sum
```

## Instalasi

### Prerequisites

- Go 1.24+
- PostgreSQL
- Redis (optional, untuk session management)

### Setup

1. Clone repository
```bash
git clone https://github.com/HSI-Boarding-School/api-shiners.git
cd api-shiners
```

2. Copy environment file
```bash
cp .env.example .env
```

3. Konfigurasi `.env`
```env
APP_PORT=3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=shiners_lms

# JWT
JWT_SECRET=your-secret-key

# Redis (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

4. Install dependencies
```bash
go mod download
```

5. Jalankan aplikasi
```bash
go run main.go
```

### Docker

Jalankan dengan Docker Compose:
```bash
docker-compose up -d --build
```

## API Documentation

Swagger UI tersedia di: `http://localhost:3000/swagger/`

### Endpoints

#### Auth
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register user baru |
| POST | `/api/auth/login` | Login dan dapatkan JWT token |
| POST | `/api/auth/logout` | Logout (invalidate token) |
| POST | `/api/auth/forgot-password` | Request reset password |
| POST | `/api/auth/reset-password` | Reset password dengan token |

#### Users
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/users` | Get semua users (paginated) | Admin |
| GET | `/api/users/:id` | Get user by ID | User |
| POST | `/api/users/:id/role` | Set role user | Admin |
| POST | `/api/users/:id/activate` | Aktivasi user | Admin |
| POST | `/api/users/:id/deactivate` | Nonaktifkan user | Admin |
| GET | `/api/profile` | Get profile sendiri | User |
| PUT | `/api/profile` | Update profile sendiri | User |

#### Feedback
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/feedback/questions` | Buat pertanyaan feedback | Teacher |
| POST | `/api/feedback/answers` | Submit jawaban feedback | User |
| GET | `/api/feedback/teacher` | Get feedback milik teacher | Teacher |
| GET | `/api/feedback/questions/:teacher_id` | Get feedback by teacher ID | User |

#### Health Check
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health/database` | Cek koneksi database |
| GET | `/api/health/redis` | Cek koneksi Redis |

## Authentication

API menggunakan JWT Bearer Token untuk autentikasi.

```bash
# Login untuk mendapatkan token
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'

# Gunakan token untuk request yang memerlukan auth
curl http://localhost:3000/api/profile \
  -H "Authorization: Bearer <your-token>"
```

## User Roles

| Role | Permissions |
|------|-------------|
| ADMIN | Full access, manage users & roles |
| TEACHER | Create feedback questions, view answers |
| STUDENT | Submit feedback answers, view profile |

## Response Format

### Success Response
```json
{
  "status": 200,
  "message": "Operation successful",
  "data": { ... },
  "timestamp": "2025-01-01T00:00:00Z",
  "path": "/api/endpoint"
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": "ErrorType",
  "statusCode": 400,
  "timestamp": "2025-01-01T00:00:00Z",
  "path": "/api/endpoint"
}
```

## Development

### Generate Swagger Docs
```bash
go run github.com/swaggo/swag/cmd/swag@latest init
```

### Run Tests
```bash
go test ./...
```

### Build
```bash
go build -o shiners-api main.go
```

## License

MIT License
