package handlers

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/gin-gonic/gin"
)

// UploadLogo handles logo file uploads
func UploadLogo(c *gin.Context) {
    // Parse multipart form
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }
    defer file.Close()

    // Validate file type
    if !isValidImageType(header.Filename) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only PNG, JPG, JPEG allowed"})
        return
    }

    // Create uploads directory if it doesn't exist
    uploadsDir := "uploads"
    if err := os.MkdirAll(uploadsDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory"})
        return
    }

    // Generate unique filename
    filename := fmt.Sprintf("logo_%d_%s", 
        c.GetInt64("timestamp"), 
        sanitizeFilename(header.Filename))
    
    filepath := filepath.Join(uploadsDir, filename)

    // Create destination file
    dst, err := os.Create(filepath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
        return
    }
    defer dst.Close()

    // Copy uploaded file to destination
    if _, err := io.Copy(dst, file); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  "File uploaded successfully",
        "filename": filename,
        "path":     filepath,
        "url":      "/" + filepath,
    })
}

// isValidImageType checks if the file has a valid image extension
func isValidImageType(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    validExts := []string{".png", ".jpg", ".jpeg", ".gif"}
    
    for _, validExt := range validExts {
        if ext == validExt {
            return true
        }
    }
    return false
}

// sanitizeFilename removes potentially dangerous characters from filename
func sanitizeFilename(filename string) string {
    // Replace spaces and special characters
    filename = strings.ReplaceAll(filename, " ", "_")
    filename = strings.ReplaceAll(filename, "..", "")
    filename = strings.ReplaceAll(filename, "/", "")
    filename = strings.ReplaceAll(filename, "\\", "")
    
    return filename
}