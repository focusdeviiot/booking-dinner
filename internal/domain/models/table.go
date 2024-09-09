package models

type Table struct {
	ID         string
	Capacity   int
	IsOccupied bool
}

func NewTable(id string, capacity int) *Table {
	return &Table{
		ID:         id,
		Capacity:   capacity,
		IsOccupied: false,
	}
}
