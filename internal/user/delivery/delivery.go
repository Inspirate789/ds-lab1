package delivery

import (
	"github.com/Inspirate789/ds-lab1/internal/models"
	"github.com/Inspirate789/ds-lab1/internal/user/delivery/errors"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"math"
	"strconv"
)

type UseCase interface { // TODO: add context
	HealthCheck() error
	GetPersons(offset, limit int64) ([]models.Person, error)
	CreatePerson(person models.PersonProperties) (models.Person, error)
	GetPerson(personID int) (models.Person, bool, error)
	UpdatePerson(person models.Person) (models.Person, bool, error)
	DeletePerson(personID int) (bool, error)
}

type delivery struct {
	useCase UseCase
	logger  *slog.Logger
}

func AddHandlers(api fiber.Router, useCase UseCase, logger *slog.Logger) {
	handler := &delivery{
		useCase: useCase,
		logger:  logger,
	}

	api.Get("/", handler.GetPersons)
	api.Post("/", handler.PostPerson)

	api.Get("/:personId", handler.GetPerson)
	api.Patch("/:personId", handler.PatchPerson)
	api.Delete("/:personId", handler.DeletePerson)
}

func (d *delivery) GetPersons(ctx *fiber.Ctx) error {
	offset, err := strconv.ParseInt(ctx.Query("offset"), 10, 64)
	if err != nil {
		d.logger.Warn("persons offset not set, use default 0")
		offset = 0
	}

	limit, err := strconv.ParseInt(ctx.Query("limit"), 10, 64)
	if err != nil {
		d.logger.Warn("persons limit not set, return all persons")
		limit = math.MaxInt64
	}

	persons, err := d.useCase.GetPersons(offset, limit)
	if err != nil {
		return err
	}

	if persons == nil {
		persons = make([]models.Person, 0)
	}

	return ctx.Status(fiber.StatusOK).JSON(NewPersonsDTO(persons)) // TODO: .Map()
}

func (d *delivery) PostPerson(ctx *fiber.Ctx) error {
	var dto PersonProperties

	err := ctx.BodyParser(&dto)
	if err != nil {
		personErr := errors.ErrInvalidPerson(err.Error())
		d.logger.Error(err.Error())

		return ctx.Status(fiber.StatusBadRequest).JSON(personErr.Map())
	}

	person, err := d.useCase.CreatePerson(dto.ToProperties())
	if err != nil {
		return err
	}

	ctx.Location(ctx.Path() + "/" + strconv.Itoa(person.ID))

	return ctx.SendStatus(fiber.StatusCreated)
}

func (d *delivery) GetPerson(ctx *fiber.Ctx) error {
	personID, err := strconv.Atoi(ctx.Params("personId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors.ErrInvalidID.Map())
	}

	person, found, err := d.useCase.GetPerson(personID)
	if err != nil {
		return err
	}

	if !found {
		return ctx.Status(fiber.StatusNotFound).JSON(errors.ErrPersonNotFound.Map())
	}

	return ctx.Status(fiber.StatusOK).JSON(NewPersonDTO(person))
}

func (d *delivery) PatchPerson(ctx *fiber.Ctx) error {
	personID, err := strconv.Atoi(ctx.Params("personId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors.ErrInvalidID.Map())
	}

	var dto PersonProperties

	err = ctx.BodyParser(&dto)
	if err != nil {
		personErr := errors.ErrInvalidPerson(err.Error())
		d.logger.Error(err.Error())

		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(personErr.Map())
	}

	person, found, err := d.useCase.UpdatePerson(dto.ToPerson(personID))
	if err != nil {
		return err
	}

	if !found {
		return ctx.Status(fiber.StatusNotFound).JSON(errors.ErrPersonNotFound.Map())
	}

	return ctx.Status(fiber.StatusOK).JSON(NewPersonDTO(person))
}

func (d *delivery) DeletePerson(ctx *fiber.Ctx) error {
	personID, err := strconv.Atoi(ctx.Params("personId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors.ErrInvalidID.Map())
	}

	found, err := d.useCase.DeletePerson(personID)
	if err != nil {
		return err
	}

	if !found {
		return ctx.Status(fiber.StatusNotFound).JSON(errors.ErrPersonNotFound.Map())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
