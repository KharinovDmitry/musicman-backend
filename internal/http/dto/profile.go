package dto

type UserProfile struct {
	UUID  string `json:"uuid"`
	Login string `json:"login"`

	SubscribeStatus bool `json:"subscribe_status"`
}
