package services

import (
	"fmt"
	"net/smtp"
)

type IEmailService interface {
	SendRegistrationEmail(to string, verificationLink string) error
}

type EmailService struct{}

// メール送信のための関数を定義
func NewEmailService() IEmailService {
	return &EmailService{}
}

// メール送信機能の実装
func (s *EmailService) SendRegistrationEmail(to string, verificationToken string) error {
	from := "info@login-go.app" // デフォルトの送信元アドレス
	smtpHost := "localhost"     // MailHogのホスト
	smtpPort := "1025"          // MailHogのデフォルトSMTPポート

	// メールの件名と内容をHTML形式で作成
	subject := "ReDesigner for Student 仮登録"
	body := fmt.Sprintf(`
    <html>
    <body>
        <div style="font-family: Arial, sans-serif; color: #333;">
            <h2 style="color: #F15A24;">ReDesigner for Student</h2>
            <p>こんにちは、</p>
            <p>ReDesigner for Studentへの仮登録を受け付けました。</p>
            <p>下記のボタンをクリックして、本登録を完了させてください。</p>
            <a href="http://localhost:8080/auth/verify?token=%s" style="padding: 10px 20px; background-color: #F15A24; color: #fff; text-decoration: none; border-radius: 5px;">本登録を完了する</a>
            <p>このリンクの有効期限は<strong>7日間</strong>です。</p>
        </div>
    </body>
    </html>`, verificationToken)

	// メールヘッダーと本文の組み立て
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	// メール送信
	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, message)
}
