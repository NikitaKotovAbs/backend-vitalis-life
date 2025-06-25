package product

import "time"

type Product struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Price       float64   `json:"price"`
    Description string    `json:"description"`
    Discount    float64   `json:"discount"`
    Image       string    `json:"image,omitempty"` // omitempty - не показывать если nil
    CreatedAt   time.Time `json:"created_at"`
}