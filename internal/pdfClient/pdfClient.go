package pdfClient

import (
	"bytes"
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/tnaucoin/pdftmpl/config"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type PdfClient struct {
	baseURI    string
	config     config.ConfGoten
	httpClient *http.Client
}

func New(gotenConfig config.ConfGoten) *PdfClient {
	c := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}
	baseUri := fmt.Sprintf("http://%s:%d", gotenConfig.Host, gotenConfig.Port)

	return &PdfClient{
		baseURI:    baseUri,
		config:     gotenConfig,
		httpClient: c,
	}
}

func (p *PdfClient) renderDocument(component templ.Component) (*bytes.Buffer, error) {
	var b bytes.Buffer
	err := component.Render(context.Background(), &b)
	return &b, err
}

func (p *PdfClient) createFormData(buffer *bytes.Buffer, fieldName, fileName string) (bytes.Buffer, string, error) {
	var formData bytes.Buffer
	multipartWriter := multipart.NewWriter(&formData)
	formFile, err := multipartWriter.CreateFormFile(fieldName, fileName)
	if err != nil {
		return bytes.Buffer{}, "", err
	}
	_, err = io.Copy(formFile, buffer)
	if err != nil {
		return bytes.Buffer{}, "", err
	}
	// Important: Close the multipart writer so it writes the ending boundary.
	err = multipartWriter.Close()
	if err != nil {
		return bytes.Buffer{}, "", err
	}
	return formData, multipartWriter.FormDataContentType(), nil
}

func (p *PdfClient) convertHTMLToPDF(formData bytes.Buffer, contentType string) error {
	endpoint := fmt.Sprintf("%s/forms/chromium/convert/html", p.baseURI)
	req, err := http.NewRequest(http.MethodPost, endpoint, &formData)
	if err != nil {
		log.Fatal(err)
	}
	// Set the Content Type header to the MediaType reported by the multipart Writer
	req.Header.Set("Content-Type", contentType)

	// Send off the request
	resp, err := p.httpClient.Do(req)
	// Always ensure response body gets closed
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code: %d\n", resp.StatusCode)
	} else {
		filePath := fmt.Sprintf("%s/test.pdf", p.config.VolumeOutPath)
		out, err := os.Create(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PdfClient) GeneratePdfFromComponent(component templ.Component) error {
	componentBytes, err := p.renderDocument(component)
	if err != nil {
		return err
	}
	form, contentType, err := p.createFormData(componentBytes, "index.html", "index.html")
	if err != nil {
		return err
	}
	err = p.convertHTMLToPDF(form, contentType)
	if err != nil {
		return err
	}
	return nil
}
