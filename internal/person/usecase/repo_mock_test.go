package usecase_test

import (
	"github.com/Inspirate789/ds-lab1/internal/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryPositiveMock struct {
	mock.Mock
}

func (r *RepositoryPositiveMock) HealthCheck() error {
	args := r.Called()
	return args.Error(0)
}

func (r *RepositoryPositiveMock) GetPersons(offset, limit int64) ([]models.Person, error) {
	args := r.Called(offset, limit)
	return args.Get(0).([]models.Person), args.Error(1)
}

func (r *RepositoryPositiveMock) CreatePerson(person models.PersonProperties) (models.Person, error) {
	args := r.Called(person)
	return args.Get(0).(models.Person), args.Error(1)
}

func (r *RepositoryPositiveMock) GetPerson(personID int) (models.Person, bool, error) {
	args := r.Called(personID)
	return args.Get(0).(models.Person), args.Bool(1), args.Error(2)
}

func (r *RepositoryPositiveMock) UpdatePerson(person models.Person) (models.Person, bool, error) {
	args := r.Called(person)
	return args.Get(0).(models.Person), args.Bool(1), args.Error(2)
}

func (r *RepositoryPositiveMock) DeletePerson(personID int) (bool, error) {
	args := r.Called(personID)
	return args.Bool(0), args.Error(1)
}
