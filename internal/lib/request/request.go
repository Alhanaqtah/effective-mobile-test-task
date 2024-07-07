package request

type CreateUser struct {
	PassportNumber string `json:"passportNumber,omitempty"`
}
