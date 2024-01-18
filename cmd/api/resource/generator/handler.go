package generator

import (
	"fmt"
	"github.com/tnaucoin/pdftmpl/internal/pdfClient"
	"github.com/tnaucoin/pdftmpl/templates"
	"github.com/tnaucoin/pdftmpl/utils"
	"log/slog"
)

func Create(logger *slog.Logger, gotenClient *pdfClient.PdfClient) func(input *GenerateInput) error {
	return func(input *GenerateInput) error {
		imagePath := fmt.Sprintf("./templates/images/%s", input.Body.Logo)
		image, err := utils.GetImageBase64(imagePath)

		if err != nil {
			logger.Error(err.Error())
			return err
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
			return err
		}
		return nil
	}
}
