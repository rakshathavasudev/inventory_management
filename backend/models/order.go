package models

import "time"

const (
    StatusCreated         = "CREATED"
    StatusMockupGenerated = "MOCKUP_GENERATED"
    StatusApproved        = "APPROVED"
    StatusReady           = "READY_FOR_FULFILLMENT"
)

type Order struct {
    ID        uint      `gorm:"primaryKey"`
    Product   string
    Color     string
    Size      string
    Status    string
    CreatedAt time.Time
}

type Asset struct {
    ID        uint   `gorm:"primaryKey"`
    OrderID   uint
    LogoURL   string
    MockupURL string
    AIGenerated bool `gorm:"default:false"`
    AIPrompt    string
}
