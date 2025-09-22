package smtp_sender

import (
	"backend/pkg/logger"
	"backend/pkg/templates"
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

// SMTPConfig конфигурация SMTP
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

// GetSMTPConfig возвращает конфигурацию SMTP из переменных окружения
func GetSMTPConfig() (*SMTPConfig, error) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return nil, fmt.Errorf("SMTP параметры не заданы")
	}

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return nil, fmt.Errorf("неверный порт SMTP: %v", err)
	}

	return &SMTPConfig{
		Host:     smtpHost,
		Port:     port,
		User:     smtpUser,
		Password: smtpPass,
	}, nil
}

// createDialer создает SMTP dialer
func createDialer(config *SMTPConfig) *gomail.Dialer {
	d := gomail.NewDialer(config.Host, config.Port, config.User, config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d
}

// SendEmail базовая функция отправки email
func SendEmail(to, subject, body string, isHTML bool) error {
	config, err := GetSMTPConfig()
	if err != nil {
		logger.Debug("не заданы SMTP параметры", zap.Error(err))
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.User)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	if isHTML {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	d := createDialer(config)
	if err := d.DialAndSend(m); err != nil {
		logger.Error("Ошибка отправки email", zap.Error(err))
		return fmt.Errorf("ошибка отправки email: %v", err)
	}

	logger.Info("Email отправлен успешно", zap.String("to", to))
	return nil
}

// SendReceiptToCustomer отправляет чек клиенту
func SendReceiptToCustomer(order templates.OrderData) error {
	config, err := GetSMTPConfig()
	if err != nil {
		return err
	}

	htmlBody := templates.GenerateReceiptHTML(order)

	m := gomail.NewMessage()
	m.SetHeader("From", config.User)
	m.SetHeader("To", order.Email)
	m.SetHeader("Subject", "Ваш заказ и чек об оплате - Vitalis Life")
	m.SetBody("text/html", htmlBody)

	d := createDialer(config)
	if err := d.DialAndSend(m); err != nil {
		logger.Error("Ошибка при отправке чека", zap.Error(err))
		return fmt.Errorf("ошибка при отправке чека: %v", err)
	}

	logger.Info("Чек отправлен успешно", zap.String("email", order.Email))
	return nil
}

// SendOrderToManager отправляет заказ менеджеру
func SendOrderToManager(order templates.OrderData, managerEmail string) error {
	config, err := GetSMTPConfig()
	if err != nil {
		return err
	}

	htmlBody := templates.GenerateManagerOrderHTML(order)

	m := gomail.NewMessage()
	m.SetHeader("From", config.User)
	m.SetHeader("To", managerEmail)
	m.SetHeader("Subject", fmt.Sprintf("🎯 Новый заказ #%s", order.PaymentID))
	m.SetBody("text/html", htmlBody)

	d := createDialer(config)
	if err := d.DialAndSend(m); err != nil {
		logger.Error("Ошибка при отправке заказа менеджеру", zap.Error(err))
		return fmt.Errorf("ошибка при отправке заказа менеджеру: %v", err)
	}

	logger.Info("Заказ отправлен менеджеру", zap.String("manager_email", managerEmail))
	return nil
}

// SendOrderEmails универсальная функция для отправки обоих писем
func SendOrderEmails(order templates.OrderData, managerEmail string) error {
	if err := SendReceiptToCustomer(order); err != nil {
		logger.Error("Ошибка отправки чека клиенту", zap.Error(err))
		return err
	}

	if err := SendOrderToManager(order, managerEmail); err != nil {
		logger.Error("Ошибка отправки заказа менеджеру", zap.Error(err))
		return err
	}

	return nil
}