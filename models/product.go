package models

import "time"

type Product struct {
	Id            int       `json:"id"`
	Sku           string    `json:"sku"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Etalase       string    `json:"etalase"`
	Image         string    `json:"image"`
	Weight        float32   `json:"weight"`
	Price         float32   `json:"price"`
	CreatedAt     time.Time `json:"created_at"`
	AverageRating float32   `json:"average_rating"`
}
