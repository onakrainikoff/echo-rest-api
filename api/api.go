// Echo-REST-API.
// Пример создания REST API на Echo framework
// Version: 0.0.1
// Schemes: http
// Host: localhost
// BasePath: /api
// Consumes:
// - application/json
// Produces:
// - application/json
// Contact: uchonyy@gmail.com
// swagger:meta
package api

import (
	"database/sql"
	"echo-rest-api/config"
	"echo-rest-api/model"
	"echo-rest-api/service"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type Api struct {
	Http     *echo.Echo
	conf     *config.Config
	cs       service.CategoryService
	ps       service.ProductService
	apiInfo  ApiInfo
	validate *validator.Validate
}

type ApiInfo struct {
	Address string
	MW      []string
	Routs   []string
}

func NewApi(conf *config.Config, cs service.CategoryService, ps service.ProductService) *Api {
	api := &Api{}
	api.validate = validator.New()
	api.conf = conf
	api.cs = cs
	api.ps = ps
	api.Http = echo.New()
	api.Http.Logger.SetLevel(log.Lvl(conf.LogLevel))
	api.apiInfo.Address = ":" + strconv.Itoa(api.conf.Api.HttpPort)
	api.Http.HideBanner = true
	api.Http.Pre(middleware.RemoveTrailingSlash())
	if conf.Api.Logging {
		api.Http.Use(middleware.Logger())
		api.apiInfo.MW = append(api.apiInfo.MW, "Logger")
	}
	api.Http.GET("/", api.index)
	api.Http.Static("/spec", "spec")
	api.Http.GET("/api/categories", api.getCategories)
	api.Http.GET("/api/categories/:id", api.getCategory)
	api.Http.POST("/api/categories", api.createCategory)
	api.Http.PUT("/api/categories/:id", api.updateCategory)
	api.Http.DELETE("/api/categories/:id", api.deleteCategory)

	api.Http.GET("/api/products", api.getProducts)
	api.Http.GET("/api/products/:id", api.getProduct)
	api.Http.POST("/api/products", api.createProduct)
	api.Http.PUT("/api/products/:id", api.updateProduct)
	api.Http.DELETE("/api/products/:id", api.deleteProduct)
	for _, r := range api.Http.Routes() {
		api.apiInfo.Routs = append(api.apiInfo.Routs, fmt.Sprintf("%s %s", r.Path, r.Method))
	}
	return api
}

// Запустить api
func (api *Api) Start() error {
	return api.Http.Start(":" + strconv.Itoa(api.conf.Api.HttpPort))
}

// Инфо об api
func (api *Api) GetApiInfo() ApiInfo {
	return api.apiInfo
}

func (api *Api) index(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/spec/api.json")
}

func (api *Api) spec(c echo.Context) error {
	return c.Inline("spec/api.json", "api.json")
}

// swagger:operation GET /categories/{id} getCategory
// ---
// description: Получить категорию
// parameters:
// - name: id
//   in: path
//   description: id необходимой категории
//   required: true
//   type: int
// responses:
//  '200':
//    schema:
//      $ref: '#/definitions/Category'
//  '400':
//     description: Bad request param `id`
//  '404':
//     description: Category `id`= not found
//
func (api *Api) getCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param `id`")
	}
	cat, err := api.cs.GetCategory(id)
	if err != nil {
		return err
	}
	if cat == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category `id` = ", id, " not found")
	}
	return c.JSON(http.StatusOK, cat)
}

// swagger:operation GET /categories getCategories
// ---
// description: Получить список категорий
// responses:
//  '200':
//    schema:
//      type: array
//      items:
//        $ref: '#/definitions/Category'
//
func (api *Api) getCategories(c echo.Context) error {
	cats, err := api.cs.GetCategories()
	if err != nil {
		return err
	}
	if cats == nil {
		cats = []*model.Category{}
	}
	return c.JSON(http.StatusOK, cats)
}

// swagger:operation POST /categories createCategory
// ---
// description: Создать категорию
// parameters:
// - name: category
//   in: body
//   description: новая категория
//   required: true
//   schema:
//     $ref: '#/definitions/Category'
// responses:
//  '201':
//    schema:
//      $ref: '#/definitions/Category'
//  '400':
//     description: Bad request param
//
func (api *Api) createCategory(c echo.Context) error {
	req := &model.Category{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
		return err
	}
	if err := api.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}
	res, err := api.cs.CreateCategory(req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]*int{"id": res})
}

