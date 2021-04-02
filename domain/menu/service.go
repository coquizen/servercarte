package menu

import (
	"context"

	"github.com/google/uuid"
)

// Service describes the expected behavior for manipulating service data.
type Service interface {
	Sections(context.Context) (*[]Section, error)
	SectionByID(context.Context, string) (*Section, error)
	NewSection(context.Context, *Section) error
	UpdateSectionContent(context.Context, *Section) error
	ReParentSection(context.Context, *Section, uuid.UUID) error
	DeleteSection(context.Context, string) error
	Items(context.Context) (*[]Item, error)
	ItemByID(context.Context, string) (*Item, error)
	NewItem(context.Context, *Item) error
	ReParentItem(context.Context, *Item, uuid.UUID) error
	UpdateItemContent(context.Context, *Item) error
	DeleteItem(context.Context, string) error
}

type service struct {
	repo Repository
}

func NewService(menuRepo Repository) *service {
	return &service{menuRepo}
}

func (m *service) NewSection(ctx context.Context, section *Section) error {
	return m.repo.CreateSection(ctx, section)
}

func (m *service) Sections(ctx context.Context) (*[]Section, error) {
	return m.repo.ListSections(ctx)
}

func (m *service) SectionByID(ctx context.Context, rawID string) (*Section, error) {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return &Section{}, err
	}
	var section Section
	section.ID = id
	if err := m.repo.FindSection(ctx, &section); err != nil {
		return &section, err
	}
	return &section, nil
}

func (m *service) UpdateSectionContent(ctx context.Context, section *Section) error {
	return m.repo.UpdateSection(ctx, section)
}

func (m *service) ReParentSection(ctx context.Context, section *Section, newParentID uuid.UUID) error {
	var newParentSection Section
	newParentSection.ID = newParentID
	if err := m.repo.FindSection(ctx, &newParentSection); err != nil {
		return err
	}

	return m.repo.UpdateSectionParent(ctx, section, &newParentSection)
}

func (m *service) DeleteSection(ctx context.Context, rawID string) error {
	newSectionParentID, err := uuid.Parse(rawID)
	if err != nil {
		return err
	}
	var newParentSection Section
	newParentSection.ID = newSectionParentID
	if err := m.repo.FindSection(ctx, &newParentSection); err != nil {
		return err
	}
	return m.repo.DeleteSection(ctx, &newParentSection)
}

func (m *service) NewItem(ctx context.Context, item *Item) error {
	return m.repo.CreateItem(ctx, item)
}

func (m *service) Items(ctx context.Context) (*[]Item, error) {
	return m.repo.ListItems(ctx)
}

func (m *service) ItemByID(ctx context.Context, rawID string) (*Item, error) {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return &Item{}, err
	}
	var item Item
	item.ID = id
	if err := m.repo.FindItem(ctx, &item); err != nil {
		return &item, err
	}

	return &item, nil
}

func (m *service) ReParentItem(ctx context.Context, item *Item, newSectionParentID uuid.UUID) error {
	var newParentSection Section
	newParentSection.ID = newSectionParentID
	if err := m.repo.FindSection(ctx, &newParentSection); err != nil {
		return err
	}

	return m.repo.UpdateItemParent(ctx, item, &newParentSection)
}

func (m *service) UpdateItemContent(ctx context.Context, item *Item) error {
	return m.repo.UpdateItem(ctx, item)
}

func (m *service) DeleteItem(ctx context.Context, rawID string) error {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return err
	}

	var item Item
	item.ID = id
	if err := m.repo.FindItem(ctx, &item); err != nil {
		return err
	}
	return m.repo.DeleteItem(ctx, &item)
}





