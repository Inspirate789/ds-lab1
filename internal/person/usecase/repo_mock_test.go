package usecase_test

import (
	"context"
	"github.com/Inspirate789/ds-lab1/internal/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryPositiveMock struct {
	mock.Mock
}

func (r *RepositoryPositiveMock) HealthCheck(_ context.Context) error {
	args := r.Called()
	return args.Error(0)
}

func (r *RepositoryPositiveMock) GetPersons(_ context.Context, offset, limit int64) ([]models.Person, error) {
	args := r.Called(offset, limit)
	return args.Get(0).([]models.Person), args.Error(1)
}

func (r *RepositoryPositiveMock) CreatePerson(_ context.Context, person models.PersonProperties) (models.Person, error) {
	args := r.Called(person)
	return args.Get(0).(models.Person), args.Error(1)
}

func (r *RepositoryPositiveMock) GetPerson(_ context.Context, personID int) (models.Person, bool, error) {
	args := r.Called(personID)
	return args.Get(0).(models.Person), args.Bool(1), args.Error(2)
}

func (r *RepositoryPositiveMock) UpdatePerson(_ context.Context, person models.Person) (models.Person, bool, error) {
	args := r.Called(person)
	return args.Get(0).(models.Person), args.Bool(1), args.Error(2)
}

func (r *RepositoryPositiveMock) DeletePerson(_ context.Context, personID int) (bool, error) {
	args := r.Called(personID)
	return args.Bool(0), args.Error(1)
}
