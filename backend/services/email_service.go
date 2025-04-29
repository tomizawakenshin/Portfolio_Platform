package services

import (
	"fmt"
	"net/smtp"
	"os"
)

// IEmailService はメール送信機能のインターフェースです。
type IEmailService interface {
	SendRegistrationEmail(to string, verificationToken string) error
	SendPasswordResetEmail(to string, resetToken string) error
	SendWelcomeEmail(to string) error
	SendPasswordResetConfirmationEmail(to string) error
}

// EmailService は IEmailService の実装です。
type EmailService struct{}

// NewEmailService は EmailService の新しいインスタンスを返します。
func NewEmailService() IEmailService {
	return &EmailService{}
}

// SMTP 設定を環境変数から取得するヘルパー
func getSMTPConfig() (host, port string, auth smtp.Auth) {
	env := os.Getenv("ENV")
	if env != "prod" {
		// 開発環境 → MailHog
		host = os.Getenv("MAILHOG_HOST")
		port = os.Getenv("MAILHOG_PORT")
		auth = nil // MailHog は認証不要
	} else {
		// 本番環境 → SMTP(Gmail 等)
		host = os.Getenv("SMTP_HOST")
		port = os.Getenv("SMTP_PORT")

		username := os.Getenv("SMTP_USERNAME")
		password := os.Getenv("SMTP_PASSWORD")
		// 認証情報を作成（コメントアウトせず残しておく）
		auth = smtp.PlainAuth("", username, password, host)
	}
	return
}

func (s *EmailService) SendRegistrationEmail(to string, verificationToken string) error {
	frontendURL := os.Getenv("FRONTEND_URL")
	from := os.Getenv("SMTP_USERNAME")

	subject := "エンジニアのポートフォリオ 仮登録"
	verificationLink := fmt.Sprintf("%s/verifyStart?token=%s", frontendURL, verificationToken)
	body := fmt.Sprintf(`
    <html>
    <body>
        <div style="font-family: Arial, sans-serif; color: #333;">
            <h2 style="color: #F15A24;">エンジニアのポートフォリオ</h2>
            <p>こんにちは、</p>
            <p>エンジニアのポートフォリオへの仮登録を受け付けました</p>
            <p>下記のボタンをクリックして本登録を完了させてください。</p>
            <a href="%s" style="padding: 10px 20px; background-color: #F15A24; color: #fff; text-decoration: none; border-radius: 5px;">本登録を完了する</a>
            <p>このリンクの有効期限は<strong>7日間</strong>です。</p>
        </div>
    </body>
    </html>`, verificationLink)

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	host, port, auth := getSMTPConfig()
	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, from, []string{to}, message)
}

// SendPasswordResetEmail はパスワードリセットの案内メールを送信します。
func (s *EmailService) SendPasswordResetEmail(to string, resetToken string) error {
	frontendURL := os.Getenv("FRONTEND_URL")
	from := os.Getenv("SMTP_USERNAME")

	subject := "パスワードリセットのご案内"
	resetLink := fmt.Sprintf("%s/PasswordReset/%s", frontendURL, resetToken)
	body := fmt.Sprintf(`
    <html>
    <body>
        <div style="font-family: Arial, sans-serif; color: #333;">
            <h2>パスワードリセット</h2>
            <p>パスワードリセットのリクエストを受け付けました。</p>
            <p>以下のリンクをクリックしてパスワードをリセットしてください：</p>
            <a href="%s" style="padding: 10px 20px; background-color: #F15A24; color: #fff; text-decoration: none; border-radius: 5px;">パスワードをリセットする</a>
            <p>このリンクの有効期限は1時間です。</p>
        </div>
    </body>
    </html>`, resetLink)

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	host, port, auth := getSMTPConfig()
	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, from, []string{to}, message)
}

func (s *EmailService) SendWelcomeEmail(to string) error {
	from := os.Getenv("SMTP_USERNAME")

	subject := "エンジニアのポートフォリオ へようこそ！"
	body := `
    <html>
    <body>
        <div style="font-family: Arial, sans-serif; color: #333;">
            <h2>エンジニアのポートフォリオ へようこそ！</h2>
            <p>この度はご登録ありがとうございます！エンジニアのポートフォリオは「エンジニアを目指す、すべての学生のためのプラットフォーム」です。</p>
            <p>他大学の学生ポートフォリオを見ることで本格的な就職活動でご利用いただけます！</p>
            <p>プロフィールや作品をアップすることによって、企業の人たちがあなたを見つけることができるようになるかもしれません。ぜひあなたのことを、あなたの作品や言葉で教えてくださいね。</p>
        </div>
    </body>
    </html>
    `

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	host, port, auth := getSMTPConfig()
	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, from, []string{to}, message)
}

// SendPasswordResetConfirmationEmail はパスワード変更完了のお知らせメールを送信します。
func (s *EmailService) SendPasswordResetConfirmationEmail(to string) error {
	from := os.Getenv("SMTP_USERNAME")

	subject := "パスワード変更完了のお知らせ"
	body := fmt.Sprintf(`
    <html>
    <body>
        <div style="font-family: Arial, sans-serif; color: #333;">
            <p>こんにちは、%s さん</p>
            <p>パスワード変更が完了しました。</p>
            <p>変更したパスワードはセキュリティの関係上、記載しておりません。<br>
            ログインIDやパスワードはサービス利用にあたり重要な情報のため、<br>
            ご自身で大切に保管していただきますようお願い致します。</p>
            <hr>
        </div>
    </body>
    </html>
    `, to)

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	host, port, auth := getSMTPConfig()
	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, from, []string{to}, message)
}
