package model

// Продукт.
// Сущность продукт
// swagger:model
type Product struct {
	// id продукта
	Id          int     `json:"id"`
	// название
	Name        string  `json:"name" validate:"required,min=3"`
	// описание
	Description string  `json:"desc"`
	// id категория
	Category    int     `json:"category" validate:"required"`
	// цена
	Price       float64 `json:"price" validate:"required,gt=0"`
}
