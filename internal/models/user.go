package models

// User представляет собой модель пользователя
type User struct {
	ID             string `json:"id,omitempty"`              // Уникальный идентификатор пользователя
	Name           string `json:"name,omitempty"`            // Имя пользователя
	Surname        string `json:"surname,omitempty"`         // Фамилия пользователя
	Patronymic     string `json:"patronymic,omitempty"`      // Отчество пользователя
	Address        string `json:"address,omitempty"`         // Адрес пользователя
	PassportSerie  int    `json:"passport_serie,omitempty"`  // Серия паспорта пользователя
	PassportNumber int    `json:"passport_number,omitempty"` // Номер паспорта пользователя
}
