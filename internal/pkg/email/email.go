package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"chaoxing/internal/globals/config"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

var defaultConfig EmailConfig

func Init() {
	defaultConfig = EmailConfig{
		Host:     config.Config.GetString("email.host"),
		Port:     config.Config.GetInt("email.port"),
		Username: config.Config.GetString("email.username"),
		Password: config.Config.GetString("email.password"),
		From:     config.Config.GetString("email.from"),
	}
}

func SendVerificationCode(to string, code string) error {
	subject := "超星自动签到 - 邮箱验证"
	body := fmt.Sprintf(`
		<h1>邮箱验证码</h1>
		<p>您的验证码是: <strong>%s</strong></p>
		<p>验证码有效期为5分钟，请及时使用。</p>
		<p>如果这不是您的操作，请忽略此邮件。</p>
	`, code)

	return sendEmail(to, subject, body)
}

func sendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", defaultConfig.Host, defaultConfig.Port)
	auth := smtp.PlainAuth("", defaultConfig.Username, defaultConfig.Password, defaultConfig.Host)

	header := make(map[string]string)
	header["From"] = defaultConfig.From
	header["To"] = to
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	// 配置TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         defaultConfig.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, defaultConfig.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(defaultConfig.From); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
