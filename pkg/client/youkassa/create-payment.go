package youkassa

import "fmt"

func (c *Client) CreatePayment(request CreatePaymentRequest) (CreatePaymentResponse, error) {
	const path = "/v3/payments"

	var response CreatePaymentResponse
	resp, err := c.client.R().
		SetBody(request).
		SetBasicAuth(c.accountID, c.secretKey).
		SetResult(&CreatePaymentResponse{}).
		Post(c.host.JoinPath(path).String())
	if err != nil {
		return CreatePaymentResponse{}, fmt.Errorf("error creating payment: %w", err)
	}

	if resp.IsError() {
		return CreatePaymentResponse{}, fmt.Errorf("error creating payment: %s, bode: %s", resp.Status(), resp.Body())
	}

	return response, nil
}

func (c *Client) GetPayment(id string) (PaymentByIDResponse, error) {
	path := fmt.Sprintf("/v3/payments/%s", id)

	var response PaymentByIDResponse

	resp, err := c.client.R().
		SetBasicAuth(c.accountID, c.secretKey).
		SetResult(&response).
		Get(c.host.JoinPath(path).String())
	if err != nil {
		return PaymentByIDResponse{}, fmt.Errorf("error getting payment: %w", err)
	}

	if resp.IsError() {
		return PaymentByIDResponse{}, fmt.Errorf("error getting payment: %s, bode: %s", resp.Status(), resp.Body())
	}

	return response, nil
}
