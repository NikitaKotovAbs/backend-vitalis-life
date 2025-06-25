package basket

type Basket struct {
    ID          int       `json:"id"`
	ProductID int `json:"product_id"`
	UID int `json:"uid"`
	Quantity int `json:"quantity"`
}