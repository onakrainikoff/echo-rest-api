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

func TestCategoryService_GetCategory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := mock.NewMockStore(mockCtrl)
	mockStore.EXPECT().GetCategory(nil, 1).Return(nil, errors.New("test")).Times(1)
	mockStore.EXPECT().GetCategory(nil, 2).Return(&model.Category{2, "test"}, nil).Times(1)
	cs := service.NewCategoryService(mockStore)
	r, e := cs.GetCategory(1)
	assert.NotNil(t, e)
	assert.Nil(t, r)
	r, e = cs.GetCategory(2)
	assert.Nil(t, e)
	assert.NotNil(t, r)
}

func TestCategoryService_GetCategories(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := mock.NewMockStore(mockCtrl)
	mockStore.EXPECT().GetCategories(nil).Return(nil, errors.New("test")).Times(1)
	cs := service.NewCategoryService(mockStore)
	r, e := cs.GetCategories()
	assert.NotNil(t, e)
	assert.Nil(t, r)
	mockStore.EXPECT().GetCategories(nil).Return([]*model.Category{}, nil).Times(1)
	r, e = cs.GetCategories()
	assert.Nil(t, e)
	assert.NotNil(t, r)
}

func TestCategoryService_CreateCategory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock.NewMockStore(mockCtrl)
	mockStore.EXPECT().Begin().Return(nil, errors.New("test")).Times(1)
	cs := service.NewCategoryService(mockStore)
	r, e := cs.CreateCategory(&model.Category{Name: "Test"})
	assert.NotNil(t, e)
	assert.Nil(t, r)

	mockStore = mock.NewMockStore(mockCtrl)
	tx := new(sql.Tx)
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().CreateCategory(tx, &model.Category{Name: "Test"}).Return(nil, errors.New("test")).Times(1)
	mockStore.EXPECT().Rollback(tx).Return(nil).Times(1)
	cs = service.NewCategoryService(mockStore)
	r, e = cs.CreateCategory(&model.Category{Name: "Test"})
	assert.NotNil(t, e)
	assert.Nil(t, r)

	mockStore = mock.NewMockStore(mockCtrl)
	tx = new(sql.Tx)
	var id = 1
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().CreateCategory(tx, &model.Category{Name: "Test"}).Return(&id, nil).Times(1)
	mockStore.EXPECT().Commit(tx).Return(nil).Times(1)
	cs = service.NewCategoryService(mockStore)
	r, e = cs.CreateCategory(&model.Category{Name: "Test"})
	assert.Nil(t, e)
	assert.NotNil(t, r)
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock.NewMockStore(mockCtrl)
	mockStore.EXPECT().Begin().Return(nil, errors.New("test")).Times(1)
	cs := service.NewCategoryService(mockStore)
	e := cs.UpdateCategory(nil)
	assert.NotNil(t, e)

	mockStore = mock.NewMockStore(mockCtrl)
	tx := new(sql.Tx)
	cat := &model.Category{1, "test"}
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().UpdateCategory(tx, cat).Return(errors.New("test")).Times(1)
	mockStore.EXPECT().Rollback(tx).Return(nil).Times(1)
	cs = service.NewCategoryService(mockStore)
	e = cs.UpdateCategory(cat)
	assert.NotNil(t, e)

	mockStore = mock.NewMockStore(mockCtrl)
	tx = new(sql.Tx)
	cat = &model.Category{1, "test"}
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().UpdateCategory(tx, cat).Return(nil).Times(1)
	mockStore.EXPECT().Commit(tx).Return(nil).Times(1)
	cs = service.NewCategoryService(mockStore)
	e = cs.UpdateCategory(cat)
	assert.Nil(t, e)
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock.NewMockStore(mockCtrl)
	mockStore.EXPECT().Begin().Return(nil, errors.New("test")).Times(1)
	cs := service.NewCategoryService(mockStore)
	e := cs.DeleteCategory(1)
	assert.NotNil(t, e)

	mockStore = mock.NewMockStore(mockCtrl)
	tx := new(sql.Tx)
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().DeleteCategory(tx, 1).Return(errors.New("test")).Times(1)
	mockStore.EXPECT().Rollback(tx).Return(nil).Times(1)
	cs = service.NewCategoryService(mockStore)
	e = cs.DeleteCategory(1)
	assert.NotNil(t, e)

	mockStore = mock.NewMockStore(mockCtrl)
	tx = new(sql.Tx)
	mockStore.EXPECT().Begin().Return(tx, nil).Times(1)
	mockStore.EXPECT().DeleteCategory(tx, 1).Return(nil).Times(1)
	mockStore.EXPECT().Commit(tx).Return(nil).Times(1)
	cs = service.NewCategoryService(mockStore)
	e = cs.DeleteCategory(1)
	assert.Nil(t, e)
}
