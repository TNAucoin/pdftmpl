package generator

import (
	"fmt"
	"github.com/tnaucoin/pdftmpl/internal/pdfClient"
	"github.com/tnaucoin/pdftmpl/templates"
	"github.com/tnaucoin/pdftmpl/utils"
	"log/slog"
)

func Create(logger *slog.Logger, gotenClient *pdfClient.PdfClient) func(input *GenerateInput) (*GenerateOutput, error) {
	return func(input *GenerateInput) (*GenerateOutput, error) {
		name := input.Body.RecipientName
		addr := input.Body.RecipientAddress
		email := input.Body.RecipientEmail
		imageName := input.Body.Logo
		imagePath := fmt.Sprintf("./templates/images/%s", imageName)
		image, err := utils.GetImageBase64(imagePath)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		component := templates.Hello(name, addr, email, image)
		err = gotenClient.GeneratePdfFromComponent(component)
		if err != nil {
			return nil, err
		}
		return &GenerateOutput{}, nil
	}
}
