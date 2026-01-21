package handlers

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "minicronk/db"
    "minicronk/models"
    "minicronk/services"
)

type CreateOrderInput struct {
    Product string `json:"product"`
    Color   string `json:"color"`
    Size    string `json:"size"`
    LogoURL string `json:"logoUrl"`
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

    mockupURL, err := services.GenerateMockup(order.ID, input.LogoURL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate mockup"})
        return
    }

    asset := models.Asset{
        OrderID:   order.ID,
        LogoURL:   input.LogoURL,
        MockupURL: mockupURL,
    }
    db.DB.Create(&asset)

    services.Transition(&order, models.StatusMockupGenerated)
    db.DB.Save(&order)

    c.JSON(http.StatusCreated, gin.H{
        "order":  order,
        "mockup": mockupURL,
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
    c.JSON(http.StatusOK, order)
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
