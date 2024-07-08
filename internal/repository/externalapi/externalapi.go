package externalapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time-tracker/internal/models"
	"time-tracker/internal/repository"
)

type PeopleInfoRepo struct {
	address string
	client  *http.Client
}

func New(address string) *PeopleInfoRepo {
	return &PeopleInfoRepo{
		address: address,
		client:  &http.Client{},
	}
}

func (p *PeopleInfoRepo) GetUserInfo(passportSerie, passportNumber int) (*models.User, error) {
	const op = "repository.externalapi.GetUserInfo"

	url := fmt.Sprintf("http://%s/info?passportSerie=%d&passportNumber=%d", p.address, passportSerie, passportNumber)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode != http.StatusBadRequest {
			return nil, fmt.Errorf("%s: %d", op, repository.ErrBadRequest)
		}
		if resp.StatusCode != http.StatusInternalServerError {
			return nil, fmt.Errorf("%s: %d", op, repository.ErrExternalAPIError)
		}
	}

	// Декодирование JSON-ответа
	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
