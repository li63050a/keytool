package core

import (
	"os"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(content, outPath string, size int) error {
	if size <= 0 {
		size = 512
	}
	return qrcode.WriteFile(content, qrcode.Medium, size, outPath)
}

func GenerateQRCodeString(content string) (string, error) {
	qr, err := qrcode.New(content, qrcode.Low)
	if err != nil {
		return "", err
	}
	return qr.ToSmallString(false), nil
}

func QRCodeFromFile(filePath, outPath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return GenerateQRCode(string(data), outPath, 512)
}
