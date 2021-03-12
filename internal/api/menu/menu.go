package menu

import (
	"context"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/internal/model"
)

type Menu struct {
	repo Repository
}

// New instantiates a new application interface
func Bind(repo Repository) *Menu {
	return &Menu{repo: repo}
}

func Initialize(repo Repository) *Menu {
	return Bind(repo)
}


func (m Menu) Sections(ctx context.Context) (*[]model.Section, error) {
	return m.repo.ListSections(ctx)
}

func (m Menu) SectionByID(ctx context.Context, rawID string) (*model.Section, error) {
 	id, err := uuid.Parse(rawID)
	if err != nil {
		return &model.Section{}, err
	}
	var section model.Section
	section.ID = id
	if err := m.repo.FindSection(ctx, &section); err != nil {
		return &section, err
	}
	return &section, nil
}

func (m Menu) NewSection(ctx context.Context, section *model.Section) error {
	return m.repo.CreateSection(ctx, section)
}

func (m Menu) UpdateSectionData(ctx context.Context, section *model.Section) error {
	return m.repo.UpdateSection(ctx, section)
}

func (m Menu) ReParentSection(ctx context.Context, section *model.Section, newParentID uuid.UUID) error {
	var newParentSection model.Section
	newParentSection.ID = newParentID
	if err := m.repo.FindSection(ctx, &newParentSection); err != nil {
		return err
	}

	return m.repo.UpdateSectionParent(ctx, section, &newParentSection)
}

func (m Menu) DeleteSection(ctx context.Context, rawID string) error {
	newSectionParentID, err := uuid.Parse(rawID)
	if err != nil {
		return err
	}
	var newParentSection model.Section
	newParentSection.ID = newSectionParentID
	if err := m.repo.FindSection(ctx, &newParentSection); err != nil {
		return err
	}
	return m.repo.DeleteSection(ctx, &newParentSection)
}

func (m Menu) Items(ctx context.Context) (*[]model.Item, error) {
	return m.repo.ListItems(ctx)
}

func (m Menu) ItemByID(ctx context.Context, rawID string) (*model.Item, error) {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return &model.Item{}, err
	}
	var item model.Item
	item.ID = id
	if err := m.repo.FindItem(ctx, &item); err != nil {
		return &item, err
	}

	return &item, nil
}

func (m Menu) ReParentItem(ctx context.Context, item *model.Item, newSectionParentID uuid.UUID) error {
	var newParentSection model.Section
	newParentSection.ID = newSectionParentID
	if err := m.repo.FindSection(ctx, &newParentSection); err != nil {
		return err
	}

	return m.repo.UpdateItemParent(ctx, item, &newParentSection)
}

func (m Menu) NewItem(ctx context.Context, item *model.Item) error {
	return m.repo.CreateItem(ctx, item)
}

func (m Menu) UpdateItemData(ctx context.Context, item *model.Item) error {
	return m.repo.UpdateItem(ctx, item)
}

func (m Menu) DeleteItem(ctx context.Context, rawID string) error {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return err
	}

	var item model.Item
	item.ID = id
	if err := m.repo.FindItem(ctx, &item); err != nil {
		return err
	}
	return m.repo.DeleteItem(ctx, &item)
}





