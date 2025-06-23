package smtp_sender

import (
	"backend/pkg/logger"
	"crypto/tls"
	"fmt"
	"math/rand"
	"go.uber.org/zap"
	"os"
	"strconv"
	"gopkg.in/gomail.v2"
)

func GenerateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendEmail(to string, code string) error {
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
		logger.Error("Ошибка при попытке конвертировать строку в число",
		zap.Error(err),
		)
		return fmt.Errorf("ошибка при попытке конвертировать строку в число")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Код подтверждения")
	m.SetBody("text/plain", fmt.Sprintf("Ваш код подтверждения: %s", code)) // Добавлен параметр code
	
	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("ошибка при отправке email: %v", err)
	}

	return nil
}