// swagger:operation PUT /categories/{id} updateCategory
// ---
// description: Обновить категорию
// parameters:
// - name: id
//   in: path
//   description: id необходимой категории
//   required: true
//   type: int
// - name: category
//   in: body
//   description: измененная категория
//   required: true
//   schema:
//     $ref: '#/definitions/Category'
// responses:
//  '204':
//     description: Категория обновлена
//  '400':
//     description: Bad request param
//
func (api *Api) updateCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param `id`")
	}
	req := &model.Category{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}
	if err := api.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}
	req.Id = id
	if err = api.cs.UpdateCategory(req); err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			return echo.NewHTTPError(http.StatusNotFound, "Category `id` = ", id, " not found")
		}

	}
	return c.NoContent(http.StatusNoContent)
}

// swagger:operation DELETE /categories/{id} deleteCategory
// ---
// description: Удалить категорию
// parameters:
// - name: id
//   in: path
//   description: id необходимой категории
//   required: true
//   type: int
// responses:
//  '204':
//     description: Категория удалена
//  '400':
//     description: Bad request param `id`
//  '404':
//     description: Category `id`= not found
//
func (api *Api) deleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param `id`")
	}
	if err = api.cs.DeleteCategory(id); err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			return echo.NewHTTPError(http.StatusNotFound, "Category `id` = ", id, " not found")
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// swagger:operation GET /products/{id} getProduct
// ---
// description: Получить продукт
// parameters:
// - name: id
//   in: path
//   description: id необходимого продукта
//   required: true
//   type: int
// responses:
//  '200':
//    schema:
//      $ref: '#/definitions/Product'
//  '400':
//     description: Bad request param `id`
//  '404':
//     description: Product `id`= not found
//
func (api *Api) getProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product `id` = ", id, " not found")
	}
	prod, err := api.ps.GetProduct(id)
	if err != nil {
		return err
	}
	if prod == nil {
		return c.String(http.StatusNotFound, "")
	}
	return c.JSON(http.StatusOK, prod)
}

// swagger:operation GET /products getProducts
// ---
// description: Получить список продуктов
// parameters:
// - name: category
//   in: path
//   description: id категории по которой выбрать продукты
//   required: false
//   type: int
// responses:
//  '200':
//    schema:
//      type: array
//      items:
//        $ref: '#/definitions/Product'
//
func (api *Api) getProducts(c echo.Context) error {
	var products []*model.Product
	var err error
	category, err := strconv.Atoi(c.QueryParam("category"))
	if err != nil {
		products, err = api.ps.GetProducts(nil)
	} else {
		products, err = api.ps.GetProducts(&category)
	}
	if err != nil {
		return err
	}
	if products == nil {
		products = []*model.Product{}
	}
	return c.JSON(http.StatusOK, products)
}

// swagger:operation POST /products createProduct
// ---
// description: Создать продукт
// parameters:
// - name: product
//   in: body
//   description: новый продукт
//   required: true
//   schema:
//     $ref: '#/definitions/Product'
// responses:
//  '201':
//    schema:
//      $ref: '#/definitions/Product'
//  '400':
//     description: Bad request param
//
func (api *Api) createProduct(c echo.Context) error {
	req := &model.Product{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}
	if err := api.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}
	res, err := api.ps.CreateProduct(req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]*int{"id": res})
}

// swagger:operation PUT /products/{id} updateProduct
// ---
// description: Обновить продукт
// parameters:
// - name: id
//   in: path
//   description: id необходимого продукта
//   required: true
//   type: int
// - name: product
//   in: body
//   description: измененный продукт
//   required: true
//   schema:
//     $ref: '#/definitions/Product'
// responses:
//  '204':
//     description: Продукт обновлена
//  '400':
//     description: Bad request param
//
func (api *Api) updateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param `id`")
	}
	req := &model.Product{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}
	if err := api.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param: "+err.Error())
	}
	req.Id = id
	if err = api.ps.UpdateProduct(req); err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			return echo.NewHTTPError(http.StatusNotFound, "Category `id` = ", id, " not found")
		}
	}
	return c.NoContent(http.StatusNoContent)
}

// swagger:operation DELETE /products/{id} deleteProduct
// ---
// description: Удалить продукт
// parameters:
// - name: id
//   in: path
//   description: id необходимого продукта
//   required: true
//   type: int
// responses:
//  '204':
//     description: Категория удалена
//  '400':
//     description: Bad request param `id`
//  '404':
//     description: Product `id`= not found
//
func (api *Api) deleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param `id`")
	}
	if err = api.ps.DeleteProduct(id); err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			return echo.NewHTTPError(http.StatusNotFound, "Product `id` = ", id, " not found")
		}
	}

	return c.NoContent(http.StatusNoContent)
}
