package menu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/CaninoDev/gastro/server/internal/logger"
	"github.com/CaninoDev/gastro/server/internal/model"
)

// GormDBRepository represents the client to its persistent storage
type GormDBRepository struct {
	db *gorm.DB
}

// NewGormDBRepository instantiates an instance for data persistence
func NewGormDBRepository(db *gorm.DB) *GormDBRepository {
	return &GormDBRepository{db}
}


// ListSections lists all the users in the db
func (r GormDBRepository) ListSections(_ context.Context,) (*[]model.Section, error) {
	var sections []model.Section

	if err := r.db.Preload(clause.Associations).Find(&sections).Error; err != nil {
		logger.Error.Printf("db connection error %v", err)
		return &sections, err
	}
	return &sections, nil
}
// FindSection finds an section by its id
func (r GormDBRepository) FindSection(_ context.Context,section *model.Section) error {
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
func (r GormDBRepository) CreateSection(_ context.Context,section *model.Section) error {
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
func (r GormDBRepository) UpdateSection(_ context.Context,section *model.Section) error {
	return r.db.Save(section).Error
}

// UpdateSectionParent re-parents a subsection
func (r GormDBRepository) UpdateSectionParent(_ context.Context, child *model.Section, newParent *model.Section) error {
	return r.db.Model(&newParent).Association("SubSections").Append(&child)
}

// DeleteSection deletes a section
func (r GormDBRepository) DeleteSection(_ context.Context,section *model.Section) error {
	return r.db.Delete(section).Error
}

// ListItems lists all the users in the db
func (r GormDBRepository) ListItems(_ context.Context,) (*[]model.Item, error) {
	var items []model.Item

	if err := r.db.Preload(clause.Associations).Find(&items).Error; err != nil {
		logger.Error.Printf("db connection error %v", err)
		return &items, err
	}
	return &items, nil
}

// FindItem finds an item by its id
func (r GormDBRepository) FindItem(_ context.Context, item *model.Item) error {
	if err := r.db.Preload(clause.Associations).First(item).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf("record not found for %v", item.ID))
	} else if err != nil {
		logger.Error.Printf("db connection error %v", err)
		return err
	}
	return nil
}

// CreateItem first checks for preexisting record, and if not found will create the specified item
func (r GormDBRepository) CreateItem(_ context.Context, item *model.Item) error {
	var checkItem = new(model.Item)
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
func (r GormDBRepository) UpdateItem(_ context.Context, item *model.Item) error {
	return r.db.Save(&item).Error
}


// UpdateItemParent re-parents an item
func (r GormDBRepository) UpdateItemParent(_ context.Context, child *model.Item, newParent *model.Section) error {
	return r.db.Model(&newParent).Association("Items").Append(&child)
}

// DeleteItem deletes an item
func (r GormDBRepository) DeleteItem(_ context.Context, item *model.Item) error {
	if err := r.db.Preload(clause.Associations).First(&item).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else if err != nil {
		logger.Error.Printf("db connection error %v", err)
		return err
	}
	return r.db.Delete(&item).Error
}
