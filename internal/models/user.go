package models

type User struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Surname        string `json:"surname,omitempty"`
	Patronymic     string `json:"patronymic,omitempty"`
	Address        string `json:"address,omitempty"`
	PassportSerie  int    `json:"passport_serie,omitempty"`
	PassportNumber int    `json:"passport_number,omitempty"`
}
