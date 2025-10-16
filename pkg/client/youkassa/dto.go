package youkassa

import "time"

type CreatePaymentRequest struct {
	Amount             Amount            `json:"amount"`
	Description        string            `json:"description"`
	Test               bool              `json:"test"`
	Confirmation       Confirmation      `json:"confirmation"`
	MerchantCustomerID string            `json:"merchant_customer_id"`
	Capture            bool              `json:"capture"`
	Metadata           map[string]string `json:"metadata"`
}

type CreatePaymentResponse struct {
	ID           string         `json:"id"`
	Status       string         `json:"status"`
	Paid         bool           `json:"paid"`
	Amount       Amount         `json:"amount"`
	Confirmation Confirmation   `json:"confirmation"`
	CreatedAt    time.Time      `json:"created_at"`
	Description  string         `json:"description"`
	Metadata     map[string]any `json:"metadata"`
	Recipient    Recipient      `json:"recipient"`
	Refundable   bool           `json:"refundable"`
	Test         bool           `json:"test"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type            string `json:"type"`
	ConfirmationURL string `json:"confirmation_url"`
}

type Recipient struct {
	AccountID string `json:"account_id"`
	GatewayID string `json:"gateway_id"`
}

type PaymentByIDResponse struct {
}
