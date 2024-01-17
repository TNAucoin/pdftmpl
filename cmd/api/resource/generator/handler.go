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
		imagePath := fmt.Sprintf("./templates/images/%s", input.Body.Logo)
		image, err := utils.GetImageBase64(imagePath)

		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		component := templates.Hello(
			input.Body.InvoiceID,
			input.Body.RecipientName,
			input.Body.RecipientAddress,
			input.Body.RecipientEmail,
			image,
		)
		err = gotenClient.GeneratePdfFromComponent(component)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		return &GenerateOutput{}, nil
	}
}
