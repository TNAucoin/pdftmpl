package utils

import (
	"encoding/base64"
	"io/ioutil"
)

func GetImageBase64(path string) (string, error) {
	img, err := ioutil.ReadFile(path)
	if err != nil {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(img), nil
}
