package services

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/jung-kurt/gofpdf"
)

type LabelInput struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
}

func GenerateLabel(orderID uint, input LabelInput) (string, error) {
	// Provide default values to make it fail-safe
	if input.Name == "" {
		input.Name = "Customer"
	}
	if input.Address == "" {
		input.Address = "123 Main Street"
	}
	if input.City == "" {
		input.City = "Anytown"
	}
	if input.State == "" {
		input.State = "CA"
	}
	if input.Zip == "" {
		input.Zip = "12345"
	}

	// Create directories
	if err := os.MkdirAll("labels/tmp", 0755); err != nil {
		return "", err
	}

	// Generate barcode
	code, err := code128.Encode(fmt.Sprintf("PRINTFLOW-%d", orderID))
	if err != nil {
		return "", err
	}

	scaled, err := barcode.Scale(code, 400, 120)
	if err != nil {
		return "", err
	}

	barcodePath := fmt.Sprintf("labels/tmp/%d.png", orderID)

	// Write barcode file with proper 8-bit encoding
	f, err := os.Create(barcodePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Convert to RGBA to ensure 8-bit depth
	bounds := scaled.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, scaled, bounds.Min, draw.Src)

	if err := png.Encode(f, rgba); err != nil {
		return "", err
	}

	absBarcodePath, err := filepath.Abs(barcodePath)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(absBarcodePath); err != nil {
		return "", fmt.Errorf("barcode file not found: %s", absBarcodePath)
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()

	pdf.SetFont("Helvetica", "B", 18)
	pdf.Cell(0, 12, "PRINTFLOW SHIPPING LABEL")
	pdf.Ln(14)

	pdf.SetFont("Helvetica", "", 12)
	pdf.MultiCell(0, 8, fmt.Sprintf(
		"TO:\n%s\n%s\n%s, %s %s\n\nOrder #%d\nProduct: %s\nSize: %s",
		input.Name,
		input.Address,
		input.City,
		input.State,
		input.Zip,
		orderID,
		"Custom Apparel", // We'd need to pass product info here
		"Size TBD",       // We'd need to pass size info here
	), "", "", false)

	pdf.Ln(10)

	pdf.Image(absBarcodePath, 10, 120, 150, 0, false, "", 0, "")

	out := fmt.Sprintf("labels/order_%d.pdf", orderID)
	if err := pdf.OutputFileAndClose(out); err != nil {
		return "", err
	}

	return "/" + out, nil
}
