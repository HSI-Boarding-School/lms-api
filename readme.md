1. inisialisasi project

    go mod init api-shiners


2. install dependencies

    go get github.com/gofiber/fiber/v2
    go get -u gorm.io/gorm
    go get gorm.io/driver/postgres
    go get github.com/joho/godotenv
    go get github.com/golang-jwt/jwt/v5
    go get -u github.com/swaggo/swag/cmd/swag
    go get -u github.com/gofiber/swagger
    go get -u github.com/swaggo/files
    go get github.com/gofiber/fiber/v2/middleware/cors
    go get github.com/redis/go-redis/v9
    go get github.com/stretchr/testify
    go get gopkg.in/yaml.v3@latest
    go mod tidy



3. Copy file .env.example menjadi .env

    cp .env.example .env



4. Sesuaikan konfigurasi di .env:

    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_NAME=db_shiners
    DB_PORT=5432

    REDIS_HOST=localhost:6379
    REDIS_PASSWORD=
    REDIS_DB=0


5. Pastikan PostgreSQL sudah terinstal. Jika belum, unduh dari https://www.postgresql.org/download/


6. Buat database baru:

    CREATE DATABASE db_shiners;


7. Menjalankan Aplikasi
    
    Pastikan berada di root project, lalu jalankan: go run main.go