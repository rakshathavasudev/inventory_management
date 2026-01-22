package db

import (
    "log"

    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
    "printflow/models"
)

var DB *gorm.DB

func Connect() {
    database, err := gorm.Open(sqlite.Open("printflow.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database")
    }

    // Auto migrate the schema
    database.AutoMigrate(&models.Order{}, &models.Asset{})

    DB = database
    log.Println("Database connected and migrated successfully")
}
