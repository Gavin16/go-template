package main

import (
	"net/smtp"
	"testing"
)

type MailSender interface {
	send(subject, from string, to []string, content string, mailServer string, addr smtp.Auth) error
}

const DISCLAIM = `-------------------------------------------------------
免责声明: 此电子邮件和任何附件可能包含特权和机密信息,仅供特定的收件人使用.如果你收到次电子邮件，
请通知发件人立即删除此电子邮件。
`

func attachDisclaim(content string) string {
	return content + "\n\n" + DISCLAIM
}

// 邮件发送行为抽象为接口，方便在单元测试时通过伪造发送者来进行测试
// 否则就需要依赖联网的三方发送模块，导致代码可测试性下降
func SendMailWithDisclaim(sender MailSender, subject, from string,
	to []string, content string, mailServer string, addr smtp.Auth) error {
	return sender.send(subject, from, to, attachDisclaim(content), mailServer, addr)
}

type FakeEmailSender struct {
	subject string
	from    string
	to      []string
	content string
}

func (fes *FakeEmailSender) send(subject, from string,
	to []string, content string, mailServer string, addr smtp.Auth) error {
	fes.subject = subject
	fes.from = from
	fes.to = to
	fes.content = content
	return nil
}

func TestSendMailWithDisclaim(t *testing.T) {
	sender := &FakeEmailSender{}
	err := SendMailWithDisclaim(sender, "gopher mail test", "mingorun@163.com",
		[]string{"goper@163.com"}, "hello,gopher", "smtp.163.com:25",
		smtp.PlainAuth("", "YOUR_EMIAL_ACCOUNT", "YOUR_EMIAL_ACCOUNT", "smtp.163.com"))

	if err != nil {
		t.Errorf("SendMailWithDisclaim err: %v", err)
		return
	}

	want := "hello,gopher" + "\n\n" + DISCLAIM
	if sender.content != want {
		t.Fatalf("want %s, got %s", want, sender.content)
	}
}
