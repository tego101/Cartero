package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tego101/cartero-smtp-catch/handlers"
	"github.com/tego101/cartero-smtp-catch/types"
)

var db *sql.DB

func main() {
	envError := godotenv.Load()

	if envError != nil {
		log.Fatal("Error loading .env file")
	}

	// Define ENV variables.
	SMTP_PORT := os.Getenv("SMTP_PORT")
	SMTP_HOST := os.Getenv("SMTP_HOST")
	WEB_PORT := os.Getenv("WEB_PORT")
	WEB_HOST := os.Getenv("WEB_HOST")

	// Set default values if not set.
	if err, ok := os.LookupEnv("SMTP_PORT"); !ok || err == "" {
		SMTP_PORT = "1025"
	}

	if err, ok := os.LookupEnv("SMTP_HOST"); !ok || err == "" {
		log.Println("SMTP_HOST not set, defaulting to localhost")
		SMTP_HOST = "localhost"
	}

	if err, ok := os.LookupEnv("WEB_PORT"); !ok || err == "" {
		log.Println("WEB_PORT not set, defaulting to 10122")
		WEB_PORT = "10122"
	}

	if err, ok := os.LookupEnv("WEB_HOST"); !ok || err == "" {
		log.Println("WEB_HOST not set, defaulting to localhost")
		WEB_HOST = "localhost"
	}

	var err error
	db, err = sql.Open("sqlite3", "./inbox.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Query the count of emails
	var emailCount int
	errEmailCount := db.QueryRow("SELECT COUNT(*) FROM emails").Scan(&emailCount)
	if errEmailCount != nil {
		log.Fatalf("Failed counting emails: %v", err)
		return
	}

	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS emails (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		"from" TEXT,
		"to" TEXT,
		subject TEXT,
		body TEXT,
		raw TEXT,
		attachments TEXT,
		timestamp TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println(`
		        _________
		      =| ~~~  (*)|
		     ==| ~~~~~   |
		    ===|_________|

		 [Cartero] SMTP CatchAll.
		 <repo> github.com/tego101/cartero-smtp-catch
		 <author> x.com/tegodotdev
		 <license> MIT.
		 <disclaimer> This app is opensource and
		 FREE for anyone to use, modify, and redistribute.
		 If you paid for this, please ask for a refund.
		=_=_=_=_=_=_=_=_=_=_=_=_=_=_=_=_=_=_=_=_=_=_=_=_=
	`)

	go startSMTPServer(SMTP_PORT, SMTP_HOST)

	// routes.
	http.HandleFunc("/", redirectToInbox)

	http.HandleFunc("/inbox", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleAllEmails(w, r, db, types.InboxConfig{
			Label: "Default",
			Host:  SMTP_HOST,
			Port:  SMTP_PORT,
		})
	})

	http.HandleFunc("/inbox/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleViewEmail(w, r, db)
	})

	http.HandleFunc("/mail/all", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleAllEmailsHTMX(w, r, db)
	})

	http.HandleFunc("/mail/search", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleSearchEmailsHTMX(w, r, db)
	})

	http.HandleFunc("/mail/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleDeleteEmailHTMX(w, r, db)
	})

	http.HandleFunc("/mail/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleViewEmail(w, r, db)
	})

	http.HandleFunc("/mail.json/all", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleAllEmailsJSON(w, r, db)
	})

	http.HandleFunc("/mail.json/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleViewEmailJSON(w, r, db)
	})

	// End routes.

	// Serve static assets.
	fs := http.FileServer(http.Dir("./assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	webHostPort := fmt.Sprintf("%s:%s", WEB_HOST, WEB_PORT)

	log.Println(fmt.Sprintf("E-Mails Captured: ➡️ %s", strconv.Itoa(emailCount)))
	log.Println(fmt.Sprintf("Web Server: ➡️ https://%s:%s", WEB_HOST, WEB_PORT))

	// AHOY!
	if err := http.ListenAndServe(webHostPort, nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func redirectToInbox(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/inbox", http.StatusFound) // 302 Found
		return
	}
	http.NotFound(w, r)
}

func startSMTPServer(port string, host string) {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("Failed to start SMTP server: %v", err)
	}
	defer listener.Close()
	log.Printf("SMTP server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleSMTPConnection(conn)
	}
}

func handleSMTPConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	fmt.Fprint(writer, "220 Welcome to the Cartero SMTP catch-all server\r\n")
	writer.Flush()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from client: %v", err)
			return
		}

		log.Printf("Client: %s", strings.TrimSpace(line))

		switch {
		case strings.HasPrefix(line, "HELO"):
			fmt.Fprint(writer, "250 Hello\r\n")
		case strings.HasPrefix(line, "EHLO"):
			fmt.Fprint(writer, "250 Hello\r\n")
		case strings.HasPrefix(line, "MAIL FROM:"):
			fmt.Fprint(writer, "250 OK\r\n")
		case strings.HasPrefix(line, "MAIL TO:"):
			fmt.Fprint(writer, "250 OK\r\n")
		case strings.HasPrefix(line, "DATA"):
			fmt.Fprint(writer, "354 End data with <CR><LF>.<CR><LF>\r\n")
			writer.Flush()

			email, err := readEmail(reader)
			if err != nil {
				log.Printf("Error reading email data: %v", err)
				return
			}

			if err := saveEmail(email); err != nil {
				log.Printf("Failed to save email: %v", err)
				fmt.Fprint(writer, "550 Failed to save email\r\n")
				writer.Flush()
				return
			}

			fmt.Fprint(writer, "250 OK: Message received\r\n")
		case strings.HasPrefix(line, "QUIT"):
			fmt.Fprint(writer, "221 Bye\r\n")
			writer.Flush()
			return
		default:
			fmt.Fprint(writer, "500 Unrecognized command\r\n")
		}
		writer.Flush()
	}
}

func readEmail(reader *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		if line == ".\r\n" {
			break
		}
		sb.WriteString(line)
	}
	fmt.Println("Email Incoming ...")
	fmt.Println("Reading Email...")

	return sb.String(), nil
}

func saveEmail(email string) error {
	msg, err := mail.ReadMessage(strings.NewReader(email))
	if err != nil {
		return fmt.Errorf("failed to parse email: %v", err)
	}

	from := msg.Header.Get("From")
	to := msg.Header.Get("To")
	subject := msg.Header.Get("Subject")

	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return fmt.Errorf("failed to read email body: %v", err)
	}

	_, err = db.Exec("INSERT INTO emails (\"from\", \"to\", subject, body, raw, \"timestamp\") VALUES (?, ?, ?, ?, ?, ?)", from, to, subject, string(body), email, time.Now())
	return err
}
