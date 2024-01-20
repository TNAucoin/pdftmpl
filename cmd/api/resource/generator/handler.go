package generator

import (
	"bytes"
	"fmt"
	"github.com/tnaucoin/pdftmpl/internal/pdfClient"
	"github.com/tnaucoin/pdftmpl/templates"
	"github.com/tnaucoin/pdftmpl/utils"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"time"
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
			strconv.FormatUint(input.Body.InvoiceID, 10),
			input.Body.RecipientName,
			input.Body.RecipientAddress,
			input.Body.RecipientEmail,
			image,
		)
		startTime := time.Now()
		document, err := gotenClient.RenderDocument(component)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		c := exec.Command("weasyprint", "-", "-")
		c.Stdin = document
		var bo bytes.Buffer
		c.Stdout = &bo
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		filePath := fmt.Sprintf("%s/test.pdf", "/var/containerOut")
		out, err := os.Create(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(out, &bo)
		if err != nil {
			return err
		}
		logger.Debug("weasyprint generation", "duration", time.Since(startTime))
		return nil
	}
}
