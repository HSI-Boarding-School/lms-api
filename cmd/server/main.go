package main

import (
	"log"
	"github.com/daffa-fawwaz/shiners-lms-backend/app"
)

func main() {
	appServer := app.SetupApp()

	log.Println("🚀 Server running on :8000")
	log.Fatal(appServer.Listen(":8000"))
}
