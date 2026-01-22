package services

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    _ "image/jpeg" // Add JPEG support
    "os"
    "path/filepath"
    "net/http"
    "strings"
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
            LogoPosition: image.Pt(240, 210), // Chest area (moved right from 220 to 240, down from 200 to 210)
            LogoMaxSize:  image.Pt(150, 150), // Max logo dimensions
            OutputDir:    "mockups",
        },
        "T-Shirt": {
            LogoPosition: image.Pt(220, 190), // Chest area (moved right from 200 to 220, down from 180 to 190)
            LogoMaxSize:  image.Pt(120, 120),
            OutputDir:    "mockups",
        },
    }
    
    if config, exists := configs[product]; exists {
        return config
    }
    
    // Default config
    return MockupConfig{
        LogoPosition: image.Pt(180, 200), // Left chest default
        LogoMaxSize:  image.Pt(150, 150),
        OutputDir:    "mockups",
    }
}

// GenerateMockupWithProduct creates a product mockup by compositing a logo onto a product template
func GenerateMockupWithProduct(orderID uint, logoURL string, product string, productColor string) (string, error) {
    if logoURL == "" {
        return "", fmt.Errorf("no logo URL provided")
    }

    // Get product-specific configuration
    config := GetMockupConfig(product)
    
    // Load the product template using actual uploaded images
    var templatePath string
    if strings.ToLower(product) == "hoodie" {
        templatePath = "assets/hoodie.png"
    } else if strings.ToLower(product) == "t-shirt" {
        templatePath = "assets/tshirt.jpg"
    } else {
        // Default to hoodie if unknown product
        templatePath = "assets/hoodie.png"
    }
    
    template, err := loadTemplate(templatePath)
    if err != nil {
        return "", fmt.Errorf("failed to load template: %v", err)
    }
    
    // Apply color to template if specified
    if productColor != "" && strings.ToLower(productColor) != "default" {
        template = applyColorToTemplate(template, productColor)
    }
    
    // Load the logo
    logoPath := "." + logoURL
    logo, err := loadLogo(logoPath)
    if err != nil {
        return "", fmt.Errorf("failed to load logo: %v", err)
    }
    
    // Composite the logo onto the template
    composite, err := compositeImages(template, logo, config)
    if err != nil {
        return "", fmt.Errorf("failed to composite images: %v", err)
    }
    
    // Save the final mockup
    outputPath := fmt.Sprintf("mockups/order_%d.png", orderID)
    if err := saveMockup(composite, outputPath); err != nil {
        return "", fmt.Errorf("failed to save mockup: %v", err)
    }

    return "/" + outputPath, nil
}

// GenerateMockup creates a product mockup (backward compatibility)
func GenerateMockup(orderID uint, logoURL string) (string, error) {
    return GenerateMockupWithProduct(orderID, logoURL, "Hoodie", "gray")
}


// loadTemplate loads the product template image
func loadTemplate(templatePath string) (image.Image, error) {
    file, err := os.Open(templatePath)
    if err != nil {
        return nil, fmt.Errorf("template file not found: %s", templatePath)
    }
    defer file.Close()
    
    img, _, err := image.Decode(file)
    if err != nil {
        return nil, fmt.Errorf("failed to decode template image: %v", err)
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


// ColorMap defines available colors for products
var ColorMap = map[string]color.RGBA{
	"black":   {0, 0, 0, 255},
	"white":   {255, 255, 255, 255},
	"red":     {220, 20, 60, 255},
	"blue":    {30, 144, 255, 255},
	"green":   {34, 139, 34, 255},
	"yellow":  {255, 215, 0, 255},
	"purple":  {128, 0, 128, 255},
	"orange":  {255, 165, 0, 255},
	"pink":    {255, 192, 203, 255},
	"gray":    {128, 128, 128, 255},
	"navy":    {0, 0, 128, 255},
	"maroon":  {128, 0, 0, 255},
}

// applyColorToTemplate changes the color of the product template
func applyColorToTemplate(template image.Image, colorName string) image.Image {
	bounds := template.Bounds()
	colored := image.NewRGBA(bounds)
	
	// Get the target color
	targetColor, exists := ColorMap[strings.ToLower(colorName)]
	if !exists {
		// Default to original template if color not found
		draw.Draw(colored, bounds, template, bounds.Min, draw.Src)
		return colored
	}
	
	// Process each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := template.At(x, y)
			r, g, b, a := originalColor.RGBA()
			
			// Convert to 8-bit values
			r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
			
			// Check if this pixel should be recolored (not transparent and not too bright/dark)
			if a8 > 128 && isColorablePixel(r8, g8, b8) {
				// Apply color while preserving brightness/shadows
				newColor := blendColorWithShading(targetColor, r8, g8, b8, a8)
				colored.Set(x, y, newColor)
			} else {
				// Keep original pixel (transparent areas, highlights, shadows)
				colored.Set(x, y, originalColor)
			}
		}
	}
	
	return colored
}

// isColorablePixel determines if a pixel should be recolored
func isColorablePixel(r, g, b uint8) bool {
	// Don't recolor very dark (shadows) or very bright (highlights) pixels
	brightness := (int(r) + int(g) + int(b)) / 3
	return brightness > 30 && brightness < 200
}

// blendColorWithShading applies the target color while preserving shading
func blendColorWithShading(targetColor color.RGBA, r, g, b, a uint8) color.RGBA {
	// Calculate the brightness factor from the original pixel
	originalBrightness := (int(r) + int(g) + int(b)) / 3
	brightnessFactor := float64(originalBrightness) / 128.0 // Normalize around middle gray
	
	// Apply brightness factor to target color
	newR := uint8(float64(targetColor.R) * brightnessFactor)
	newG := uint8(float64(targetColor.G) * brightnessFactor)
	newB := uint8(float64(targetColor.B) * brightnessFactor)
	
	// Ensure values don't exceed 255
	if newR > 255 { newR = 255 }
	if newG > 255 { newG = 255 }
	if newB > 255 { newB = 255 }
	
	return color.RGBA{newR, newG, newB, a}
}

// GetAvailableColors returns list of available colors
func GetAvailableColors() []string {
	colors := make([]string, 0, len(ColorMap))
	for colorName := range ColorMap {
		colors = append(colors, colorName)
	}
	return colors
}