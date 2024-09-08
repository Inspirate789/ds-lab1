package delivery

import (
	"github.com/Inspirate789/ds-lab1/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Person struct {
	ID int `json:"id"`
	PersonProperties
}

type PersonProperties struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

func NewPersonDTO(person models.Person) Person {
	return Person{
		ID: person.ID,
		PersonProperties: PersonProperties{
			Name:    person.Name,
			Age:     person.Age,
			Address: person.Address,
			Work:    person.Work,
		},
	}
}

func (p PersonProperties) ToPerson(id int) models.Person {
	return models.Person{
		ID: id,
		PersonProperties: models.PersonProperties{
			Name:    p.Name,
			Age:     p.Age,
			Address: p.Address,
			Work:    p.Work,
		},
	}
}

func (p PersonProperties) ToProperties() models.PersonProperties {
	return models.PersonProperties{
		Name:    p.Name,
		Age:     p.Age,
		Address: p.Address,
		Work:    p.Work,
	}
}

type Persons []Person

func NewPersonsDTO(persons []models.Person) Persons {
	dto := make([]Person, 0, len(persons))

	for _, person := range persons {
		dto = append(dto, NewPersonDTO(person))
	}

	return dto
}

func (p Persons) Map() map[string]any {
	return fiber.Map{
		"persons": p,
		"count":   len(p),
	}
}
