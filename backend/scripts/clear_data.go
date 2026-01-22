package main

import (
	"log"
	"printflow/db"
	"printflow/models"
)

func main() {
	// Connect to database
	db.Connect()

	log.Println("Clearing all data from database...")

	// Delete all records from tables (keeps table structure)
	if err := db.DB.Exec("DELETE FROM assets").Error; err != nil {
		log.Printf("Error clearing assets: %v", err)
	} else {
		log.Println("✓ Cleared assets table")
	}

	if err := db.DB.Exec("DELETE FROM orders").Error; err != nil {
		log.Printf("Error clearing orders: %v", err)
	} else {
		log.Println("✓ Cleared orders table")
	}

	// Reset auto-increment counters (SQLite specific)
	if err := db.DB.Exec("DELETE FROM sqlite_sequence WHERE name IN ('orders', 'assets')").Error; err != nil {
		log.Printf("Warning: Could not reset auto-increment counters: %v", err)
	} else {
		log.Println("✓ Reset auto-increment counters")
	}

	// Verify tables are empty
	var orderCount, assetCount int64
	db.DB.Model(&models.Order{}).Count(&orderCount)
	db.DB.Model(&models.Asset{}).Count(&assetCount)

	log.Printf("Database cleared successfully!")
	log.Printf("Orders: %d records", orderCount)
	log.Printf("Assets: %d records", assetCount)
	log.Println("Tables structure preserved - ready for fresh data")
}