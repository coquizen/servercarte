package repositories

import (
	"github.com/CaninoDev/gastro/server/models"
)

type ItemsRepositoryInterface interface {
	Items() ([]models.Item, error)
	FindByID(id string) (*models.Item, error)
	FindOrCreateItem(i *models.Item) (error)
	Update(i *models.Item) (*models.Item, error)
	Delete(i *models.Item) (error)
}


