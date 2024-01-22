package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/a-h/templ"
	"io"
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

func RenderTemplate(input templ.Component) (*bytes.Buffer, error) {
	var out bytes.Buffer
	err := input.Render(context.Background(), &out)
	return &out, err
}

func WriteToFile(doc *bytes.Buffer, outPath, fileName string) error {
	filePath := fmt.Sprintf("%s/%s", outPath, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, doc)
	if err != nil {
		return err
	}
	return nil
}
