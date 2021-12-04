package domain

type Transaction struct {
	ID                  int    `json:"id"`
	SenderIIN           string `json:"sender"`
	SenderAccountNumber string `json:"sender_number"`
	RecipientAccNumber  string `json:"recipient_number"`
	RecipientIIN        string `json:"recipient"`
	Date                string `json:"date"`
}
