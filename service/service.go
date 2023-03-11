package service

import (
	"jalabs.kz/bot/storage"
)

type Service struct {
	stg storage.Storage
}

func New(stg storage.Storage) Service {
	return Service{
		stg: stg,
	}
}
