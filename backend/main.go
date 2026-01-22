package main

import (
    "log"
    "os"
    "time"
    "printflow/db"
    "printflow/handlers"
    "printflow/models"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    db.Connect()
    db.DB.AutoMigrate(&models.Order{}, &models.Asset{})

    r := gin.Default()

    // Add middleware for timestamp
    r.Use(func(c *gin.Context) {
        c.Set("timestamp", time.Now().Unix())
        c.Next()
    })

	r.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
    c.Next()
})

    // API routes
    r.POST("/orders", handlers.CreateOrder)
    r.GET("/orders", handlers.ListOrders)
    r.GET("/orders/:ID", handlers.GetOrder)
    r.POST("/orders/:ID/approve", handlers.ApproveOrder)
    r.POST("/orders/:ID/mockup", handlers.GenerateMockupHandler)
	r.POST("/orders/:ID/label", handlers.GenerateLabel)
	r.GET("/colors", handlers.GetAvailableColors)


    
    // Upload route
    r.POST("/upload/logo", handlers.UploadLogo)

    // Static file serving
    r.Static("/mockups", "./mockups")
    r.Static("/uploads", "./uploads")
    r.Static("/assets", "./assets")
	r.Static("/labels", "./labels")



    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    r.Run(":" + port)
}


