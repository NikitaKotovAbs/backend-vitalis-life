package templates

import "fmt"

// calculateDeliveryCost вычисляет стоимость доставки
func CalculateDeliveryCost(itemsTotal float64, deliveryType string) float64 {
	if deliveryType == "pickup" {
		return 0
	}

	switch {
	case itemsTotal >= 5000:
		return 0
	case itemsTotal >= 3000:
		return itemsTotal * 0.10
	case itemsTotal >= 1000:
		return itemsTotal * 0.15
	default:
		return itemsTotal * 0.20
	}
}

// calculateItemsTotal вычисляет общую стоимость товаров
func CalculateItemsTotal(cartItems []CartItem) float64 {
	total := 0.0
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

type CartItem struct {
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Name      string  `json:"name"`
}

type OrderData struct {
	CustomerName    string
	Email           string
	Phone           string
	DeliveryType    string
	DeliveryAddress string
	Comment         string
	PaymentID       string
	Amount          string
	Currency        string
	Description     string
	CartItems       []CartItem
}

// GenerateReceiptHTML генерирует HTML для чека клиента
func GenerateReceiptHTML(order OrderData) string {
	itemsTotal := CalculateItemsTotal(order.CartItems)
	deliveryCost := CalculateDeliveryCost(itemsTotal, order.DeliveryType)
	totalAmount := itemsTotal + deliveryCost

	deliveryText := "Самовывоз"
	if order.DeliveryType == "delivery" {
		deliveryText = "Доставка"
	}

	itemsTable := ""
	for _, item := range order.CartItems {
		itemTotal := item.Price * float64(item.Quantity)
		itemsTable += fmt.Sprintf(`
        <tr>
            <td style="padding: 12px; border-bottom: 1px solid #e0e0e0;">%s</td>
            <td style="padding: 12px; border-bottom: 1px solid #e0e0e0; text-align: center;">%d шт.</td>
            <td style="padding: 12px; border-bottom: 1px solid #e0e0e0; text-align: right;">%.2f ₽</td>
            <td style="padding: 12px; border-bottom: 1px solid #e0e0e0; text-align: right;">%.2f ₽</td>
        </tr>`, item.Name, item.Quantity, item.Price, itemTotal)
	}

	deliveryRow := ""
	if order.DeliveryType == "pickup" {
		deliveryRow = fmt.Sprintf(`
        <tr>
            <td colspan="3" style="padding: 12px; border-bottom: 1px solid #e0e0e0;"><strong>%s</strong></td>
            <td style="padding: 12px; border-bottom: 1px solid #e0e0e0; text-align: right;"><strong>Бесплатно</strong></td>
        </tr>`, deliveryText)
	} else if deliveryCost > 0 {
		deliveryRow = fmt.Sprintf(`
        <tr>
            <td colspan="3" style="padding: 12px; border-bottom: 1px solid #e0e0e0;"><strong>%s</strong></td>
            <td style="padding: 12px; border-bottom: 1px solid #e0e0e0; text-align: right;"><strong>%.2f ₽</strong></td>
        </tr>`, deliveryText, deliveryCost)
	} else {
		deliveryRow = fmt.Sprintf(`
        <tr>
            <td colspan="3" style="padding: 12px; border-bottom: 1px solid #e0e0e0;"><strong>%s</strong></td>
            <td style="padding: 12px; border-bottom: 1px solid #e0e0e0; text-align: right;"><strong>Бесплатно</strong></td>
        </tr>`, deliveryText)
	}

	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <style>/* стили остаются такими же */</style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Vitalis Life</h1>
                <h2>Заказ успешно оплачен!</h2>
            </div>
            
            <div class="content">
                <!-- остальная верстка -->
                %s
                %s
                <tr class="total-row">
                    <td colspan="3" style="text-align: right;"><strong>Итого к оплате:</strong></td>
                    <td style="text-align: right;"><strong>%.2f ₽</strong></td>
                </tr>
            </div>
        </div>
    </body>
    </html>
    `, order.PaymentID, order.CustomerName, order.Phone, deliveryText,
		order.DeliveryAddress, order.Comment, itemsTable, deliveryRow, totalAmount)
}

// GenerateManagerOrderHTML генерирует HTML для уведомления менеджера
func GenerateManagerOrderHTML(order OrderData) string {
	itemsTotal := CalculateItemsTotal(order.CartItems)
	deliveryCost := CalculateDeliveryCost(itemsTotal, order.DeliveryType)
	totalAmount := itemsTotal + deliveryCost

	deliveryText := "Самовывоз"
	if order.DeliveryType == "delivery" {
		deliveryText = "Доставка"
	}

	itemsList := ""
	for _, item := range order.CartItems {
		itemsList += fmt.Sprintf("• %s: %d шт. x %.2f ₽ = %.2f ₽\n",
			item.Name, item.Quantity, item.Price, item.Price*float64(item.Quantity))
	}

	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <style>/* стили для менеджера */</style>
    </head>
    <body>
        <!-- верстка для менеджера -->
        %s
        %.2f
    </body>
    </html>
    `, order.PaymentID, order.CustomerName, order.Phone, order.Email,
		deliveryText, order.DeliveryAddress, order.Comment, itemsList,
		totalAmount, itemsTotal, deliveryText, deliveryCost)
}