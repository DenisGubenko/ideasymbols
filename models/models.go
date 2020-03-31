package models

type Order struct {
	Content string
	Active  bool
	Counter uint64
}
