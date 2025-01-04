package tests

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"testing"
)

func TestSendEmailWithAttachment(t *testing.T) {
	// Define the SMTP server and email addresses
	smtpHost := "localhost"
	smtpPort := "1025"
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
