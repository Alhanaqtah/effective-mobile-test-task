package externalapi

import (
	"net/http"
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
