package services

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "os"
    "path/filepath"
    "net/http"
    // "io"
)

// MockupConfig holds configuration for mockup generation
type MockupConfig struct {
    LogoPosition image.Point
    LogoMaxSize  image.Point
    OutputDir    string
}

// GetMockupConfig returns configuration based on product type
func GetMockupConfig(product string) MockupConfig {
    configs := map[string]MockupConfig{
        "Hoodie": {
            LogoPosition: image.Pt(220, 200), // Center chest area
            LogoMaxSize:  image.Pt(150, 150), // Max logo dimensions
            OutputDir:    "mockups",
        },
        "T-Shirt": {
            LogoPosition: image.Pt(200, 180),
            LogoMaxSize:  image.Pt(120, 120),
            OutputDir:    "mockups",
        },
    }
    
    if config, exists := configs[product]; exists {
        return config
    }
    
    // Default config
    return MockupConfig{
        LogoPosition: image.Pt(220, 200),
        LogoMaxSize:  image.Pt(150, 150),
        OutputDir:    "mockups",
    }
}

// GenerateMockup creates a product mockup by compositing a logo onto a product template
func GenerateMockup(orderID uint, logoURL string) (string, error) {
    src := "." + logoURL

    dst := fmt.Sprintf("mockups/order_%d.png", orderID)

    // copy logo into mockup location (simulate mockup)
    data, err := os.ReadFile(src)
    if err != nil {
        return "", err
    }

    if err := os.WriteFile(dst, data, 0644); err != nil {
        return "", err
    }

    return "/" + dst, nil
}


// loadTemplate loads the product template image
func loadTemplate(templatePath string) (image.Image, error) {
    // Check if template exists, if not create a placeholder
    if _, err := os.Stat(templatePath); os.IsNotExist(err) {
        return createPlaceholderTemplate(), nil
    }
    
    file, err := os.Open(templatePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    img, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }
    
    return img, nil
}

// loadLogo loads the logo image from file or URL
func loadLogo(logoPath string) (image.Image, error) {
    // Handle URL logos
    if isURL(logoPath) {
        return loadLogoFromURL(logoPath)
    }
    
    // Handle local file logos
    return loadLogoFromFile(logoPath)
}

// isURL checks if the path is a URL
func isURL(path string) bool {
    return len(path) > 4 && (path[:4] == "http" || path[:5] == "https")
}

// loadLogoFromURL downloads and loads logo from URL
func loadLogoFromURL(url string) (image.Image, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    img, _, err := image.Decode(resp.Body)
    return img, err
}

// loadLogoFromFile loads logo from local file
func loadLogoFromFile(logoPath string) (image.Image, error) {
    // If logo doesn't exist, create a placeholder
    if _, err := os.Stat(logoPath); os.IsNotExist(err) {
        return createPlaceholderLogo(), nil
    }
    
    file, err := os.Open(logoPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    img, _, err := image.Decode(file)
    return img, err
}

// compositeImages combines the template and logo images
func compositeImages(template, logo image.Image, config MockupConfig) (image.Image, error) {
    // Create a new RGBA image based on template bounds
    bounds := template.Bounds()
    composite := image.NewRGBA(bounds)
    
    // Draw the template as background
    draw.Draw(composite, bounds, template, image.Point{}, draw.Src)
    
    // Resize logo if needed
    resizedLogo := resizeLogo(logo, config.LogoMaxSize)
    
    // Calculate logo position (center it at the specified position)
    logoBounds := resizedLogo.Bounds()
    logoRect := image.Rectangle{
        Min: config.LogoPosition.Sub(image.Pt(logoBounds.Dx()/2, logoBounds.Dy()/2)),
        Max: config.LogoPosition.Add(image.Pt(logoBounds.Dx()/2, logoBounds.Dy()/2)),
    }
    
    // Draw logo onto composite
    draw.Draw(composite, logoRect, resizedLogo, logoBounds.Min, draw.Over)
    
    return composite, nil
}

// resizeLogo resizes logo to fit within max dimensions while maintaining aspect ratio
func resizeLogo(logo image.Image, maxSize image.Point) image.Image {
    bounds := logo.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    
    // Calculate scale factor to fit within max dimensions
    scaleX := float64(maxSize.X) / float64(width)
    scaleY := float64(maxSize.Y) / float64(height)
    scale := scaleX
    if scaleY < scaleX {
        scale = scaleY
    }
    
    // If logo is already smaller, don't upscale
    if scale > 1.0 {
        scale = 1.0
    }
    
    newWidth := int(float64(width) * scale)
    newHeight := int(float64(height) * scale)
    
    // Create resized image
    resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
    
    // Simple nearest neighbor scaling
    for y := 0; y < newHeight; y++ {
        for x := 0; x < newWidth; x++ {
            srcX := int(float64(x) / scale)
            srcY := int(float64(y) / scale)
            resized.Set(x, y, logo.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
        }
    }
    
    return resized
}

// saveMockup saves the composite image to file
func saveMockup(img image.Image, outputPath string) error {
    // Ensure directory exists
    dir := filepath.Dir(outputPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    return png.Encode(file, img)
}

// createPlaceholderTemplate creates a simple placeholder template
func createPlaceholderTemplate() image.Image {
    img := image.NewRGBA(image.Rect(0, 0, 500, 600))
    
    // Fill with gray background
    gray := color.RGBA{128, 128, 128, 255}
    for y := 0; y < 600; y++ {
        for x := 0; x < 500; x++ {
            img.Set(x, y, gray)
        }
    }
    
    // Add hoodie shape (simplified)
    hoodie := color.RGBA{64, 64, 64, 255}
    for y := 100; y < 500; y++ {
        for x := 100; x < 400; x++ {
            img.Set(x, y, hoodie)
        }
    }
    
    return img
}

// createPlaceholderLogo creates a simple placeholder logo
func createPlaceholderLogo() image.Image {
    img := image.NewRGBA(image.Rect(0, 0, 100, 100))
    
    // Create a simple colored square logo
    logoColor := color.RGBA{255, 0, 0, 255}
    for y := 20; y < 80; y++ {
        for x := 20; x < 80; x++ {
            img.Set(x, y, logoColor)
        }
    }
    
    return img
}
