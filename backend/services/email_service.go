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

// SendRegistrationEmail は仮登録メールを送信します。
func (s *EmailService) SendRegistrationEmail(to string, verificationToken string) error {
	from := "info@login-go.app" // 送信元アドレス

	// BACKEND_URL を環境変数から取得（未設定ならデフォルト値）
	backendURL := os.Getenv("BACKEND_URL")

	// SMTP_HOST, SMTP_PORT を環境変数から取得（未設定ならデフォルト値）
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	subject := "ReDesigner for Student 仮登録"
	// verificationLink を BACKEND_URL から動的に組み立て
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

	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, message)
}

// SendPasswordResetEmail はパスワードリセットの案内メールを送信します。
func (s *EmailService) SendPasswordResetEmail(to string, resetToken string) error {
	from := "info@login-go.app"

	// FRONTEND_URL を環境変数から取得（未設定ならデフォルト値）
	frontendURL := os.Getenv("FRONTEND_URL")

	// SMTP_HOST, SMTP_PORT を環境変数から取得
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	subject := "パスワードリセットのご案内"
	// リセットリンクを FRONTEND_URL から組み立てる
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

	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, message)
}

// SendWelcomeEmail は歓迎メールを送信します。
func (s *EmailService) SendWelcomeEmail(to string) error {
	from := "info@login-go.app"

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	subject := "ReDesigner for Student へようこそ！"
	body := `
    <html>
    <body>
        <div style="font-family: Arial, sans-serif; color: #333;">
            <h2>ReDesigner for Student へようこそ！</h2>
            <p>この度はご登録ありがとうございます！ReDesigner for Studentは「デザイナーを目指す、すべての学生のためのプラットフォーム」です。</p>
            <p>デザイナーを大切にしている企業のインターンや本採用の情報を集めたり、他大学の学生ポートフォリオを見たりと、デザインを学ぶことから本格的な就職活動をすることまで、さまざまなシーンでご利用いただけます！</p>
            <p>プロフィールや作品をアップすることによって、企業の人たちがあなたを見つけることができるようになります。ぜひあなたのことを、あなたの作品や言葉で教えてくださいね。</p>
            <p>作品をきっかけにあなたと企業がつながり、そこから得たフィードバックがあなたを支える。そんな出会いをお届けできたら、私たちも嬉しいです。</p>
        </div>
    </body>
    </html>
    `

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, message)
}

// SendPasswordResetConfirmationEmail はパスワード変更完了のお知らせメールを送信します。
func (s *EmailService) SendPasswordResetConfirmationEmail(to string) error {
	from := "info@login-go.app"

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

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

	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, message)
}
