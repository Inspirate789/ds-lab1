package repository

import "github.com/Inspirate789/ds-lab1/internal/models"

type Person struct {
	ID int `sqlx:"id"`
	PersonProperties
}

type PersonProperties struct {
	Name    string `sqlx:"name"`
	Age     int    `sqlx:"age"`
	Address string `sqlx:"address"`
	Work    string `sqlx:"work"`
}

func (p Person) UpdateBy(properties PersonProperties) Person {
	if properties.Name != "" {
		p.Name = properties.Name
	}

	if properties.Age != 0 {
		p.Age = properties.Age
	}

	if properties.Address != "" {
		p.Address = properties.Address
	}

	if properties.Work != "" {
		p.Work = properties.Work
	}

	return p
}

func NewPerson(person models.Person) Person {
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

func NewPersonProperties(person models.PersonProperties) PersonProperties {
	return PersonProperties{
		Name:    person.Name,
		Age:     person.Age,
		Address: person.Address,
		Work:    person.Work,
	}
}

func (p Person) ToModel() models.Person {
	return models.Person{
		ID: p.ID,
		PersonProperties: models.PersonProperties{
			Name:    p.Name,
			Age:     p.Age,
			Address: p.Address,
			Work:    p.Work,
		},
	}
}

type Persons []Person

func (p Persons) ToModel() []models.Person {
	dto := make([]models.Person, 0, len(p))

	for _, person := range p {
		dto = append(dto, person.ToModel())
	}

	return dto
}
