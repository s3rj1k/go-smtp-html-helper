package sendmail

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

// Config - SMTP config
type Config struct {
	Headers struct {
		From    string
		To      string
		ReplyTo string
		Subject string
		Cc      []string
		Bcc     []string
		IsText  bool
	}
	Body struct {
		Message string
	}
	SMTP struct {
		Server   string
		Port     int
		Email    string
		Password string
	}
}

// Send - send HTML email via SMTP
func (c *Config) Send() error {

	header := make(map[string]string)
	// prepare headers map
	header["From"] = c.Headers.From
	header["To"] = c.Headers.To
	header["Subject"] = c.Headers.Subject
	if len(c.Headers.ReplyTo) > 0 {
		header["Reply-To"] = c.Headers.ReplyTo
	}
	header["MIME-Version"] = "1.0"
	if c.Headers.IsText {
		header["Content-Type"] = "text/plain; charset=\"utf-8\""
	} else {
		header["Content-Type"] = "text/html; charset=\"utf-8\""
	}
	header["Content-Transfer-Encoding"] = "base64"
	if len(c.Headers.Cc) > 0 {
		header["cc"] = strings.Join(c.Headers.Cc, ", ")
	}
	if len(c.Headers.Bcc) > 0 {
		header["bcc"] = strings.Join(c.Headers.Bcc, ", ")
	}

	message := ""
	// assemble headers
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// add base64 encoded message body
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(c.Body.Message))

	return smtp.SendMail(
		net.JoinHostPort(c.SMTP.Server, strconv.Itoa(c.SMTP.Port)),
		smtp.PlainAuth("", c.SMTP.Email, c.SMTP.Password, c.SMTP.Server),
		c.SMTP.Email,
		[]string{c.Headers.To},
		[]byte(message))

}
