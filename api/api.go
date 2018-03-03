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
	mas := map[string]string{"msg": "ok"}
	return c.JSON(http.StatusOK, mas)
}

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

func (api *Api) deleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request param `id`")
	}
	if err = api.ps.DeleteProduct(id); err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			return echo.NewHTTPError(http.StatusNotFound, "Category `id` = ", id, " not found")
		}
	}

	return c.NoContent(http.StatusNoContent)
}
