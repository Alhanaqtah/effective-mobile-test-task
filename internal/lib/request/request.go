package request

// CreateUser содержит данные для создания нового пользователя
type CreateUser struct {
	PassportNumber string `json:"passportNumber,omitempty"` // Номер паспорта пользователя
}
