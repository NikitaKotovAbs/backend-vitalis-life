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
	CartItems       string // JSON строка с товарами
}

// Функция для отправки чека клиенту
func SendReceiptToCustomer(order OrderData) error {
	subject := "Ваш заказ и чек об оплате - Vitalis Life"

	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background: #4CAF50; color: white; padding: 20px; text-align: center; border-radius: 5px; }
			.content { background: #f9f9f9; padding: 20px; margin: 20px 0; border-radius: 5px; }
			.footer { text-align: center; color: #666; font-size: 12px; margin-top: 20px; }
			.info-item { margin: 10px 0; }
			.info-label { font-weight: bold; color: #555; }
		</style>
	</head>
	<body>
		<div class="header">
			<h1>Vitalis Life</h1>
			<h2>Заказ успешно оплачен!</h2>
		</div>
		
		<div class="content">
			<p>Уважаемый(ая) <strong>%s</strong>,</p>
			<p>Ваш заказ успешно оплачен. Ниже приведены детали заказа:</p>
			
			<div class="info-item">
				<span class="info-label">ID платежа:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Сумма:</span> %s %s
			</div>
			<div class="info-item">
				<span class="info-label">Описание:</span> %s
			</div>
			
			<h3>Контактная информация:</h3>
			<div class="info-item">
				<span class="info-label">Email:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Телефон:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Способ получения:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Адрес:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Комментарий:</span> %s
			</div>
		</div>
		
		<div class="footer">
			<p>С уважением,<br>Команда Vitalis Life</p>
			<p>Если у вас есть вопросы, свяжитесь с нами: support@vitalis-life.ru</p>
		</div>
	</body>
	</html>
	`, order.CustomerName, order.PaymentID, order.Amount, order.Currency, 
	   order.Description, order.Email, order.Phone, order.DeliveryType, 
	   order.DeliveryAddress, order.Comment)

	return SendEmail(order.Email, subject, htmlBody, true)
}

// Функция для отправки заказа менеджеру
func SendOrderToManager(order OrderData, managerEmail string) error {
	subject := fmt.Sprintf("Новый заказ #%s - Vitalis Life", order.PaymentID)

	// Парсим товары из JSON (упрощенно)
	var cartItemsInfo string
	if order.CartItems != "" {
		cartItemsInfo = "Товары:\n" + order.CartItems
	} else {
		cartItemsInfo = "Информация о товарах не указана"
	}

	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background: #ff6b35; color: white; padding: 20px; text-align: center; border-radius: 5px; }
			.content { background: #fff3cd; padding: 20px; margin: 20px 0; border-radius: 5px; border: 1px solid #ffeaa7; }
			.urgent { color: #d63031; font-weight: bold; }
			.info-item { margin: 10px 0; padding: 8px; background: #f8f9fa; border-radius: 3px; }
			.info-label { font-weight: bold; color: #2d3436; }
		</style>
	</head>
	<body>
		<div class="header">
			<h1>Vitalis Life</h1>
			<h2>НОВЫЙ ЗАКАЗ!</h2>
		</div>
		
		<div class="content">
			<p class="urgent">💡 Требуется обработка заказа</p>
			
			<div class="info-item">
				<span class="info-label">ID заказа:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Сумма:</span> %s %s
			</div>
			<div class="info-item">
				<span class="info-label">Клиент:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Email:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Телефон:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Способ получения:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Адрес:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Комментарий клиента:</span> %s
			</div>
			<div class="info-item">
				<span class="info-label">Товары:</span><br>
				<pre>%s</pre>
			</div>
		</div>
		
		<p><strong>Время сбора заказа:</strong> Как можно скорее!</p>
	</body>
	</html>
	`, order.PaymentID, order.Amount, order.Currency, order.CustomerName, 
	   order.Email, order.Phone, order.DeliveryType, order.DeliveryAddress, 
	   order.Comment, cartItemsInfo)

	return SendEmail(managerEmail, subject, htmlBody, true)
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