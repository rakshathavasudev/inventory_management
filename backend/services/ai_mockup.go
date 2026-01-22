package services

import (
	"bytes"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Hugging Face API configuration (FREE!)
const (
	HuggingFaceAPIURL = "https://api-inference.huggingface.co/models/runwayml/stable-diffusion-v1-5"
	DefaultModel      = "stable-diffusion-v1-5"
)

// AIPromptRequest represents the request for AI-generated mockup
type AIPromptRequest struct {
	Prompt  string `json:"prompt"`
	Product string `json:"product"`
	Color   string `json:"color"`
	Size    string `json:"size"`
}

// HuggingFaceImageRequest represents the request to Hugging Face API
type HuggingFaceImageRequest struct {
	Inputs string `json:"inputs"`
}

// GenerateAIMockup creates a mockup using AI based on a text prompt (FREE with Hugging Face!)
func GenerateAIMockup(orderID uint, request AIPromptRequest) (string, error) {
	fmt.Printf("Starting AI mockup generation for order %d with prompt: %s\n", orderID, request.Prompt)
	
	// Hugging Face API key is required for the new router endpoint
	apiKey := os.Getenv("HUGGINGFACE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("HF_TOKEN") // Also check for standard HF_TOKEN
	}
	fmt.Printf("API Key present: %t\n", apiKey != "")
	
	if apiKey == "" {
		fmt.Println("No API key provided, using fallback")
		return GenerateAIMockupFallback(orderID, request)
	}
	
	// Create a detailed prompt for the AI
	fullPrompt := buildMockupPrompt(request)
	fmt.Printf("Full prompt: %s\n", fullPrompt)

	// Try the new Hugging Face Inference Providers API
	success, mockupPath, err := tryNewHuggingFaceAPI(orderID, fullPrompt, apiKey)
	if success {
		return mockupPath, err
	}
	
	// If new API fails, fall back to HTML mockup
	fmt.Printf("New API failed: %v, using fallback\n", err)
	return GenerateAIMockupFallback(orderID, request)
}

// tryNewHuggingFaceAPI attempts to use the new Hugging Face Inference Providers API
func tryNewHuggingFaceAPI(orderID uint, prompt string, apiKey string) (bool, string, error) {
	// Try the FLUX.2-klein-9B model and other working models using the new router
	models := []struct {
		name     string
		endpoint string
	}{
		{"black-forest-labs/FLUX.2-klein-9B", "https://router.huggingface.co/models/black-forest-labs/FLUX.2-klein-9B"},
		{"ostris/OpenFLUX.1", "https://router.huggingface.co/models/ostris/OpenFLUX.1"},
		{"lodestones/Chroma", "https://router.huggingface.co/models/lodestones/Chroma"},
		{"runwayml/stable-diffusion-v1-5", "https://router.huggingface.co/models/runwayml/stable-diffusion-v1-5"},
	}
	
	for _, model := range models {
		fmt.Printf("Trying model: %s\n", model.name)
		
		success, mockupPath, err := makeHuggingFaceDirectRequest(model.endpoint, prompt, apiKey, orderID)
		if success {
			return true, mockupPath, nil
		}
		fmt.Printf("Model %s failed: %v\n", model.name, err)
	}
	
	return false, "", fmt.Errorf("all models failed")
}

// makeHuggingFaceDirectRequest makes a direct request to Hugging Face Inference API
func makeHuggingFaceDirectRequest(url, prompt, apiKey string, orderID uint) (bool, string, error) {
	// Prepare simple request format for direct API
	requestData := map[string]interface{}{
		"inputs": prompt,
		"parameters": map[string]interface{}{
			"num_inference_steps": 20,
			"guidance_scale":      7.5,
		},
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return false, "", fmt.Errorf("failed to marshal request: %v", err)
	}

	fmt.Printf("Making request to: %s\n", url)
	
	// Make request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return false, "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("API response status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("API error response: %s\n", string(body))
		
		// Handle specific error cases
		if resp.StatusCode == 503 {
			return false, "", fmt.Errorf("model is loading, please try again in a few minutes")
		}
		
		return false, "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Read the image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("failed to read image data: %v", err)
	}

	// Check if we got valid image data
	if len(imageData) < 100 {
		return false, "", fmt.Errorf("received invalid image data (too small)")
	}

	// Save the generated image
	mockupPath, err := saveImageData(imageData, orderID)
	if err != nil {
		return false, "", fmt.Errorf("failed to save image: %v", err)
	}

	fmt.Printf("Successfully generated AI mockup: %s\n", mockupPath)
	return true, mockupPath, nil
}

// buildMockupPrompt creates a detailed prompt for AI image generation
func buildMockupPrompt(request AIPromptRequest) string {
	basePrompt := fmt.Sprintf(
		"A professional product photo of a %s %s in %s color, ",
		request.Color, request.Product, request.Color,
	)

	// Add user's custom prompt
	if request.Prompt != "" {
		basePrompt += request.Prompt + ". "
	}

	// Add style guidelines optimized for Stable Diffusion
	basePrompt += "Professional product photography, clean white background, " +
		"studio lighting, high quality, detailed, realistic, " +
		"e-commerce style, front view, centered composition, " +
		"photorealistic, 4k resolution"

	return basePrompt
}

