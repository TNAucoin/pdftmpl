package utils

import (
	"encoding/base64"
	"fmt"
	"os"
)

func GetImageBase64(path string) (string, error) {
	img, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(img), nil
}
