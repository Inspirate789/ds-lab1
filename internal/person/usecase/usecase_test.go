package usecase_test

import (
	"github.com/Inspirate789/ds-lab1/internal/models"
	"github.com/Inspirate789/ds-lab1/internal/person/usecase"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"log/slog"
	"os"
	"strconv"
	"testing"
)

type UseCaseSuite struct {
	suite.Suite
}

func (s *UseCaseSuite) TestHealthCheck(t provider.T) {
	t.Epic("MVP")
	t.Severity(allure.NORMAL)

	// arrange
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := new(RepositoryPositiveMock)
	repo.On("HealthCheck").Return(nil)
	useCase := usecase.New(repo, logger)
	// act
	err := useCase.HealthCheck()
	// assert
	t.Require().NoError(err)
	repo.AssertExpectations(t)
	repo.AssertNumberOfCalls(t, "HealthCheck", 1)
}

func (*UseCaseSuite) newPerson(id int) models.Person {
	return models.Person{
		ID: id,
		PersonProperties: models.PersonProperties{
			Name:    "Aboba " + strconv.Itoa(id),
			Age:     id,
			Address: "Address " + strconv.Itoa(id),
			Work:    "Work " + strconv.Itoa(id),
		},
	}
}

func (s *UseCaseSuite) TestGetPersons(t provider.T) {
	t.Epic("MVP")
	t.Severity(allure.NORMAL)

	// arrange
	const offset int64 = 1
	const limit int64 = 3
	persons := []models.Person{
		s.newPerson(1),
		s.newPerson(2),
		s.newPerson(3),
		s.newPerson(4),
		s.newPerson(5),
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := new(RepositoryPositiveMock)
	repo.On("GetPersons", offset, limit).Return(persons[offset:][:limit], nil)
	useCase := usecase.New(repo, logger)
	// act
	res, err := useCase.GetPersons(offset, limit)
	// assert
	t.Require().NoError(err)
	t.Assert().Len(res, int(limit))
	t.Assert().ElementsMatch(persons[offset:][:limit], res)
	repo.AssertExpectations(t)
	repo.AssertNumberOfCalls(t, "GetPersons", 1)
}

func (s *UseCaseSuite) TestCreatePerson(t provider.T) {
	t.Epic("MVP")
	t.Severity(allure.NORMAL)

	// arrange
	const personID = 5
	person := s.newPerson(personID)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := new(RepositoryPositiveMock)
	repo.On("CreatePerson", person.PersonProperties).Return(person, nil)
	useCase := usecase.New(repo, logger)
	// act
	res, err := useCase.CreatePerson(person.PersonProperties)
	// assert
	t.Require().NoError(err)
	t.Require().Equal(person, res)
	repo.AssertExpectations(t)
	repo.AssertNumberOfCalls(t, "CreatePerson", 1)
}

func (s *UseCaseSuite) TestGetPerson(t provider.T) {
	t.Epic("MVP")
	t.Severity(allure.NORMAL)

	// arrange
	const personID = 5
	person := s.newPerson(personID)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := new(RepositoryPositiveMock)
	repo.On("GetPerson", person.ID).Return(person, true, nil)
	useCase := usecase.New(repo, logger)
	// act
	res, found, err := useCase.GetPerson(person.ID)
	// assert
	t.Require().NoError(err)
	t.Require().True(found)
	t.Require().Equal(person, res)
	repo.AssertExpectations(t)
	repo.AssertNumberOfCalls(t, "GetPerson", 1)
}

func (s *UseCaseSuite) TestUpdatePerson(t provider.T) {
	t.Epic("MVP")
	t.Severity(allure.NORMAL)

	// arrange
	const personID = 5
	person := s.newPerson(personID)
	personOverride := models.Person{
		ID:               personID,
		PersonProperties: models.PersonProperties{Work: "New Work"},
	}
	newPerson := person
	newPerson.Work = "New Work"
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := new(RepositoryPositiveMock)
	repo.On("UpdatePerson", personOverride).Return(newPerson, true, nil)
	useCase := usecase.New(repo, logger)
	// act
	res, found, err := useCase.UpdatePerson(personOverride)
	// assert
	t.Require().NoError(err)
	t.Require().True(found)
	t.Require().Equal(newPerson, res)
	repo.AssertExpectations(t)
	repo.AssertNumberOfCalls(t, "UpdatePerson", 1)
}

func (s *UseCaseSuite) TestDeletePerson(t provider.T) {
	t.Epic("MVP")
	t.Severity(allure.NORMAL)

	// arrange
	const personID = 5
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := new(RepositoryPositiveMock)
	repo.On("DeletePerson", personID).Return(true, nil)
	useCase := usecase.New(repo, logger)
	// act
	found, err := useCase.DeletePerson(personID)
	// assert
	t.Require().NoError(err)
	t.Require().True(found)
	repo.AssertExpectations(t)
	repo.AssertNumberOfCalls(t, "DeletePerson", 1)
}

func TestUseCase(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	suite.RunSuite(t, new(UseCaseSuite))
}
