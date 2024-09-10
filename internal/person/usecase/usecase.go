package usecase

import (
	"github.com/Inspirate789/ds-lab1/internal/models"
	"log/slog"
)

type Repository interface { // TODO: Перейти на указатели
	HealthCheck() error
	GetPersons(offset, limit int64) ([]models.Person, error)
	CreatePerson(person models.PersonProperties) (models.Person, error)
	GetPerson(personID int) (models.Person, bool, error)
	UpdatePerson(person models.Person) (models.Person, bool, error)
	DeletePerson(personID int) (bool, error)
}

type UseCase struct {
	repo   Repository
	logger *slog.Logger
}

func New(repo Repository, logger *slog.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (u *UseCase) HealthCheck() error {
	return u.repo.HealthCheck()
}

func (u *UseCase) GetPersons(offset, limit int64) ([]models.Person, error) {
	return u.repo.GetPersons(offset, limit)
}

func (u *UseCase) CreatePerson(person models.PersonProperties) (models.Person, error) {
	return u.repo.CreatePerson(person)
}

func (u *UseCase) GetPerson(personID int) (models.Person, bool, error) {
	return u.repo.GetPerson(personID)
}

func (u *UseCase) UpdatePerson(person models.Person) (models.Person, bool, error) {
	return u.repo.UpdatePerson(person)
}

func (u *UseCase) DeletePerson(personID int) (bool, error) {
	return u.repo.DeletePerson(personID)
}
