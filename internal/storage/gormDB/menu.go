package gormDB

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/CaninoDev/gastro/server/api"
	"github.com/CaninoDev/gastro/server/internal/logger"
)

// MenuRepository represents the client to its persistent storage
type MenuRepository struct {
	db *gorm.DB
}

// NewMenuRepository instantiates an instance for data persistence
func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db}
}


// ListSections lists all the users in the db
func (r MenuRepository) ListSections(_ context.Context,) (*[]api.Section, error) {
	var sections []api.Section

	if err := r.db.Preload(clause.Associations).Find(&sections).Error; err != nil {
		logger.Error.Printf("db connection error %v", err)
		return &sections, err
	}
	return &sections, nil
}
// FindSection finds an section by its id
func (r MenuRepository) FindSection(_ context.Context,section *api.Section) error {
	if err := r.db.Preload("Items").Preload(clause.Associations).First(section).Error; errors.Is(err,
		gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf("record not found for %v", section.ID))
	} else if err != nil {
		logger.Error.Printf("db connection error %v", err)
		return err
	}
	return nil
}

// CreateSection first checks for preexisting record, and if not found will create the specified section
func (r MenuRepository) CreateSection(_ context.Context,section *api.Section) error {
	if err := r.db.Where(
		"lower(title) = ?",
		strings.ToLower(section.Title)).First(section).Error; err == nil {
		return errors.New("title already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := r.db.Create(&section).Error;  err != nil {
		return err
	}
	return nil
}

// UpdateSection updates section data
func (r MenuRepository) UpdateSection(_ context.Context,section *api.Section) error {
	return r.db.Save(section).Error
}

// UpdateSectionParent re-parents a subsection
func (r MenuRepository) UpdateSectionParent(_ context.Context, child *api.Section, newParent *api.Section) error {
	return r.db.Model(&newParent).Association("SubSections").Append(&child)
}

// DeleteSection deletes a section
func (r MenuRepository) DeleteSection(_ context.Context,section *api.Section) error {
	return r.db.Delete(section).Error
}

// ListItems lists all the users in the db
func (r MenuRepository) ListItems(_ context.Context,) (*[]api.Item, error) {
	var items []api.Item

	if err := r.db.Preload(clause.Associations).Find(&items).Error; err != nil {
		logger.Error.Printf("db connection error %v", err)
		return &items, err
	}
	return &items, nil
}

// FindItem finds an item by its id
func (r MenuRepository) FindItem(_ context.Context, item *api.Item) error {
	if err := r.db.Preload(clause.Associations).First(item).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf("record not found for %v", item.ID))
	} else if err != nil {
		logger.Error.Printf("db connection error %v", err)
		return err
	}
	return nil
}

// CreateItem first checks for preexisting record, and if not found will create the specified item
func (r MenuRepository) CreateItem(_ context.Context, item *api.Item) error {
	var checkItem = new(api.Item)
	if err := r.db.Where(
		"lower(title) = ?",
		strings.ToLower(item.Title)).First(&checkItem).Error; err == nil {
		return errors.New("title already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := r.db.Create(&item).Error;  err != nil {
		return err
	}

	return nil
}

// UpdateItem updates an item
func (r MenuRepository) UpdateItem(_ context.Context, item *api.Item) error {
	return r.db.Save(&item).Error
}


// UpdateItemParent re-parents an item
func (r MenuRepository) UpdateItemParent(_ context.Context, child *api.Item, newParent *api.Section) error {
	return r.db.Model(&newParent).Association("Items").Append(&child)
}

// DeleteItem deletes an item
func (r MenuRepository) DeleteItem(_ context.Context, item *api.Item) error {
	if err := r.db.Preload(clause.Associations).First(&item).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else if err != nil {
		logger.Error.Printf("db connection error %v", err)
		return err
	}
	return r.db.Delete(&item).Error
}
