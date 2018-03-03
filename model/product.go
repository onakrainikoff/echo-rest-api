package model

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" validate:"required,min=3"`
	Description string  `json:"desc"`
	Category    int     `json:"category" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}
