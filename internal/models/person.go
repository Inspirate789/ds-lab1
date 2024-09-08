package models

type Person struct {
	ID int
	PersonProperties
}

type PersonProperties struct {
	Name    string
	Age     int
	Address string
	Work    string
}
