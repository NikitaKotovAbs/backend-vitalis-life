package basket

type BasketRepository interface {
	Add(data Basket) error
}