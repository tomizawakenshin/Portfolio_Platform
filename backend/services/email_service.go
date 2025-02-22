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

func (s *EmailService) SendRegistrationEmail(to string, verificationToken string) error {
	// 環境変数から必要な情報を取得
	backendURL := os.Getenv("BACKEND_URL")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// 送信元は SMTP_USERNAME を利用
	from := smtpUsername

	subject := "ReDesigner for Student 仮登録"
	verificationLink := fmt.Sprintf("%s/auth/verify?token=%s", backendURL, verificationToken)
	body := fmt.Sprintf(`
    <html>
    <body>
        <div style="font-family: Arial, sans-serif; color: #333;">
            <h2 style="color: #F15A24;">ReDesigner for Student</h2>
            <p>こんにちは、</p>
            <p>ReDesigner for Studentへの仮登録を受け付けました。</p>
            <p>下記のボタンをクリックして、本登録を完了させてください。</p>
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

	// SMTP認証情報の設定
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// SMTPサーバーに接続してメール送信
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}

// SendPasswordResetEmail はパスワードリセットの案内メールを送信します。
func (s *EmailService) SendPasswordResetEmail(to string, resetToken string) error {
	// 環境変数から必要な情報を取得
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	frontendURL := os.Getenv("FRONTEND_URL")

	// 送信元は SMTP_USERNAME を利用
	from := smtpUsername

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

	// SMTP認証情報の設定
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}

func (s *EmailService) SendWelcomeEmail(to string) error {
	// 環境変数から必要な情報を取得
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// 送信元は SMTP_USERNAME を利用
	from := smtpUsername

	subject := "ReDesigner for Student へようこそ！"
	body := `
    <html>
    <body>
        <div style="font-family: Arial, sans-serif; color: #333;">
            <h2>ReDesigner for Student へようこそ！</h2>
            <p>この度はご登録ありがとうございます！ReDesigner for Studentは「デザイナーを目指す、すべての学生のためのプラットフォーム」です。</p>
            <p>デザイナーを大切にしている企業のインターンや本採用の情報を集めたり、他大学の学生ポートフォリオを見たりと、デザインを学ぶことから本格的な就職活動をすることまで、さまざまなシーンでご利用いただけます！</p>
            <p>プロフィールや作品をアップすることによって、企業の人たちがあなたを見つけることができるようになります。ぜひあなたのことを、あなたの作品や言葉で教えてくださいね。</p>
        </div>
    </body>
    </html>
    `

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	// SMTP認証情報の設定
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}

// SendPasswordResetConfirmationEmail はパスワード変更完了のお知らせメールを送信します。
func (s *EmailService) SendPasswordResetConfirmationEmail(to string) error {
	// 環境変数から必要な情報を取得
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// 送信元は SMTP_USERNAME を利用
	from := smtpUsername

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

	// SMTP認証情報の設定
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}
