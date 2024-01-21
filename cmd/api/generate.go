package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/tnaucoin/pdftmpl/cmd/api/resource/types"
	"github.com/tnaucoin/pdftmpl/internal/weasyPrintClient"
	"github.com/tnaucoin/pdftmpl/templates"
	"github.com/tnaucoin/pdftmpl/utils"
	"io"
	"os"
	"strconv"
)

func generatePdfPost(app *Application) func(input *types.GenerateInput) error {
	return func(input *types.GenerateInput) error {
		imagePath := fmt.Sprintf("./templates/images/%s", input.Body.Logo)
		image, err := utils.GetImageBase64(imagePath)

		if err != nil {
			app.logger.Error(err.Error())
			return err
		}
		component := templates.Hello(
			strconv.FormatUint(input.Body.InvoiceID, 10),
			input.Body.RecipientName,
			input.Body.RecipientAddress,
			input.Body.RecipientEmail,
			image,
		)
		var doc bytes.Buffer
		err = component.Render(context.Background(), &doc)
		if err != nil {
			app.logger.Error(err.Error())
			return err
		}
		var pdf bytes.Buffer
		err = app.weasyPrint.GeneratePDF(&doc, &pdf)
		if err != nil {
			if errors.Is(err, weasyPrintClient.ErrWeasyPrintGeneration) {
				return err
			}
			// Log errors not caused by weasyPrintClient
			app.logger.Error(err.Error())
			return err
		}
		filePath := fmt.Sprintf("%s/test.pdf", app.config.Server.VolumeOutPath)
		out, err := os.Create(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(out, &pdf)
		if err != nil {
			return err
		}
		return nil
	}
}
