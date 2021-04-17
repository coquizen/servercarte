package gorm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/CaninoDev/gastro/server/domain/menu"
	"github.com/CaninoDev/gastro/server/internal/logger"
)

var (
	ErrSectionNotFound = errors.New("section could not be found in the database")
	ErrItemNotFound = errors.New("item could not be found in the database")
)
// menuRepository represents the client to its persistent repository
type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository instantiates an instance for data persistence
func NewMenuRepository(db *gorm.DB) *menuRepository {
	return &menuRepository{db}
}

// ListMenus lists all the menus in the db
func (r *menuRepository) ListMenus(_ context.Context,) (*[]menu.Section, error) {
	var sections []menu.Section


	if err := r.db.Preload("Items").Preload("SubSections.Items").Preload(clause.Associations).Where(
		"type = ?",
		menu.Meal).Find(&sections).Error; err != nil {
		logger.Error.Printf("db connection error %v", err)
		return &sections, err
	}
	return &sections, nil
}
// ListSections lists all the sections in the db
func (r *menuRepository) ListSections(_ context.Context,) (*[]menu.Section, error) {
	var sections []menu.Section

	if err := r.db.Preload(clause.Associations).Find(&sections).Error; err != nil {
		logger.Error.Printf("db connection error %v", err)
		return &sections, err
	}
	return &sections, nil
}
// FindSection finds an section by its id
func (r *menuRepository) FindSection(_ context.Context,section *menu.Section) error {
	if err := r.db.Preload("Items").Preload(clause.Associations).First(section).Error; errors.Is(err,
		gorm.ErrRecordNotFound) {
		return fmt.Errorf("record not found for %v", section.ID)
	} else if err != nil {
		logger.Error.Printf("db connection error %v", err)
		return err
	}
	return nil
}

// CreateSection first checks for preexisting record, and if not found will create the specified section
func (r *menuRepository) CreateSection(_ context.Context,section *menu.Section) error {
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
func (r *menuRepository) UpdateSection(_ context.Context,section *menu.Section) error {
	return r.db.Save(section).Error
}

// UpdateSectionParent re-parents a subsection
func (r *menuRepository) UpdateSectionParent(_ context.Context, child *menu.Section, newParent *menu.Section) error {
	return r.db.Model(&newParent).Association("SubSections").Append(&child)
}

// DeleteSection deletes a section
func (r *menuRepository) DeleteSection(_ context.Context,section *menu.Section) error {
	return r.db.Delete(section).Error
}

// ListItems lists all the users in the db
func (r *menuRepository) ListItems(_ context.Context,) (*[]menu.Item, error) {
	var items []menu.Item

	if err := r.db.Preload(clause.Associations).Find(&items).Error; err != nil {
		logger.Error.Printf("db connection error %v", err)
		return &items, err
	}
	return &items, nil
}

// FindItem finds an item by its id
func (r *menuRepository) FindItem(_ context.Context, item *menu.Item) error {
	if err := r.db.Preload(clause.Associations).First(item).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("record not found for %v", item.ID)
	} else if err != nil {
		logger.Error.Printf("db connection error %v", err)
		return err
	}
	return nil
}

// CreateItem first checks for preexisting record, and if not found will create the specified item
func (r *menuRepository) CreateItem(_ context.Context, item *menu.Item) error {
	var checkItem = new(menu.Item)
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
func (r *menuRepository) UpdateItem(_ context.Context, item *menu.Item) error {
	return r.db.Save(&item).Error
}


// UpdateItemParent re-parents an item
func (r *menuRepository) UpdateItemParent(_ context.Context, child *menu.Item, newParent *menu.Section) error {
	return r.db.Model(&newParent).Association("Items").Append(&child)
}

// DeleteItem deletes an item
func (r *menuRepository) DeleteItem(_ context.Context, item *menu.Item) error {
	if err := r.db.Preload(clause.Associations).First(&item).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else if err != nil {
		logger.Error.Printf("db connection error %v", err)
		return err
	}
	return r.db.Delete(&item).Error
}
