package tests

import (
	"bookmarket/database"
	"bookmarket/internal/app/users/models"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// load .env file
	godotenv.Load("../../../../.env")
	db := database.SetupTestDB()

	// Migration table user
	db.AutoMigrate(&models.Users{})

	// Run test
	m.Run()

	db.Migrator().DropTable(&models.Users{})
}
