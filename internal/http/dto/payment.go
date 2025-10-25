package dto

type CreatePaymentRequest struct {
	// Amount сумма платежа в КОПЕЙКАХ
	Amount int `json:"amount"`
	// ReturnURI ссылка на которую вернуть после оплаты
	ReturnURI string `json:"return_uri"`
}
