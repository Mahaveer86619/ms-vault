package main

import (
	"github.com/Mahaveer86619/ms/auth/pkg/config"
	"github.com/Mahaveer86619/ms/auth/pkg/db"
	"github.com/Mahaveer86619/ms/auth/pkg/models"
	"github.com/labstack/gommon/log"
)

var version string = "dev"

func main() {
	log.Printf("Starting migration, version %s...", version)

	config.InitConfig()
	db.InitDB()

	tables := []interface{}{
		&models.UserProfile{},
		&models.UserSecurity{},
	}

	log.Info("Running AutoMigrate...")
	if err := db.DB.AutoMigrate(tables...); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Info("DB Migration completed successfully!")
}
