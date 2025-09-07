package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// 使用 163 邮箱发送邮件
func main() {
	// 检查命令行参数
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <sender_email> <receiver_email>")
	}

	// 从命令行参数读取发件人和收件人邮箱
	from := os.Args[1]
	receiverEmail := os.Args[2]

	// 从环境变量读取密码（注意：这里用的是授权码，而不是登录密码）
	password := os.Getenv("SMTP_PASSWORD")
	if password == "" {
		log.Fatal("SMTP_PASSWORD environment variable is not set")
	}

	// 收件人
	to := []string{receiverEmail}

	// SMTP 服务器配置
	smtpHost := "smtp.163.com"
	smtpPort := "25" // 可用 25 或 587 (STARTTLS)

	// 邮件内容
	message := []byte("Subject: 测试邮件\r\n" +
		"From: " + from + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"\r\n" +
		"你好，这是一封 Go 语言发出的测试邮件！\r\n")

	// 认证信息
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 发送邮件
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("邮件发送成功！")
}
