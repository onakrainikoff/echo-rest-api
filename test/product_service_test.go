package test

import (
	"database/sql"
	"echo-rest-api/model"
	"echo-rest-api/service"
	"echo-rest-api/test/mock"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductService_GetProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := mock.NewMockStore(mockCtrl)
	mockStore.EXPECT().GetProduct(nil, 1).Return(nil, errors.New("test")).Times(1)
	mockStore.EXPECT().GetProduct(nil, 2).Return(&model.Product{Id: 2, Name: "test"}, nil).Times(1)
	ps := service.NewProductService(mockStore)
	r, e := ps.GetProduct(1)
	assert.NotNil(t, e)
	assert.Nil(t, r)
	r, e = ps.GetProduct(2)
	assert.Nil(t, e)
	assert.NotNil(t, r)
}

func TestProductService_GetProducts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := mock.NewMockStore(mockCtrl)
	mockStore.EXPECT().GetProducts(nil, nil).Return(nil, errors.New("test")).Times(1)
	ps := service.NewProductService(mockStore)
	r, e := ps.GetProducts(nil)
	assert.NotNil(t, e)
	assert.Nil(t, r)
	mockStore.EXPECT().GetProducts(nil, nil).Return([]*model.Product{}, nil).Times(1)
	r, e = ps.GetProducts(nil)
	assert.Nil(t, e)
	assert.NotNil(t, r)
}

func TestProductService_CreateProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock.NewMockStore(mockCtrl)
	mockStore.EXPECT().Begin().Return(nil, errors.New("test")).Times(1)
	ps := service.NewProductService(mockStore)
	r, e := ps.CreateProduct(&model.Product{Name: "test"})
	assert.NotNil(t, e)
	assert.Nil(t, r)

	mockStore = mock.NewMockStore(mockCtrl)
	tx := new(sql.Tx)
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().CreateProduct(tx, &model.Product{Name: "test"}).Return(nil, errors.New("test")).Times(1)
	mockStore.EXPECT().Rollback(tx).Return(nil).Times(1)
	ps = service.NewProductService(mockStore)
	r, e = ps.CreateProduct(&model.Product{Name: "test"})
	assert.NotNil(t, e)
	assert.Nil(t, r)

	mockStore = mock.NewMockStore(mockCtrl)
	tx = new(sql.Tx)
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	var id = 1
	mockStore.EXPECT().CreateProduct(tx, &model.Product{Name: "test"}).Return(&id, nil).Times(1)
	mockStore.EXPECT().Commit(tx).Return(nil).Times(1)
	ps = service.NewProductService(mockStore)
	r, e = ps.CreateProduct(&model.Product{Name: "test"})
	assert.Nil(t, e)
	assert.NotNil(t, r)
}

func TestProductService_UpdateProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock.NewMockStore(mockCtrl)
	mockStore.EXPECT().Begin().Return(nil, errors.New("test")).Times(1)
	ps := service.NewProductService(mockStore)
	e := ps.UpdateProduct(nil)
	assert.NotNil(t, e)

	mockStore = mock.NewMockStore(mockCtrl)
	tx := new(sql.Tx)
	prod := &model.Product{Id: 1, Name: "test"}
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().UpdateProduct(tx, prod).Return(errors.New("test")).Times(1)
	mockStore.EXPECT().Rollback(tx).Return(nil).Times(1)
	ps = service.NewProductService(mockStore)
	e = ps.UpdateProduct(prod)
	assert.NotNil(t, e)

	mockStore = mock.NewMockStore(mockCtrl)
	tx = new(sql.Tx)
	prod = &model.Product{Id: 1, Name: "test"}
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().UpdateProduct(tx, prod).Return(nil).Times(1)
	mockStore.EXPECT().Commit(tx).Return(nil).Times(1)
	ps = service.NewProductService(mockStore)
	e = ps.UpdateProduct(prod)
	assert.Nil(t, e)
}

func TestProductService_DeleteProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock.NewMockStore(mockCtrl)
	mockStore.EXPECT().Begin().Return(nil, errors.New("test")).Times(1)
	ps := service.NewProductService(mockStore)
	e := ps.DeleteProduct(1)
	assert.NotNil(t, e)

	mockStore = mock.NewMockStore(mockCtrl)
	tx := new(sql.Tx)
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().DeleteProduct(tx, 1).Return(errors.New("test")).Times(1)
	mockStore.EXPECT().Rollback(tx).Return(nil).Times(1)
	ps = service.NewProductService(mockStore)
	e = ps.DeleteProduct(1)
	assert.NotNil(t, e)

	mockStore = mock.NewMockStore(mockCtrl)
	tx = new(sql.Tx)
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().DeleteProduct(tx, 1).Return(nil).Times(1)
	mockStore.EXPECT().Commit(tx).Return(nil).Times(1)
	ps = service.NewProductService(mockStore)
	e = ps.DeleteProduct(1)
	assert.Nil(t, e)
}
