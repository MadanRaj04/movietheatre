package main

import (
	"os"
	"path/filepath"

	qrcode "github.com/skip2/go-qrcode"
)

func GenerateQRCodePNG(content, filename string) (string, error) {
	if err := os.MkdirAll("uploads/qrcodes", 0o755); err != nil {
		return "", err
	}
	path := filepath.Join("uploads", "qrcodes", filename)
	if err := qrcode.WriteFile(content, qrcode.Medium, 256, path); err != nil {
		return "", err
	}
	return path, nil
}
