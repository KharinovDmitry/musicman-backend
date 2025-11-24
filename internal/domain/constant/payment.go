package constant

const (
	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusCanceled  = "canceled"

	PaymentStatusPendingRU   = "В обработке"
	PaymentStatusSucceededRU = "Завершен"
	PaymentStatusCanceledRU  = "Отменен"
)

var (
	PaymentStatusTranslate = map[string]string{
		PaymentStatusPending:   PaymentStatusPendingRU,
		PaymentStatusSucceeded: PaymentStatusSucceededRU,
		PaymentStatusCanceled:  PaymentStatusCanceledRU,
	}
)
