package types

type EmailProps struct {
	ID          int
	Timestamp   string
	From        string
	To          string
	Subject     string
	Body        string
	Raw         string
	Attachments []Attachment
}

type Attachment struct {
	ID       int
	EmailID  int
	Filename string
	MIMEType string
	Content  []byte
}
