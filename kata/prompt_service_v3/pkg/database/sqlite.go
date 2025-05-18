package database

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/walterfan/prompt-service/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dbfile string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	if err := DB.AutoMigrate(&models.Prompt{}, &models.User{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	InitData()
}

func InitData() {
	var count int64
	DB.Model(&models.Prompt{}).Count(&count)
	if count == 0 {
		// Load prompts from config
		var prompts []models.Prompt
		if err := viper.UnmarshalKey("prompts", &prompts); err != nil {
			log.Fatalf("Unable to decode into struct: %v", err)
		}

		result := DB.Create(&prompts)
		if result.Error != nil {
			log.Fatal("Failed to insert initial prompt data: ", result.Error)
		}
	}

	defaultUsername := os.Getenv("DEFAULT_USERNAME")
	defaultPassword := os.Getenv("DEFAULT_PASSWORD")
	defaultEmail := os.Getenv("DEFAULT_EMAIL")

	pwdHash, _ := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	// Check if users already exist
	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		users := []models.User{
			{
				Username:  defaultUsername,
				Password:  string(pwdHash),
				Email:     defaultEmail,
				Role:      "admin",
				ExpiredAt: time.Now().AddDate(1, 0, 0), // valid for 1 year
			},
		}
		DB.Create(&users)
		log.Println("Initialized database with sample users.")
	}

	log.Println("Initialized database with sample Go-related prompts.")
}
