package storage

import (
	"context"

	"jalabs.kz/bot/storage/postgres"
	"jalabs.kz/bot/storage/postgres/jalabs"
	"jalabs.kz/bot/storage/postgres/todays_jalabs"
	"jalabs.kz/bot/storage/postgres/users"
	"jalabs.kz/bot/storage/postgres/yaxshis"
)

type Storage struct {
	Jalabs       jalabs.Repo
	TodaysJalabs todays_jalabs.Repo
	Users        users.Repo
	Yaxshis      yaxshis.Repo
	Repo         postgres.Repo
}

func New(ctx context.Context) (Storage, error) {
	database, errNewDB := postgres.New()
	if errNewDB != nil {
		return Storage{}, errNewDB
	}

	jalabsRepo, errNewJalabsRepo := jalabs.New(ctx, database)
	if errNewJalabsRepo != nil {
		return Storage{}, errNewJalabsRepo
	}

	todaysJalabsRepo, errNewTodaysJalabsRepo := todays_jalabs.New(ctx, database)
	if errNewTodaysJalabsRepo != nil {
		return Storage{}, errNewTodaysJalabsRepo
	}

	repo, errNewRepo := postgres.NewRepo(ctx, database)
	if errNewRepo != nil {
		return Storage{}, errNewRepo
	}

	usersRepo, errNewUsersRepo := users.New(ctx, database)
	if errNewUsersRepo != nil {
		return Storage{}, errNewUsersRepo
	}

	yaxshisRepo, errNewYaxshisRepo := yaxshis.New(ctx, database)
	if errNewYaxshisRepo != nil {
		return Storage{}, errNewYaxshisRepo
	}

	return Storage{
		Jalabs:       jalabsRepo,
		TodaysJalabs: todaysJalabsRepo,
		Users:        usersRepo,
		Yaxshis:      yaxshisRepo,
		Repo:         repo,
	}, nil
}
