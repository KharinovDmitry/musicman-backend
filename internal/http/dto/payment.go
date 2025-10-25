package dto

import (
	"github.com/musicman-backend/internal/domain/entity"
	"time"
)

type CreatePaymentRequest struct {
	// Amount сумма платежа в КОПЕЙКАХ
	Amount int `json:"amount" validate:"required"`
	// ReturnURI ссылка на которую вернуть после оплаты
	ReturnURI string `json:"return_uri" validate:"required"`
}

type UserPayment struct {
	ID            string    `json:"id"`
	PaymentStatus string    `json:"payment_status"`
	Description   string    `json:"description"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

func UserPaymentsFromEntities(payments []entity.Payment) []UserPayment {
	res := make([]UserPayment, len(payments))
	for i, payment := range payments {
		res[i] = UserPayment{
			ID:            payment.ID,
			PaymentStatus: string(payment.PaymentStatus),
			Description:   payment.Description,
			Amount:        payment.Amount,
			CreatedAt:     payment.CreatedAt,
		}
	}

	return res
}
