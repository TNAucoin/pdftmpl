package generator

type GenerateInput struct {
	Body struct {
		InvoiceID        uint64 `json:"invoice-id,required" doc:"Invoice ID number" example:"1024"`
		RecipientName    string `json:"recipient-name,required" maxLength:"50" doc:"Name of invoice recipient" example:"John Doe"`
		RecipientAddress string `json:"recipient-address,required" doc:"Address of invoice recipient" example:"123 Somewhere Ln."`
		RecipientEmail   string `json:"recipient-email,required" doc:"Email of invoice recipient" example:"example@gmail.com" pattern:"^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"`
		Logo             string `json:"logo,required" enum:"logo.png,logo2.png" doc:"Logo to include on invoice"`
	}
}

type GenerateOutput struct {
}
