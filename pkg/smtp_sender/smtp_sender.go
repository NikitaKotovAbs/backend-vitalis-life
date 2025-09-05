package smtp_sender

import (
	"backend/pkg/logger"
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
	"gopkg.in/gomail.v2"
)

// calculateDeliveryCost вычисляет стоимость доставки с учетом типа доставки
func calculateDeliveryCost(itemsTotal float64, deliveryType string) float64 {
	// Если это самовывоз, доставка бесплатна
	if deliveryType == "pickup" {
		return 0
	}

	// Для доставки применяем стандартные правила
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
func calculateItemsTotal(cartItems []CartItem) float64 {
	total := 0.0
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// Базовая функция отправки email
func SendEmail(to, subject, body string, isHTML bool) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		logger.Debug("не заданы SMTP параметры",
			zap.String("SMTP_HOST", smtpHost),
			zap.String("SMTP_PORT", smtpPort),
			zap.String("SMTP_USER", smtpUser),
			zap.String("SMTP_PASSWORD", smtpPass),
		)
		return fmt.Errorf("не заданы SMTP параметры")
	}

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		logger.Error("Ошибка при конвертации порта", zap.Error(err))
		return fmt.Errorf("ошибка при конвертации порта: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	if isHTML {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		logger.Error("Ошибка отправки email", zap.Error(err))
		return fmt.Errorf("ошибка отправки email: %v", err)
	}

	logger.Info("Email отправлен успешно", zap.String("to", to))
	return nil
}

type CartItem struct {
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Name      string  `json:"name"`
}

// Структура для данных заказа
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

