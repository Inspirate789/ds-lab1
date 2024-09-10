package usecase

import (
	"context"
	"github.com/Inspirate789/ds-lab1/internal/models"
	"log/slog"
)

type Repository interface {
	HealthCheck(ctx context.Context) error
	GetPersons(ctx context.Context, offset, limit int64) ([]models.Person, error)
	CreatePerson(ctx context.Context, person models.PersonProperties) (models.Person, error)
	GetPerson(ctx context.Context, personID int) (models.Person, bool, error)
	UpdatePerson(ctx context.Context, person models.Person) (models.Person, bool, error)
	DeletePerson(ctx context.Context, personID int) (bool, error)
}

type UseCase struct {
	repo   Repository
	logger *slog.Logger
}

func New(repo Repository, logger *slog.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (u *UseCase) HealthCheck(ctx context.Context) error {
	return u.repo.HealthCheck(ctx)
}

func (u *UseCase) GetPersons(ctx context.Context, offset, limit int64) ([]models.Person, error) {
	return u.repo.GetPersons(ctx, offset, limit)
}

func (u *UseCase) CreatePerson(ctx context.Context, person models.PersonProperties) (models.Person, error) {
	return u.repo.CreatePerson(ctx, person)
}

func (u *UseCase) GetPerson(ctx context.Context, personID int) (models.Person, bool, error) {
	return u.repo.GetPerson(ctx, personID)
}

func (u *UseCase) UpdatePerson(ctx context.Context, person models.Person) (models.Person, bool, error) {
	return u.repo.UpdatePerson(ctx, person)
}

func (u *UseCase) DeletePerson(ctx context.Context, personID int) (bool, error) {
	return u.repo.DeletePerson(ctx, personID)
}
