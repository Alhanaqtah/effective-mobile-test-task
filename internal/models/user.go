package models

type User struct {
	Name           string `json:"name,omitempty"`
	Surname        string `json:"surname,omitempty"`
	Patronymic     string `json:"patronymic,omitempty"`
	Address        string `json:"address,omitempty"`
	PassportSerie  int64  `json:"passport_serie,omitempty"`
	PassportNumber int64  `json:"passport_number,omitempty"`
}
