package config

import (
	"api-shiners/pkg/entities"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: .env file not found, using system environment variables")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	timezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}

	DB = db

	log.Println("🚀 Running AutoMigrate...")
	err = db.AutoMigrate(
		&entities.User{},
		&entities.Role{},
		&entities.UserRole{},
		&entities.Enrollment{},
		&entities.Material{},
		&entities.Quiz{},
		&entities.QuizAttempt{},
		&entities.Answer{},
		&entities.Choice{},
		&entities.Question{},
		&entities.Course{},
		&entities.CourseModule{},
		&entities.FeedbackQuestion{},
		&entities.FeedbackAnswer{},
	)
	if err != nil {
		log.Fatal("❌ Failed to migrate:", err)
	}

	log.Println("✅ Database connected and migrated successfully!")

	seedRoles(db)
}

func seedRoles(db *gorm.DB) {
	roles := []entities.Role{
		{
			Name:        entities.ADMIN,
			Description: "Administrator dengan akses penuh",
		},
		{
			Name:        entities.TEACHER,
			Description: "Guru yang dapat mengelola materi dan nilai",
		},
		{
			Name:        entities.STUDENT,
			Description: "Siswa yang dapat mengakses pembelajaran",
		},
	}

	for _, role := range roles {
		var existing entities.Role
		if err := db.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&role).Error; err != nil {
					log.Printf("❌ Gagal menambahkan role %s: %v", role.Name, err)
				} else {
					log.Printf("✅ Role %s berhasil ditambahkan", role.Name)
				}
			}
		}
	}
}
