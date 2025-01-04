package handlers

/*
*	Email Handler.
*	-----------------
*	@author https://x.com/tegodotdev <- -> https://github.com/tego101
*	@license MIT
 */

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	helpers "github.com/tego101/cartero-smtp-catch/lib"
	"github.com/tego101/cartero-smtp-catch/types"
	"github.com/tego101/cartero-smtp-catch/views"
	"github.com/tego101/cartero-smtp-catch/views/components"
)

/*
* func PaginateEmails
* @param w http.ResponseWriter
* @param r *http.Request
* @param db *sql.DB
* @param limit int
* @param offset int
* @param isSearched bool
* @param searchQuery string
*
* This function paginates emails in the database based on limit and offset.
* If isSearched is true, it will paginate searched emails based on searchQuery.
 */
func PaginateEmails(w http.ResponseWriter, r *http.Request, db *sql.DB, limit int, offset int, isSearched bool, searchQuery string) {

}

func HandleAllEmails(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Ensure the handler supports only GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Query the database
	rows, err := db.Query("SELECT id, \"from\", \"to\", subject, body, raw, \"timestamp\" FROM emails ORDER BY timestamp DESC")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query emails: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Initialize slice for all emails
	var allEmails []types.EmailProps

	// Iterate through the rows
	for rows.Next() {
		var email types.EmailProps
		if scanErr := rows.Scan(&email.ID, &email.From, &email.To, &email.Subject, &email.Body, &email.Raw, &email.Timestamp); scanErr != nil {
			http.Error(w, fmt.Sprintf("Failed to scan email row: %v", scanErr), http.StatusInternalServerError)
			return
		}
		allEmails = append(allEmails, email)
	}

	// Check for errors during iteration
	if rows.Err() != nil {
		http.Error(w, fmt.Sprintf("Error during rows iteration: %v", rows.Err()), http.StatusInternalServerError)
		return
	}

	helpers.Render(w, r, views.Inbox(allEmails))

}

/*
* func HandleAllEmailsHTMX
* @param w http.ResponseWriter
* @param r *http.Request
* @param db *sql.DB
*
* This function handles all emails and renders them using HTMX.
 */
func HandleAllEmailsHTMX(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	pageNumber, err := strconv.Atoi(page)
	limitPerPage, err := strconv.Atoi(limit)

	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	// If limit is not empty return emails based on limit.
	// If limit is empty return emails based on default limit which is 10.
	if limit != "" {
		rows, err := db.Query("SELECT id, \"from\", \"to\", subject, body, raw, \"timestamp\" FROM emails ORDER BY timestamp DESC LIMIT ?", limitPerPage)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to query emails: %v", err), http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var allEmails []types.EmailProps

		for rows.Next() {
			var email types.EmailProps
			if scanErr := rows.Scan(&email.ID, &email.From, &email.To, &email.Subject, &email.Body, &email.Raw, &email.Timestamp); scanErr != nil {
				http.Error(w, fmt.Sprintf("Failed to scan email row: %v", scanErr), http.StatusInternalServerError)
				return
			}
			allEmails = append(allEmails, email)
		}

		if rows.Err() != nil {
			http.Error(w, fmt.Sprintf("Error during rows iteration: %v", rows.Err()), http.StatusInternalServerError)
			return
		}

		for _, email := range allEmails {
			helpers.Render(w, r, components.EmailRow(email))
		}

	} else {

		rows, err := db.Query("SELECT id, \"from\", \"to\", subject, body, raw, \"timestamp\" FROM emails ORDER BY timestamp DESC")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to query emails: %v", err), http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var allEmails []types.EmailProps

		for rows.Next() {
			var email types.EmailProps
			if scanErr := rows.Scan(&email.ID, &email.From, &email.To, &email.Subject, &email.Body, &email.Raw, &email.Timestamp); scanErr != nil {
				http.Error(w, fmt.Sprintf("Failed to scan email row: %v", scanErr), http.StatusInternalServerError)
				return
			}
			allEmails = append(allEmails, email)
		}

		if rows.Err() != nil {
			http.Error(w, fmt.Sprintf("Error during rows iteration: %v", rows.Err()), http.StatusInternalServerError)
			return
		}

		// pagination
		var emailsPerPage = 10
		var start = (pageNumber - 1) * emailsPerPage
		var end = start + emailsPerPage

		if start > len(allEmails) {
			start = len(allEmails)
		}

		if end > len(allEmails) {
			end = len(allEmails)
		}

		if limitPerPage > 0 {
			emailsPerPage = limitPerPage
		}

		for i := start; i < end; i++ {
			helpers.Render(w, r, components.EmailRow(allEmails[i]))
		}

	}
}

/*
* func HandleSearchEmailsHTMX
* @param w http.ResponseWriter
* @param r *http.Request
* @param db *sql.DB
*
* This function handles search emails and renders them using HTMX.
 */
