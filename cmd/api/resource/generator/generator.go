package generator

type GenerateInput struct {
	Body struct {
		RecipientName    string `json:"recipient-name,required" maxLength:"50" doc:"Name of invoice recipient"`
		RecipientAddress string `json:"recipient-address,required" doc:"Address of invoice recipient"`
		RecipientEmail   string `json:"recipient-email,required" doc:"Email of invoice recipient" pattern:"^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"`
		Logo             string `json:"logo,required" enum:"logo.png,logo2.png" doc:"Logo to include on invoice"`
	}
}

type GenerateOutput struct {
}
