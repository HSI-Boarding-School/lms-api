# 🛠️ Shiners LMS API

**Shiners LMS** is the backend API for a **Learning Management System**. Built with **Go + Fiber**, it exposes endpoints for authentication and core LMS features backed by **PostgreSQL** and SQL migrations via **golang-migrate**.

***

## 🚀 Tech Stack

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-2C3E50?style=for-the-badge&logo=fiber&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-6A1B9A?style=for-the-badge)
![golang-migrate](https://img.shields.io/badge/migrate-SQL_Migrations-4B8BBE?style=for-the-badge)

```
lms-api/
│── cmd/server/main.go   # Application entry
│── app/
│   │── controllers/     # HTTP handlers
│   │── services/        # Business logic
│   │── repositories/    # Data access (GORM)
│   │── models/          # Domain models
│   │── routes/          # Route groups
│   │── middlewares/     # (optional) middleware
│   │── utils/           # Helpers (JWT, etc.)
│   │── app.go           # App wiring
│── config/config.go     # DB connection
│── migrations/          # SQL migrations
│── .env.example         # Env template
│── go.mod
```

***

## ⚙️ Installation & Setup

1. Clone and enter the repo
   ```bash
   git clone git@github.com:HSI-Boarding-School/lms-api.git
   cd lms-api
   ```

2. Configure environment
   ```bash
   cp .env.example .env
   # Update DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
   # Optionally set DATABASE_URL (for migrations)
   ```

3. Prepare database
   ```bash
   createdb shiners_lms   # or create via your DB tool
   ```

4. Install dependencies
   ```bash
   go mod download
   ```

5. Run migrations (golang-migrate)
   - Install: `brew install golang-migrate` (macOS) or download from releases
   - Run:
     ```bash
     migrate -path migrations -database "$DATABASE_URL" up
     # Example DATABASE_URL:
     # postgres://user:pass@localhost:5432/shiners_lms?sslmode=disable
     ```
   - If you see an error about `pgcrypto`, enable it once:
     ```bash
     psql "$DATABASE_URL" -c 'CREATE EXTENSION IF NOT EXISTS "pgcrypto";'
     ```

6. Start the API
   ```bash
   go run ./cmd/server
   # Server listens at http://localhost:8000
   ```

***

## 🔐 Auth Endpoints

- POST `/auth/register`
  - Body: `{ "name": string, "email": string, "password": string }`

- POST `/auth/login`
  - Body: `{ "email": string, "password": string }`
  - Returns: `{ access_token, refresh_token, info: { role, email, name } }`

Note: Seed roles first (ADMIN, TEACHER, STUDENT) if your environment is empty.

```sql
-- Seed roles (run once)
INSERT INTO roles (name, description) VALUES
  ('ADMIN','Administrator'),
  ('TEACHER','Teacher'),
  ('STUDENT','Student');
```

***

## 🏫 About the Project

Shiners LMS is designed to support digital learning and management activities within the **HSI Boarding School network** across Indonesia.

It provides a centralized platform for educational materials, quizzes, and course progress tracking for both teachers and students.

***

## 📌 Notes

- Server port is currently hardcoded to `8000` in `cmd/server/main.go`.
- DB connection uses `DB_*` variables; migrations use `DATABASE_URL`.
- Consider aligning JWT secrets in code with `JWT_SECRET` from environment.

