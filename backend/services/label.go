package services

import (
	"fmt"
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
	// Create directories
	if err := os.MkdirAll("labels/tmp", 0755); err != nil {
		return "", err
	}

	// Generate barcode
	code, err := code128.Encode(fmt.Sprintf("CRONK-%d", orderID))
	if err != nil {
		return "", err
	}

	scaled, err := barcode.Scale(code, 400, 120)
	if err != nil {
		return "", err
	}

	barcodePath := fmt.Sprintf("labels/tmp/%d.png", orderID)

	// Write barcode file
	f, err := os.Create(barcodePath)
	if err != nil {
		return "", err
	}
	if err := png.Encode(f, scaled); err != nil {
		f.Close()
		return "", err
	}
	f.Close()

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
	pdf.Cell(0, 12, "CRONK SHIPPING LABEL")
	pdf.Ln(14)

	pdf.SetFont("Helvetica", "", 12)
	pdf.MultiCell(0, 8, fmt.Sprintf(
		"TO:\n%s\n%s\n%s, %s %s\n\nOrder #%d",
		input.Name,
		input.Address,
		input.City,
		input.State,
		input.Zip,
		orderID,
	), "", "", false)

	pdf.Ln(10)

	pdf.Image(absBarcodePath, 10, 120, 150, 0, false, "", 0, "")

	out := fmt.Sprintf("labels/order_%d.pdf", orderID)
	if err := pdf.OutputFileAndClose(out); err != nil {
		return "", err
	}

	return "/" + out, nil
}
