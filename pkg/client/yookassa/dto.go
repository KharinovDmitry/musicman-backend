package yookassa

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
	ID                   string                 `json:"id,omitempty"`
	Status               string                 `json:"status,omitempty"`
	Amount               Amount                 `json:"amount,omitempty"`
	IncomeAmount         Amount                 `json:"income_amount,omitempty"`
	Description          string                 `json:"description,omitempty" binding:"max=128"`
	Capture              bool                   `json:"capture,omitempty"`
	Recipient            *Recipient             `json:"recipient,omitempty"`
	PaymentMethod        map[string]interface{} `json:"payment_method,omitempty"`
	CapturedAt           *time.Time             `json:"captured_at,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	ExpiresAt            *time.Time             `json:"expires_at,omitempty"`
	Confirmation         map[string]interface{} `json:"confirmation,omitempty"`
	Test                 bool                   `json:"test,omitempty"`
	RefundedAmount       *Amount                `json:"refunded_amount,omitempty"`
	Paid                 bool                   `json:"paid,omitempty"`
	Refundable           bool                   `json:"refundable,omitempty"`
	ReceiptRegistration  string                 `json:"receipt_registration,omitempty"`
	Metadata             interface{}            `json:"metadata,omitempty"`
	CancellationDetails  *CancellationDetails   `json:"cancellation_details,omitempty"`
	AuthorizationDetails *AuthorizationDetails  `json:"authorization_details,omitempty"`
	MerchantCustomerID   string                 `json:"merchant_customer_id,omitempty" binding:"max=200"`
}

type CancellationDetails struct {
	Party  string `json:"party,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type AuthorizationDetails struct {
	RRN          string `json:"rrn,omitempty"`
	AuthCode     string `json:"auth_code,omitempty"`
	ThreeDSecure struct {
		Applied bool `json:"applied,omitempty"`
	} `json:"three_d_secure,omitempty"`
}
