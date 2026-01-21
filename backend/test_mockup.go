package main

import (
    "fmt"
    "log"
    "minicronk/services"
)

func main() {
    // Test mockup generation
    mockupURL, err := services.GenerateMockup(1, "test-logo.png")
    if err != nil {
        log.Fatal("Error generating mockup:", err)
    }
    
    fmt.Printf("Mockup generated successfully: %s\n", mockupURL)
}