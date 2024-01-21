package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"github.com/tnaucoin/pdftmpl/cmd/api/resource/types"
	"github.com/tnaucoin/pdftmpl/internal/weasyPrintClient"
	"github.com/tnaucoin/pdftmpl/templates"
	"io"
	"os"
	"strconv"
)

var (
	ErrTemplateNotFound = errors.New("provided input doesn't match any templates")
)

func selectTemplate(i any) (templ.Component, error) {
	switch v := i.(type) {
	case *types.GenerateInput:
		return templates.Hello(
			strconv.FormatUint(v.Body.InvoiceID, 10),
			v.Body.RecipientName,
			v.Body.RecipientAddress,
			v.Body.RecipientEmail,
			"",
		), nil
	default:
		return nil, ErrTemplateNotFound
	}
}

func generatePdfPost(app *Application) func(input *types.GenerateInput) error {
	return func(input *types.GenerateInput) error {
		component, err := selectTemplate(input)
		if err != nil {
			return err
		}
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
