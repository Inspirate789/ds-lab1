package errors

import "github.com/gofiber/fiber/v2"

type PersonError string

func (e PersonError) Error() string {
	return string(e)
}

func (e PersonError) Map() map[string]any {
	return fiber.Map{"message": string(e)}
}

const (
	ErrInvalidID      PersonError = "invalid person ID"
	ErrPersonNotFound PersonError = "person not found"
)

func ErrInvalidPerson(msg string) PersonError {
	return PersonError("cannot parse person from request body: " + msg)
}
