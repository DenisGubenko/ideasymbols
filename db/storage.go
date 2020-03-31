package db

import "github.com/DenisGubenko/ideasymbols/models"

type Storage interface {
	CreateOrder(request *models.Order) error
	GetRandomOrderContent() (*models.Order, error)
	InactiveRandomOrder() error
	GetStatisticsOrder() (*uint64, []*models.Order, error)
}
