package tests

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestSendEmailWithAttachment(t *testing.T) {
	// Define the SMTP server and email addresses
	envError := godotenv.Load("../.env")

	if envError != nil {
		log.Fatal("Error loading .env file")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if smtpHost == "" || smtpPort == "" {
		t.Fatalf("SMTP_HOST or SMTP_PORT environment variable is not set")
	}
	from := "sender@example.com"
	to := "receiver@example.com"
	subject := "Hello World"
	body := "Hello, this is the body."

	// Read the file to attach
	fileName := "sample_text.txt"
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", fileName, err)
	}

	// Encode the file content to Base64
	encodedFileContent := base64.StdEncoding.EncodeToString(fileContent)

	// Construct the MIME message
	boundary := "boundary123"
	var message bytes.Buffer
	message.WriteString(fmt.Sprintf("From: %s\r\n", from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	message.WriteString("\r\n")
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	message.WriteString("\r\n")
	message.WriteString(body + "\r\n")
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString(fmt.Sprintf("Content-Type: text/plain; name=\"%s\"\r\n", fileName))
	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName))
	message.WriteString("\r\n")
	message.WriteString(encodedFileContent + "\r\n")
	message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	// Send the email
	//auth := smtp.PlainAuth("", "", "", smtpHost)
	err = smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), nil, from, []string{to}, message.Bytes())
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	t.Log("Email sent successfully")
}

func TestSendEmailBasic(t *testing.T) {
	envError := godotenv.Load("../.env")

	if envError != nil {
		log.Fatal("Error loading .env file")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if smtpHost == "" || smtpPort == "" {
		t.Fatalf("SMTP_HOST or SMTP_PORT environment variable is not set")
	}
	from := "basic_sender@example.com"
	to := "basic_receiver@example.com"
	subject := "Basic Test Email"
	body := "This is a basic test email without attachments."

	// Construct the email message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), nil, from, []string{to}, []byte(message))
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	t.Log("Basic email sent successfully")
}

func TestSendEmailWithLargeAttachment(t *testing.T) {
	envError := godotenv.Load("../.env")

	if envError != nil {
		log.Fatal("Error loading .env file")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if smtpHost == "" || smtpPort == "" {
		t.Fatalf("SMTP_HOST or SMTP_PORT environment variable is not set")
	}
	from := "attachment_sender@example.com"
	to := "attachment_receiver@example.com"
	subject := "Large Attachment Test Email"
	body := "This email contains a large attachment."

	// Simulate a large file by creating an in-memory buffer
	fileName := "large_attachment.txt"
	fileContent := bytes.Repeat([]byte("This is a line of content.\n"), 10000)
	encodedFileContent := base64.StdEncoding.EncodeToString(fileContent)

	// Construct the MIME message
	boundary := "boundary123"
	var message bytes.Buffer
	message.WriteString(fmt.Sprintf("From: %s\r\n", from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	message.WriteString("\r\n")
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	message.WriteString("\r\n")
	message.WriteString(body + "\r\n")
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString(fmt.Sprintf("Content-Type: text/plain; name=\"%s\"\r\n", fileName))
	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName))
	message.WriteString("\r\n")
	message.WriteString(encodedFileContent + "\r\n")
	message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), nil, from, []string{to}, message.Bytes())
	if err != nil {
		t.Fatalf("Failed to send email with large attachment: %v", err)
	}

	t.Log("Email with large attachment sent successfully")
}

func TestSendEmailToMultipleRecipients(t *testing.T) {
	envError := godotenv.Load("../.env")

	if envError != nil {
		log.Fatal("Error loading .env file")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if smtpHost == "" || smtpPort == "" {
		t.Fatalf("SMTP_HOST or SMTP_PORT environment variable is not set")
	}
	from := "multi_sender@example.com"
	to := []string{"receiver1@example.com", "receiver2@example.com"}
	subject := "Multiple Recipients Test Email"
	body := "This email is sent to multiple recipients."

	// Construct the email message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to[0]+", "+to[1], subject, body)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), nil, from, to, []byte(message))
	if err != nil {
		t.Fatalf("Failed to send email to multiple recipients: %v", err)
	}

	t.Log("Email to multiple recipients sent successfully")
}

func TestSendEmailHTML(t *testing.T) {
	envPath, err := filepath.Abs("../.env")
	if err != nil {
		t.Fatalf("Failed to determine absolute path for .env file: %v", err)
	}
	envError := godotenv.Load(envPath)

	if envError != nil {
		log.Fatal("Error loading .env file")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if smtpHost == "" || smtpPort == "" {
		t.Fatalf("SMTP_HOST or SMTP_PORT environment variable is not set")
	}
	from := "html_sender@example.com"
	to := "html_receiver@example.com"
	subject := "HTML Test Email"
	body := `<html><body><h1>Hello!</h1><p>This is an HTML email.</p></body></html>`

	// Construct the email message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html\r\n\r\n%s", from, to, subject, body)

	err = smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), nil, from, []string{to}, []byte(message))
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	t.Log("HTML email sent successfully")
}

func TestSendEmailWithAlternatePort(t *testing.T) {
	envError := godotenv.Load("../.env")

	if envError != nil {
		log.Fatal("Error loading .env file")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if smtpHost == "" || smtpPort == "" {
		t.Fatalf("SMTP_HOST or SMTP_PORT environment variable is not set")
	}
	from := "alt_port_sender@example.com"
	to := "alt_port_receiver@example.com"
	subject := "Alternate Port Test Email"
	body := "This email is sent using an alternate SMTP port."

	// Construct the email message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), nil, from, []string{to}, []byte(message))
	if err != nil {
		t.Fatalf("Failed to send email via alternate port: %v", err)
	}

	t.Log("Email sent via alternate port successfully")
}
