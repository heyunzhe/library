package mode

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"
)

// EmailConfig 邮件配置，从环境变量读取
type EmailConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	FromName string
}

// GetEmailConfig 获取邮件配置
func GetEmailConfig() EmailConfig {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "smtp.qq.com"
	}
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}
	return EmailConfig{
		Host:     host,
		Port:     port,
		User:     os.Getenv("SMTP_USER"),
		Pass:     os.Getenv("SMTP_PASS"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
	}
}

// loginAuth 实现 smtp.Auth 接口，用于 QQ SMTP 的 LOGIN 认证
type loginAuth struct {
	user, pass string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.user), nil
		case "Password:":
			return []byte(a.pass), nil
		default:
			// 某些服务器会发送 base64 编码的提示
			return []byte(a.user), nil
		}
	}
	return nil, nil
}

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	cfg := GetEmailConfig()
	if cfg.User == "" || cfg.Pass == "" {
		return fmt.Errorf("SMTP_USER 或 SMTP_PASS 未设置")
	}

	fromName := cfg.FromName
	if fromName == "" {
		fromName = "图书馆管理系统"
	}

	auth := &loginAuth{cfg.User, cfg.Pass}

	contentType := "Content-Type: text/plain; charset=UTF-8"
	msg := []byte(fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\n%s\r\n\r\n%s",
		fromName, cfg.User, to, subject, contentType, body))

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return smtp.SendMail(addr, auth, cfg.User, []string{to}, msg)
}

// GenerateVerificationCode 生成6位数字验证码
func GenerateVerificationCode() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := rng.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

// SendVerificationCode 发送验证码邮件
func SendVerificationCode(email, code string) error {
	subject := "图书馆注册验证码"
	body := fmt.Sprintf(`您好！

感谢您注册图书馆管理系统。

您的验证码是：%s

该验证码有效期为10分钟，请尽快完成注册。

如果这不是您本人的操作，请忽略此邮件。`, code)
	return SendEmail(email, subject, body)
}
