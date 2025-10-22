package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendResetEmail(to, token string) error {
	from := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), token)
	subject := "Subject: Reset Password Request\r\n"
	body := fmt.Sprintf("Click the link below to reset your password:\r\n%s\r\n", resetLink)
	message := []byte(subject + "\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}
