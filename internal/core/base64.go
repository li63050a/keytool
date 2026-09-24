package core

import (
	"encoding/base64"
	"os"
)

func Base64EncodeFile(inPath, outPath string) error {
	data, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return os.WriteFile(outPath, []byte(encoded), 0644)
}

func Base64DecodeFile(inPath, outPath string) error {
	data, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, decoded, 0644)
}

func Base64EncodeString(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Base64DecodeString(s string) (string, error) {
	d, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