// saveImageData saves raw image data to a file
func saveImageData(imageData []byte, orderID uint) (string, error) {
	// Create mockups directory if it doesn't exist
	mockupsDir := "mockups"
	if err := os.MkdirAll(mockupsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create mockups directory: %v", err)
	}

	// Save the image
	filename := fmt.Sprintf("ai_mockup_order_%d.png", orderID)
	filepath := filepath.Join(mockupsDir, filename)

	err := os.WriteFile(filepath, imageData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	return "/" + filepath, nil
}

// GenerateAIMockupFallback creates a simple AI-style mockup without external API
// This is used when Hugging Face API is unavailable or rate limited
func GenerateAIMockupFallback(orderID uint, request AIPromptRequest) (string, error) {
	// Create a more sophisticated fallback that simulates AI generation
	mockupsDir := "mockups"
	if err := os.MkdirAll(mockupsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create mockups directory: %v", err)
	}

	filename := fmt.Sprintf("ai_mockup_order_%d.html", orderID)
	filepath := filepath.Join(mockupsDir, filename)

	// Create an HTML mockup that looks like an AI-generated design
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>AI Mockup - Order %d</title>
    <style>
        body { 
            margin: 0; 
            padding: 20px; 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        .mockup-container {
            background: white;
            padding: 40px;
            border-radius: 20px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            text-align: center;
            max-width: 500px;
        }
        .product-outline {
            width: 300px;
            height: 350px;
            margin: 20px auto;
            border: 3px solid #333;
            border-radius: 20px;
            position: relative;
            background: %s;
            display: flex;
            align-items: center;
            justify-content: center;
            box-shadow: inset 0 0 20px rgba(0,0,0,0.1);
        }
        .design-area {
            width: 180px;
            height: 120px;
            background: rgba(255,255,255,0.95);
            border-radius: 15px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 600;
            color: #333;
            text-align: center;
            padding: 15px;
            box-sizing: border-box;
            font-size: 14px;
            border: 2px solid rgba(0,0,0,0.1);
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
        }
        .ai-badge {
            background: linear-gradient(45deg, #ff6b6b, #4ecdc4);
            color: white;
            padding: 10px 20px;
            border-radius: 25px;
            font-size: 14px;
            margin-bottom: 20px;
            display: inline-block;
            font-weight: 600;
        }
        .prompt-display {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 15px;
            margin-top: 20px;
            font-style: italic;
            color: #666;
            border-left: 4px solid #4ecdc4;
        }
        .product-title {
            color: #333;
            margin: 10px 0;
            font-size: 24px;
            font-weight: 600;
        }
        .note {
            background: #e3f2fd;
            color: #1565c0;
            padding: 20px;
            border-radius: 10px;
            margin-top: 15px;
            font-size: 13px;
            border: 1px solid #90caf9;
            line-height: 1.5;
        }
        .setup-steps {
            text-align: left;
            margin-top: 10px;
        }
        .setup-steps ol {
            margin: 10px 0;
            padding-left: 20px;
        }
        .setup-steps li {
            margin: 5px 0;
        }
    </style>
</head>
<body>
    <div class="mockup-container">
        <div class="ai-badge">🤖 AI Mockup Preview</div>
        <h2 class="product-title">%s %s</h2>
        <div class="product-outline">
            <div class="design-area">
                %s
            </div>
        </div>
        <div class="prompt-display">
            <strong>Your Design Vision:</strong><br>
            "%s"
        </div>
        <div class="note">
            <strong>🤖 AI Image Generation Setup Required</strong>
            <div class="setup-steps">
                <p>To generate real AI images, you need a free Hugging Face API key:</p>
                <ol>
                    <li>Go to <strong>https://huggingface.co/settings/tokens</strong></li>
                    <li>Create a free account if you don't have one</li>
                    <li>Generate a new token (select "Read" permissions)</li>
                    <li>Add it to your <code>backend/.env</code> file: <code>HUGGINGFACE_API_KEY=your_token_here</code></li>
                    <li>Restart the backend server</li>
                </ol>
                <p><em>This preview shows what your AI-generated mockup would look like!</em></p>
            </div>
        </div>
    </div>
</body>
</html>`, 
		orderID,
		getColorCode(request.Color),
		request.Color, 
		request.Product,
		request.Prompt,
		request.Prompt,
	)

	err := os.WriteFile(filepath, []byte(htmlContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create fallback mockup: %v", err)
	}

	return "/" + filepath, nil
}

// getColorCode returns a CSS color code for the given color name
func getColorCode(color string) string {
	colorMap := map[string]string{
		"black":  "#2d3748",
		"white":  "#f7fafc",
		"red":    "#e53e3e",
		"blue":   "#3182ce",
		"green":  "#38a169",
		"yellow": "#d69e2e",
		"purple": "#805ad5",
		"pink":   "#d53f8c",
		"gray":   "#718096",
		"navy":   "#2c5282",
	}
	
	if code, exists := colorMap[color]; exists {
		return code
	}
	return "#718096" // default gray
}