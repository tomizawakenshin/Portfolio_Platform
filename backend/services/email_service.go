package services

import (
	"fmt"
	"net/smtp"
)

type IEmailService interface {
	SendRegistrationEmail(to string, verificationLink string) error
	SendPasswordResetEmail(to string, resetToken string) error
	SendWelcomeEmail(to string) error
	SendPasswordResetConfirmationEmail(to string) error
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

func (s *EmailService) SendPasswordResetEmail(to string, resetToken string) error {
	from := "info@login-go.app"
	smtpHost := "localhost"
	smtpPort := "1025"

	subject := "パスワードリセットのご案内"
	// リセットリンクを構築
	resetLink := fmt.Sprintf("http://localhost:3000/PasswordReset/%s", resetToken)

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

	// メールヘッダーとメッセージ
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	// メール送信
	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, message)
}

func (s *EmailService) SendWelcomeEmail(to string) error {
	from := "info@login-go.app"
	smtpHost := "localhost"
	smtpPort := "1025"

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

	// メールヘッダーと本文の組み立て
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	// メール送信
	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, message)
}

func (s *EmailService) SendPasswordResetConfirmationEmail(to string) error {
	from := "info@login-go.app"
	smtpHost := "localhost"
	smtpPort := "1025"

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

	// メールヘッダーと本文の組み立て
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	// メール送信
	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, message)
}
