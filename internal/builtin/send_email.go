package builtin

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

func SendEmailTask(params map[string]string) error {
	required := []string{
		"smtp_host", "smtp_port", "smtp_user", "smtp_password",
		"from", "to", "subject", "body_html",
	}
	for _, key := range required {
		if _, ok := params[key]; !ok {
			return fmt.Errorf("missing required parameter: %s", key)
		}
	}

	smtpHost := params["smtp_host"]
	smtpPort := params["smtp_port"]
	smtpUser := params["smtp_user"]
	smtpPass := params["smtp_password"]
	from := params["from"]
	to := strings.Split(params["to"], ",")
	subject := params["subject"]
	bodyHTML := params["body_html"]
	bodyText := params["body_text"]

	message := buildMimeEmail(from, to, subject, bodyHTML, bodyText)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return fmt.Errorf("tls dial failed: %w", err)
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("smtp client failed: %w", err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth failed: %w", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("smtp MAIL FROM failed: %w", err)
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("smtp RCPT TO failed for %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA failed: %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("smtp write failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("smtp close failed: %w", err)
	}

	return nil
}

func buildMimeEmail(from string, to []string, subject, html, text string) string {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	if text != "" {
		// Add multipart with text fallback
		boundary := "mixed-boundary"
		return fmt.Sprintf(`From: %s
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="%s"

--%s
Content-Type: text/plain; charset="UTF-8"

%s

--%s
Content-Type: text/html; charset="UTF-8"

%s

--%s--`, from, strings.Join(to, ","), subject, boundary, boundary, text, boundary, html, boundary)
	}

	// HTML-only
	return fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n%s\r\n%s",
		from, strings.Join(to, ","), subject, mime, html)
}