func HandleSearchEmailsHTMX(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// search post params
	search := r.FormValue("q")
	searchQuery := "%" + search + "%"

	// If the search is empty by default it will respond with all emails.
	if search == "" {
		HandleAllEmailsHTMX(w, r, db)
		return
	}

	rows, err := db.Query("SELECT id, \"from\", \"to\", subject, body, raw, \"timestamp\" FROM emails WHERE \"from\" LIKE ? OR \"to\" LIKE ? OR subject LIKE ? OR body LIKE ? ORDER BY timestamp DESC", searchQuery, searchQuery, searchQuery, searchQuery)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query emails: %v", err), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var allEmails []types.EmailProps

	// if there are no emails found from the search criteria it will return an empty array.
	if rows.Next() == false {
		helpers.Render(w, r, components.EmptyEmailRow(
			components.EmptyEmailRowProps{
				Message: "Nothing found for " + search,
			},
		))
		return
	}

	for rows.Next() {
		var email types.EmailProps
		if scanErr := rows.Scan(&email.ID, &email.From, &email.To, &email.Subject, &email.Body, &email.Raw, &email.Timestamp); scanErr != nil {
			http.Error(w, fmt.Sprintf("Failed to scan email row: %v", scanErr), http.StatusInternalServerError)
			return
		}
		allEmails = append(allEmails, email)
	}

	if rows.Err() != nil {
		http.Error(w, fmt.Sprintf("Error during rows iteration: %v", rows.Err()), http.StatusInternalServerError)
		return
	}

	for _, email := range allEmails {
		helpers.Render(w, r, components.EmailRow(email))
	}

}

/*
* func HandleViewEmailHTMX
* @param w http.ResponseWriter
* @param r *http.Request
* @param db *sql.DB
*
 */
func HandleViewEmailHTMX(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.URL.Path[len("/mail/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	var email types.EmailProps

	err = db.QueryRow("SELECT id, \"from\", \"to\", subject, body, raw, \"timestamp\" FROM emails WHERE id = ?", id).Scan(&email.ID, &email.From, &email.To, &email.Subject, &email.Body, &email.Raw, &email.Timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Println(err)
			http.Error(w, "Failed to query email", http.StatusInternalServerError)
		}
		return
	}

	helpers.Render(w, r, components.Email(email, w))
}

/*
* func HandleDeleteEmailHTMX
* @param w http.ResponseWriter
* @param r *http.Request
* @param db *sql.DB
*
* This function deletes an email.
 */
func HandleDeleteEmailHTMX(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.URL.Path[len("/mail/delete/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	_, err = db.Exec("DELETE FROM emails WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Println(err)
			http.Error(w, "Failed to delete email", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("HX-Redirect", "/inbox")
	w.WriteHeader(http.StatusOK)
}

func HandleAllEmailsJSON(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT id, \"from\", \"to\", subject, body, raw, \"timestamp\" FROM emails ORDER BY timestamp DESC")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query emails: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var allEmails []types.EmailProps

	for rows.Next() {
		var email types.EmailProps
		if scanErr := rows.Scan(&email.ID, &email.From, &email.To, &email.Subject, &email.Body, &email.Raw, &email.Timestamp); scanErr != nil {
			http.Error(w, fmt.Sprintf("Failed to scan email row: %v", scanErr), http.StatusInternalServerError)
			return
		}
		allEmails = append(allEmails, email)
	}

	if rows.Err() != nil {
		http.Error(w, fmt.Sprintf("Error during rows iteration: %v", rows.Err()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(allEmails); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode emails as JSON: %v", err), http.StatusInternalServerError)
	}
}

func HandleViewEmailJSON(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.URL.Path[len("/inbox.json/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	var email types.EmailProps

	err = db.QueryRow("SELECT id, \"from\", \"to\", subject, body, raw, \"timestamp\" FROM emails WHERE id = ?", id).Scan(&email.ID, &email.From, &email.To, &email.Subject, &email.Body, &email.Raw, &email.Timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Println(err)
			http.Error(w, "Failed to query email", http.StatusInternalServerError)
		}
		return
	}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode emails as JSON and write to response
	if err := json.NewEncoder(w).Encode(email); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode email as JSON: %v", err), http.StatusInternalServerError)
	}
}

func HandleDeleteEmail() string {
	return ""
}

func HandleViewEmail(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.URL.Path[len("/inbox/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	var email types.EmailProps

	err = db.QueryRow("SELECT id, \"from\", \"to\", subject, body, raw, \"timestamp\" FROM emails WHERE id = ?", id).Scan(&email.ID, &email.From, &email.To, &email.Subject, &email.Body, &email.Raw, &email.Timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Println(err)
			http.Error(w, "Failed to query email", http.StatusInternalServerError)
		}
		return
	}

	helpers.Render(w, r, views.EmailView(email, w))
}
