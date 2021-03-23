package menu

import (
	"context"

	"github.com/google/uuid"

	"github.com/CaninoDev/gastro/server/api"
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


func (m Menu) Sections(ctx context.Context) (*[]api.Section, error) {
	return m.repo.ListSections(ctx)
}

func (m Menu) SectionByID(ctx context.Context, rawID string) (*api.Section, error) {
 	id, err := uuid.Parse(rawID)
	if err != nil {
		return &api.Section{}, err
	}
	var section api.Section
	section.ID = id
	if err := m.repo.FindSection(ctx, &section); err != nil {
		return &section, err
	}
	return &section, nil
}

func (m Menu) NewSection(ctx context.Context, section *api.Section) error {
	return m.repo.CreateSection(ctx, section)
}

func (m Menu) UpdateSectionData(ctx context.Context, section *api.Section) error {
	return m.repo.UpdateSection(ctx, section)
}

func (m Menu) ReParentSection(ctx context.Context, section *api.Section, newParentID uuid.UUID) error {
	var newParentSection api.Section
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
	var newParentSection api.Section
	newParentSection.ID = newSectionParentID
	if err := m.repo.FindSection(ctx, &newParentSection); err != nil {
		return err
	}
	return m.repo.DeleteSection(ctx, &newParentSection)
}

func (m Menu) Items(ctx context.Context) (*[]api.Item, error) {
	return m.repo.ListItems(ctx)
}

func (m Menu) ItemByID(ctx context.Context, rawID string) (*api.Item, error) {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return &api.Item{}, err
	}
	var item api.Item
	item.ID = id
	if err := m.repo.FindItem(ctx, &item); err != nil {
		return &item, err
	}

	return &item, nil
}

func (m Menu) ReParentItem(ctx context.Context, item *api.Item, newSectionParentID uuid.UUID) error {
	var newParentSection api.Section
	newParentSection.ID = newSectionParentID
	if err := m.repo.FindSection(ctx, &newParentSection); err != nil {
		return err
	}

	return m.repo.UpdateItemParent(ctx, item, &newParentSection)
}

func (m Menu) NewItem(ctx context.Context, item *api.Item) error {
	return m.repo.CreateItem(ctx, item)
}

func (m Menu) UpdateItemData(ctx context.Context, item *api.Item) error {
	return m.repo.UpdateItem(ctx, item)
}

func (m Menu) DeleteItem(ctx context.Context, rawID string) error {
	id, err := uuid.Parse(rawID)
	if err != nil {
		return err
	}

	var item api.Item
	item.ID = id
	if err := m.repo.FindItem(ctx, &item); err != nil {
		return err
	}
	return m.repo.DeleteItem(ctx, &item)
}





