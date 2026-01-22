package handlers

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "printflow/db"
    "printflow/models"
    "printflow/services"
)

type CreateOrderInput struct {
    Product   string `json:"product"`
    Color     string `json:"color"`
    Size      string `json:"size"`
    LogoURL   string `json:"logoUrl"`
    AIPrompt  string `json:"aiPrompt"`
    UseAI     bool   `json:"useAI"`
}

func CreateOrder(c *gin.Context) {
    var input CreateOrderInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    order := models.Order{
        Product: input.Product,
        Color:   input.Color,
        Size:    input.Size,
        Status:  models.StatusCreated,
    }

    if err := db.DB.Create(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Store the mockup request data for later generation
    if input.UseAI || input.LogoURL != "" {
        asset := models.Asset{
            OrderID:     order.ID,
            LogoURL:     input.LogoURL,
            AIGenerated: input.UseAI,
            AIPrompt:    input.AIPrompt,
            MockupURL:   "", // Will be generated when user clicks "Generate Mockup"
        }
        db.DB.Create(&asset)
    }

    c.JSON(http.StatusCreated, gin.H{
        "order": order,
        "message": "Order created successfully. Click 'Generate Mockup' to create your design.",
    })
}

func ListOrders(c *gin.Context) {
    var orders []models.Order
    db.DB.Find(&orders)
    c.JSON(http.StatusOK, orders)
}

func GetOrder(c *gin.Context) {
    var order models.Order
    if err := db.DB.First(&order, c.Param("ID")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
        return
    }
    
    var asset models.Asset
    // Don't return error if asset doesn't exist, just return empty asset
    db.DB.Where("order_id = ?", order.ID).First(&asset)
    
    c.JSON(http.StatusOK, gin.H{
        "order": order,
        "asset": asset,
    })
}

func ApproveOrder(c *gin.Context) {
    var order models.Order
    if err := db.DB.First(&order, c.Param("ID")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
        return
    }

    if err := services.Transition(&order, models.StatusApproved); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.DB.Save(&order)
    c.JSON(http.StatusOK, order)
}

type MockupInput struct {
	LogoURL  string `json:"logoUrl"`
	AIPrompt string `json:"aiPrompt"`
	UseAI    bool   `json:"useAI"`
}

func GenerateMockupHandler(c *gin.Context) {
	var order models.Order
	if err := db.DB.First(&order, c.Param("ID")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	var input MockupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var mockupURL string
	var err error

	if input.UseAI && input.AIPrompt != "" {
		// Generate AI-powered mockup
		aiRequest := services.AIPromptRequest{
			Prompt:  input.AIPrompt,
			Product: order.Product,
			Color:   order.Color,
			Size:    order.Size,
		}
		
		mockupURL, err = services.GenerateAIMockup(order.ID, aiRequest)
		if err != nil {
			// Fallback to AI fallback mockup if AI fails
			fmt.Printf("AI mockup generation failed: %v, using AI fallback\n", err)
			mockupURL, err = services.GenerateAIMockupFallback(order.ID, aiRequest)
			if err != nil {
				// Final fallback to simple mockup
				fmt.Printf("AI fallback failed: %v, using simple mockup\n", err)
				mockupURL, err = services.GenerateMockupWithProduct(order.ID, input.LogoURL, order.Product, order.Color)
			}
		}
	} else {
		// Generate traditional mockup with logo
		mockupURL, err = services.GenerateMockupWithProduct(order.ID, input.LogoURL, order.Product, order.Color)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate mockup"})
		return
	}

	// Update or create asset
	var asset models.Asset
	result := db.DB.Where("order_id = ?", order.ID).First(&asset)
	
	if result.Error != nil {
		// Create new asset
		asset = models.Asset{
			OrderID:     order.ID,
			LogoURL:     input.LogoURL,
			MockupURL:   mockupURL,
			AIGenerated: input.UseAI,
			AIPrompt:    input.AIPrompt,
		}
		if err := db.DB.Create(&asset).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create asset"})
			return
		}
	} else {
		// Update existing asset
		asset.LogoURL = input.LogoURL
		asset.MockupURL = mockupURL
		asset.AIGenerated = input.UseAI
		asset.AIPrompt = input.AIPrompt
		if err := db.DB.Save(&asset).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update asset"})
			return
		}
	}

	// Update order status
	services.Transition(&order, models.StatusMockupGenerated)
	db.DB.Save(&order)

	c.JSON(http.StatusOK, gin.H{
		"order":  order,
		"asset":  asset,
		"mockup": mockupURL,
	})
}

type LabelInput struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
}

func GenerateLabel(c *gin.Context) {
	var order models.Order
	if err := db.DB.First(&order, c.Param("ID")).Error; err != nil {
		c.JSON(404, gin.H{"error": "order not found"})
		return
	}

	var input LabelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	// Convert handler LabelInput to services LabelInput
	serviceInput := services.LabelInput{
		Name:    input.Name,
		Address: input.Address,
		City:    input.City,
		State:   input.State,
		Zip:     input.Zip,
	}

	url, err := services.GenerateLabel(order.ID, serviceInput)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("label generation failed: %v", err)})
		return
	}

	c.JSON(200, gin.H{"label": url})
}

// GetAvailableColors returns list of available product colors
func GetAvailableColors(c *gin.Context) {
	colors := services.GetAvailableColors()
	c.JSON(http.StatusOK, gin.H{
		"colors": colors,
	})
}