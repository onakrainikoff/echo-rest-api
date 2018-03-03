package service

import (
	"echo-rest-api/model"
	"echo-rest-api/store"
)

type ProductService interface {
	// Получить продукт по id
	GetProduct(id int) (*model.Product, error)
	// Получить все продукты
	GetProducts(category *int) ([]*model.Product, error)
	// Создать продукт
	CreateProduct(product *model.Product) (*int, error)
	// Обновить продукт
	UpdateProduct(product *model.Product) error
	// Удалить продукт
	DeleteProduct(id int) error
}

func NewProductService(store store.Store) ProductService {
	return &ProductServiceContext{store: store}
}

type ProductServiceContext struct {
	store store.Store
}

func (psc *ProductServiceContext) GetProduct(id int) (*model.Product, error) {
	return psc.store.GetProduct(nil, id)
}

func (psc *ProductServiceContext) GetProducts(category *int) ([]*model.Product, error) {
	return psc.store.GetProducts(nil, category)
}

func (psc *ProductServiceContext) CreateProduct(product *model.Product) (*int, error) {
	tx, err := psc.store.Begin()
	if err != nil {
		return nil, err
	}
	cat, err := psc.store.CreateProduct(tx, product)
	if err != nil {
		psc.store.Rollback(tx)
		return nil, err
	}
	if err = psc.store.Commit(tx); err != nil {
		return nil, err
	}
	return cat, nil
}

func (psc *ProductServiceContext) UpdateProduct(product *model.Product) error {
	tx, err := psc.store.Begin()
	if err != nil {
		return err
	}
	err = psc.store.UpdateProduct(tx, product)
	if err != nil {
		psc.store.Rollback(tx)
		return err
	}
	if err = psc.store.Commit(tx); err != nil {
		return err
	}
	return nil
}

func (psc *ProductServiceContext) DeleteProduct(id int) error {
	tx, err := psc.store.Begin()
	if err != nil {
		return err
	}
	err = psc.store.DeleteProduct(tx, id)
	if err != nil {
		psc.store.Rollback(tx)
		return err
	}
	if err = psc.store.Commit(tx); err != nil {
		return err
	}
	return nil
}
