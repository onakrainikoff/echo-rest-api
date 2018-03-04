package test

import (
	"database/sql"
	"echo-rest-api/api"
	"echo-rest-api/config"
	"echo-rest-api/model"
	"echo-rest-api/test/mock"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApi_GetCategories(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conf := &config.Config{LogLevel: 0}
	cs := mock.NewMockCategoryService(mockCtrl)
	api := api.NewApi(conf, cs, nil)
	req := httptest.NewRequest(echo.GET, "/api/categories", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// 200 [] - ничего не найдено
	rec := httptest.NewRecorder()
	var cats []*model.Category
	cs.EXPECT().GetCategories().Return(cats, nil).Times(1)
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, rec.Body.String(), "[]")
	// 200 - ок
	cats = append(cats, &model.Category{Id: 1, Name: "Name1"})
	cats = append(cats, &model.Category{Id: 2, Name: "Name2"})
	cs.EXPECT().GetCategories().Return(cats, nil).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	res, _ := json.Marshal(cats)
	assert.Equal(t, rec.Body.String(), string(res))
}

func TestApi_GetCategory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conf := &config.Config{LogLevel: 5}
	cs := mock.NewMockCategoryService(mockCtrl)
	api := api.NewApi(conf, cs, nil)
	req := httptest.NewRequest(echo.GET, "/api/categories/2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// 404
	rec := httptest.NewRecorder()
	cs.EXPECT().GetCategory(2).Return(nil, nil).Times(1)
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	// 200
	cat := &model.Category{Id: 2, Name: "Name2"}
	cs.EXPECT().GetCategory(2).Return(cat, nil).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	res, _ := json.Marshal(cat)
	assert.Equal(t, rec.Body.String(), string(res))
}

func TestApi_CreateCategory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conf := &config.Config{LogLevel: 5}
	cs := mock.NewMockCategoryService(mockCtrl)
	api := api.NewApi(conf, cs, nil)
	// 400
	catJSON := `{"name": "te"}`
	req := httptest.NewRequest(echo.POST, "/api/categories/", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// 200
	catJSON = `{"name": "test"}`
	id := 2
	req = httptest.NewRequest(echo.POST, "/api/categories/", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	cs.EXPECT().CreateCategory(gomock.Any()).Return(&id, nil).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
	var res map[string]int
	json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(t, res["id"], 2)
}

func TestApi_UpdateCategory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conf := &config.Config{LogLevel: 5}
	cs := mock.NewMockCategoryService(mockCtrl)
	api := api.NewApi(conf, cs, nil)
	// 400
	catJSON := `{"name": "te"}`
	req := httptest.NewRequest(echo.PUT, "/api/categories/2", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// 404
	catJSON = `{"name": "test"}`
	req = httptest.NewRequest(echo.PUT, "/api/categories/1", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	cs.EXPECT().UpdateCategory(gomock.Any()).Return(sql.ErrNoRows).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	// 204
	catJSON = `{"name": "test"}`
	req = httptest.NewRequest(echo.PUT, "/api/categories/2", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	cs.EXPECT().UpdateCategory(gomock.Any()).Return(nil).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestApi_DeleteCategory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conf := &config.Config{LogLevel: 5}
	cs := mock.NewMockCategoryService(mockCtrl)
	api := api.NewApi(conf, cs, nil)
	// 404
	req := httptest.NewRequest(echo.DELETE, "/api/categories/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	cs.EXPECT().DeleteCategory(gomock.Any()).Return(sql.ErrNoRows).Times(1)
	rec := httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	// 201
	req = httptest.NewRequest(echo.DELETE, "/api/categories/2", nil)
	cs.EXPECT().DeleteCategory(gomock.Any()).Return(nil).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestApi_GetProducts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conf := &config.Config{LogLevel: 5}
	ps := mock.NewMockProductService(mockCtrl)
	api := api.NewApi(conf, nil, ps)
	req := httptest.NewRequest(echo.GET, "/api/products", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// 200 [] - ничего не найдено
	rec := httptest.NewRecorder()
	var cats []*model.Product
	ps.EXPECT().GetProducts(gomock.Any()).Return(cats, nil).Times(1)
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, rec.Body.String(), "[]")
	// 200 - ок
	cats = append(cats, &model.Product{Id: 1, Name: "Name1"})
	cats = append(cats, &model.Product{Id: 2, Name: "Name2"})
	ps.EXPECT().GetProducts(gomock.Any()).Return(cats, nil).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	res, _ := json.Marshal(cats)
	assert.Equal(t, rec.Body.String(), string(res))
	// 200
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(echo.GET, "/api/products?category=2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	id := 2
	ps.EXPECT().GetProducts(&id).Return(cats, nil).Times(1)
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	res, _ = json.Marshal(cats)
	assert.Equal(t, rec.Body.String(), string(res))
}

func TestApi_GetProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conf := &config.Config{LogLevel: 5}
	ps := mock.NewMockProductService(mockCtrl)
	api := api.NewApi(conf, nil, ps)
	req := httptest.NewRequest(echo.GET, "/api/products/2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// 404
	rec := httptest.NewRecorder()
	ps.EXPECT().GetProduct(2).Return(nil, nil).Times(1)
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	// 200
	cat := &model.Product{Id: 2, Name: "Name2"}
	ps.EXPECT().GetProduct(2).Return(cat, nil).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	res, _ := json.Marshal(cat)
	assert.Equal(t, rec.Body.String(), string(res))
}

func TestApi_CreateProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conf := &config.Config{LogLevel: 5}
	ps := mock.NewMockProductService(mockCtrl)
	api := api.NewApi(conf, nil, ps)
	// 400
	catJSON := `{"name": "test","description":"test","category":1}`
	req := httptest.NewRequest(echo.POST, "/api/products/", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// 200
	catJSON = `{"name": "test","description":"test","category":1,"price":10.1}`
	id := 2
	req = httptest.NewRequest(echo.POST, "/api/products/", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ps.EXPECT().CreateProduct(gomock.Any()).Return(&id, nil).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
	var res map[string]int
	json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(t, res["id"], 2)
}

func TestApi_UpdateProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conf := &config.Config{LogLevel: 5}
	ps := mock.NewMockProductService(mockCtrl)
	api := api.NewApi(conf, nil, ps)
	// 400
	catJSON := `{"name": "test","description":"test","category":1}`
	req := httptest.NewRequest(echo.PUT, "/api/products/2", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// 404
	catJSON = `{"name": "test","description":"test","category":1,"price":10.1}`
	req = httptest.NewRequest(echo.PUT, "/api/products/1", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ps.EXPECT().UpdateProduct(gomock.Any()).Return(sql.ErrNoRows).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	// 204
	catJSON = `{"name": "test","description":"test","category":1,"price":10.1}`
	req = httptest.NewRequest(echo.PUT, "/api/products/2", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ps.EXPECT().UpdateProduct(gomock.Any()).Return(nil).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestApi_DeleteProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conf := &config.Config{LogLevel: 5}
	ps := mock.NewMockProductService(mockCtrl)
	api := api.NewApi(conf, nil, ps)
	// 404
	req := httptest.NewRequest(echo.DELETE, "/api/products/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ps.EXPECT().DeleteProduct(gomock.Any()).Return(sql.ErrNoRows).Times(1)
	rec := httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	// 201
	req = httptest.NewRequest(echo.DELETE, "/api/products/2", nil)
	ps.EXPECT().DeleteProduct(gomock.Any()).Return(nil).Times(1)
	rec = httptest.NewRecorder()
	api.Http.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
