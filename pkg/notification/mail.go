package notification

import (
	"gopkg.in/gomail.v2"
	"strconv"
)

func SendMail(userName, authCode, host, portStr, mailTo, sendName string, subject, body string) error {
	port, _ := strconv.Atoi(portStr)
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(userName, sendName))
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	m.Embed("D:\\Project\\emolve\\autoclick\\storage\\images\\auto_click.png") // 图片路径
	m.SetBody("text/html", `<img src="cid:auto_click.png" alt="My image" />`)  //设置邮件正文

	d := gomail.NewDialer(host, port, userName, authCode)
	err := d.DialAndSend(m)
	return err
}

func Send163Mail(token, mailTo, subject string) error {
	err := SendMail(mailTo, token, "smtp.163.com", "25", mailTo, mailTo, subject, "11111")
	return err
}
