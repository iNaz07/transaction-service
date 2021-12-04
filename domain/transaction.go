package domain

type Transaction struct {
	ID                  int    `json:"id"`
	SenderIIN           string `json:"sender"`
	SenderAccountNumber string `json:"sender_number"`
	RecipientAccNumber  string `json:"recipient_number"`
	RecipientIIN        string `json:"recipient"`
	Amount              int64  `json:"amount"`
	Date                string `json:"date"`
}

type Deposit struct {
	ID     int    `json:"id"`
	IIN    string `json:"iin"`
	Number string `json:"number"`
	Amount int64  `json:"amount"`
	Date   string `json:"date"`
}
