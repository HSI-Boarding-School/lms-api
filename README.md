1. Install Dependencies

Pastikan sudah install Go 1.24+.
Jalankan:

go mod download

2. Setup Environment

Buat file .env di root project. Contoh:

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=shiners_lms
JWT_SECRET=your-secret-key

3. Setup Database & Migration

Buat database baru (misal: shiners_lms) sesuai dengan .env.

Jalankan migration dengan golang-migrate
:

migrate -path migrations -database "postgres://<username>:<password>@localhost:5432/shiners_lms?sslmode=disable" up

4. Jalankan Server
   go run cmd/server/main.go

Server default berjalan di:
👉 http://localhost:8000

📌 List API (sementara)
Auth

POST /auth/register → Register user baru (default role: STUDENT)

POST /auth/login → Login user → return access_token & refresh_token & info

🧪 Example Request & Response

Register

Request

POST /auth/register
Content-Type: application/json

{
"name": "John Doe",
"email": "john@example.com",
"password": "secret123"
}

Login

Request

POST /auth/login
Content-Type: application/json

{
"email": "john@example.com",
"password": "secret123"
}
