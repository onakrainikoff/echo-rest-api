package service

import (
	"echo-rest-api/model"
	"echo-rest-api/store"
)

type CategoryService interface {
	// Получить категорию по ид
	GetCategory(id int) (*model.Category, error)
	// Получить все категории
	GetCategories() ([]*model.Category, error)
	// Создать категорию
	CreateCategory(category *model.Category) (*int, error)
	// Обновить категорию
	UpdateCategory(category *model.Category) error
	// Удалить категорию
	DeleteCategory(id int) error
}

func NewCategoryService(store store.Store) CategoryService {
	return &CategoryServiceContext{store: store}
}

type CategoryServiceContext struct {
	store store.Store
}

func (csc *CategoryServiceContext) GetCategory(id int) (*model.Category, error) {
	return csc.store.GetCategory(nil, id)
}

func (csc *CategoryServiceContext) GetCategories() ([]*model.Category, error) {
	return csc.store.GetCategories(nil)
}

func (csc *CategoryServiceContext) CreateCategory(category *model.Category) (*int, error) {
	tx, err := csc.store.Begin()
	if err != nil {
		return nil, err
	}
	cat, err := csc.store.CreateCategory(tx, category)
	if err != nil {
		csc.store.Rollback(tx)
		return nil, err
	}
	if err = csc.store.Commit(tx); err != nil {
		return nil, err
	}
	return cat, nil
}

func (csc *CategoryServiceContext) UpdateCategory(category *model.Category) error {
	tx, err := csc.store.Begin()
	if err != nil {
		return err
	}
	err = csc.store.UpdateCategory(tx, category)
	if err != nil {
		csc.store.Rollback(tx)
		return err
	}
	if err = csc.store.Commit(tx); err != nil {
		return err
	}
	return nil
}

func (csc *CategoryServiceContext) DeleteCategory(id int) error {
	tx, err := csc.store.Begin()
	if err != nil {
		return err
	}
	err = csc.store.DeleteCategory(tx, id)
	if err != nil {
		csc.store.Rollback(tx)
		return err
	}
	if err = csc.store.Commit(tx); err != nil {
		return err
	}
	return nil
}