// Функция для отправки чека клиенту
func SendReceiptToCustomer(order OrderData) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP параметры не заданы")
	}

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("неверный порт SMTP: %v", err)
	}

	// Вычисляем стоимость товаров и доставки с учетом типа доставки
	itemsTotal := calculateItemsTotal(order.CartItems)
	deliveryCost := calculateDeliveryCost(itemsTotal, order.DeliveryType)
	totalAmount := itemsTotal + deliveryCost

	// Преобразуем способ получения в читаемый формат
	deliveryText := "Самовывоз"
	if order.DeliveryType == "delivery" {
		deliveryText = "Доставка"
	}

	// Формируем таблицу товаров
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

	// Строка с доставкой
	deliveryRow := ""
	if order.DeliveryType == "pickup" {
		// Самовывоз - всегда бесплатно
		deliveryRow = fmt.Sprintf(`
        <tr>
            <td colspan="3" style="padding: 12px; border-bottom: 1px solid #e0e0e0;"><strong>%s</strong></td>
            <td style="padding: 12px; border-bottom: 1px solid #e0e0e0; text-align: right;"><strong>Бесплатно</strong></td>
        </tr>`, deliveryText)
	} else if deliveryCost > 0 {
		// Доставка с стоимостью
		deliveryRow = fmt.Sprintf(`
        <tr>
            <td colspan="3" style="padding: 12px; border-bottom: 1px solid #e0e0e0;"><strong>%s</strong></td>
            <td style="padding: 12px; border-bottom: 1px solid #e0e0e0; text-align: right;"><strong>%.2f ₽</strong></td>
        </tr>`, deliveryText, deliveryCost)
	} else {
		// Бесплатная доставка (при заказе от 5000)
		deliveryRow = fmt.Sprintf(`
        <tr>
            <td colspan="3" style="padding: 12px; border-bottom: 1px solid #e0e0e0;"><strong>%s</strong></td>
            <td style="padding: 12px; border-bottom: 1px solid #e0e0e0; text-align: right;"><strong>Бесплатно</strong></td>
        </tr>`, deliveryText)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", order.Email)
	m.SetHeader("Subject", "Ваш заказ и чек об оплате - Vitalis Life")

	// Красивый HTML шаблон чека в зеленых тонах
	htmlBody := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <style>
            body { 
                font-family: 'Arial', sans-serif; 
                line-height: 1.6; 
                color: #2d3748; 
                background: linear-gradient(135deg, #f0fff4 0%%, #e6fffa 100%%);
                margin: 0; 
                padding: 20px; 
            }
            .container { 
                max-width: 600px; 
                margin: 0 auto; 
                background: white; 
                border-radius: 16px; 
                overflow: hidden; 
                box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
            }
            .header { 
                background: linear-gradient(135deg, #38a169 0%%, #2f855a 100%%); 
                color: white; 
                padding: 30px; 
                text-align: center; 
            }
            .header h1 { 
                margin: 0; 
                font-size: 28px; 
                font-weight: bold; 
            }
            .header h2 { 
                margin: 10px 0 0 0; 
                font-size: 20px; 
                font-weight: 300; 
                opacity: 0.9; 
            }
            .content { 
                padding: 30px; 
            }
            .section { 
                margin-bottom: 25px; 
                padding: 20px; 
                background: #f7fafc; 
                border-radius: 12px; 
                border-left: 4px solid #38a169;
            }
            .section-title { 
                color: #2d3748; 
                margin-bottom: 15px; 
                font-size: 18px; 
                font-weight: 600; 
            }
            table { 
                width: 100%%; 
                border-collapse: collapse; 
                margin: 15px 0; 
            }
            th { 
                background: #e6fffa; 
                padding: 12px; 
                text-align: left; 
                font-weight: 600; 
                color: #2d3748;
                border-bottom: 2px solid #38a169;
            }
            td { 
                padding: 12px; 
                border-bottom: 1px solid #e2e8f0; 
            }
            .total-row { 
                font-weight: bold; 
                background: #f0fff4; 
                font-size: 16px; 
                border-top: 2px solid #38a169;
            }
            .info-grid { 
                display: grid; 
                grid-template-columns: 120px 1fr; 
                gap: 10px; 
                margin: 10px 0; 
            }
            .info-label { 
                font-weight: 600; 
                color: #4a5568; 
            }
            .footer { 
                text-align: center; 
                padding: 25px; 
                background: #f7fafc; 
                color: #718096; 
                font-size: 14px; 
                border-top: 1px solid #e2e8f0;
            }
            .logo { 
                font-size: 20px; 
                font-weight: bold; 
                color: #38a169; 
                margin-bottom: 10px; 
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Vitalis Life</h1>
                <h2>Заказ успешно оплачен!</h2>
            </div>
            
            <div class="content">
                <div class="section">
                    <div class="section-title">Детали заказа</div>
                    <div class="info-grid">
                        <div class="info-label">ID заказа:</div>
                        <div>%s</div>
                        <div class="info-label">Клиент:</div>
                        <div>%s</div>
                        <div class="info-label">Телефон:</div>
                        <div>%s</div>
                        <div class="info-label">Способ:</div>
                        <div>%s</div>
                        <div class="info-label">Адрес:</div>
                        <div>%s</div>
                        <div class="info-label">Комментарий:</div>
                        <div>%s</div>
                    </div>
                </div>

                <div class="section">
                    <div class="section-title">Состав заказа</div>
                    <table>
                        <thead>
                            <tr>
                                <th>Товар</th>
                                <th style="text-align: center;">Кол-во</th>
                                <th style="text-align: right;">Цена</th>
                                <th style="text-align: right;">Сумма</th>
                            </tr>
                        </thead>
                        <tbody>
                            %s
                            %s
                            <tr class="total-row">
                                <td colspan="3" style="text-align: right;"><strong>Итого к оплате:</strong></td>
                                <td style="text-align: right;"><strong>%.2f ₽</strong></td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            
            <div class="footer">
                <div class="logo">Vitalis Life</div>
                <p>Спасибо за ваш заказ! Мы свяжемся с вами в ближайшее время.</p>
                <p>Если у вас есть вопросы: support@vitalis-life.ru</p>
                <p>© 2024 Vitalis Life. Все права защищены.</p>
            </div>
        </div>
    </body>
    </html>
    `, order.PaymentID, order.CustomerName, order.Phone, deliveryText,
		order.DeliveryAddress, order.Comment, itemsTable, deliveryRow, totalAmount)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		logger.Error("Ошибка при отправке чека", zap.Error(err))
		return fmt.Errorf("ошибка при отправке чека: %v", err)
	}

	logger.Info("Чек отправлен успешно", zap.String("email", order.Email))
	return nil
}

// Функция для отправки заказа менеджеру
func SendOrderToManager(order OrderData, managerEmail string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP параметры не заданы")
	}

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("неверный порт SMTP: %v", err)
	}

	// Вычисляем стоимость товаров и доставки с учетом типа доставки
	itemsTotal := calculateItemsTotal(order.CartItems)
	deliveryCost := calculateDeliveryCost(itemsTotal, order.DeliveryType)
	totalAmount := itemsTotal + deliveryCost

	// Преобразуем способ получения в читаемый формат
	deliveryText := "Самовывоз"
	if order.DeliveryType == "delivery" {
		deliveryText = "Доставка"
	}

	// Формируем список товаров для менеджера
	itemsList := ""
	for _, item := range order.CartItems {
		itemsList += fmt.Sprintf("• %s: %d шт. x %.2f ₽ = %.2f ₽\n",
			item.Name, item.Quantity, item.Price, item.Price*float64(item.Quantity))
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", managerEmail)
	m.SetHeader("Subject", fmt.Sprintf("🎯 Новый заказ #%s - %.2f ₽", order.PaymentID, totalAmount))

	// Красивый HTML для менеджера
	htmlBody := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <style>
            body { 
                font-family: 'Arial', sans-serif; 
                line-height: 1.6; 
                color: #2d3748; 
                background: linear-gradient(135deg, #f0fff4 0%%, #e6fffa 100%%);
                margin: 0; 
                padding: 20px; 
            }
            .container { 
                max-width: 600px; 
                margin: 0 auto; 
                background: white; 
                border-radius: 16px; 
                overflow: hidden; 
                box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
            }
            .header { 
                background: linear-gradient(135deg, #e53e3e 0%%, #c53030 100%%); 
                color: white; 
                padding: 25px; 
                text-align: center; 
            }
            .header h1 { 
                margin: 0; 
                font-size: 24px; 
                font-weight: bold; 
            }
            .content { 
                padding: 25px; 
            }
            .section { 
                margin-bottom: 20px; 
                padding: 20px; 
                background: #fff5f5; 
                border-radius: 12px; 
                border-left: 4px solid #e53e3e;
            }
            .info-grid { 
                display: grid; 
                grid-template-columns: 100px 1fr; 
                gap: 8px; 
                margin: 10px 0; 
            }
            .info-label { 
                font-weight: 600; 
                color: #4a5568; 
            }
            .items-list { 
                background: #f7fafc; 
                padding: 15px; 
                border-radius: 8px; 
                margin: 10px 0; 
            }
            .total { 
                font-weight: bold; 
                font-size: 18px; 
                color: #2d3748; 
                margin-top: 15px;
                padding-top: 15px;
                border-top: 2px solid #e53e3e;
            }
            .footer { 
                text-align: center; 
                padding: 20px; 
                background: #f7fafc; 
                color: #718096; 
                font-size: 14px; 
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>🚨 НОВЫЙ ЗАКАЗ!</h1>
            </div>
            
            <div class="content">
                <div class="section">
                    <div class="info-grid">
                        <div class="info-label">ID:</div>
                        <div><strong>%s</strong></div>
                        <div class="info-label">Клиент:</div>
                        <div>%s</div>
                        <div class="info-label">Телефон:</div>
                        <div>%s</div>
                        <div class="info-label">Email:</div>
                        <div>%s</div>
                        <div class="info-label">Способ:</div>
                        <div>%s</div>
                        <div class="info-label">Адрес:</div>
                        <div>%s</div>
                        <div class="info-label">Комментарий:</div>
                        <div>%s</div>
                    </div>
                </div>

                <div class="section">
                    <h3>📦 Товары:</h3>
                    <div class="items-list">
                        <pre style="margin: 0; font-family: Arial; line-height: 1.4;">%s</pre>
                    </div>
                    
                    <div class="total">
                        💰 Итого: <strong>%.2f ₽</strong><br>
                        📦 Товары: %.2f ₽<br>
                        🚚 %s: %.2f ₽
                    </div>
                </div>
            </div>
            
            <div class="footer">
                <p>💡 Требуется срочная обработка!</p>
            </div>
        </div>
    </body>
    </html>
    `, order.PaymentID, order.CustomerName, order.Phone, order.Email,
		deliveryText, order.DeliveryAddress, order.Comment, itemsList,
		totalAmount, itemsTotal, deliveryText, deliveryCost)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		logger.Error("Ошибка при отправке заказа менеджеру", zap.Error(err))
		return fmt.Errorf("ошибка при отправке заказа менеджеру: %v", err)
	}

	logger.Info("Заказ отправлен менеджеру", zap.String("manager_email", managerEmail))
	return nil
}

// Универсальная функция для отправки обоих писем
func SendOrderEmails(order OrderData, managerEmail string) error {
	// Отправляем чек клиенту
	if err := SendReceiptToCustomer(order); err != nil {
		logger.Error("Ошибка отправки чека клиенту", zap.Error(err))
		return err
	}

	// Отправляем заказ менеджеру
	if err := SendOrderToManager(order, managerEmail); err != nil {
		logger.Error("Ошибка отправки заказа менеджеру", zap.Error(err))
		return err
	}

	return nil
}