package model
// Категория.
// Сущность категория продукта
// swagger:model
type Category struct {
	// id категории
	Id   int    `json:"id"`
	// название категории
	Name string `json:"name" validate:"required,min=3"`
}
