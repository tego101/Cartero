package helpers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

// Render Helper for templ components
func Render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	w.Header().Set("Content-Type", "text/html")
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render component", http.StatusInternalServerError)
	}
}

func GenerateUniqueKey(id string) string {
	return strings.ToLower(strings.Replace(id, " ", "_", -1)) + "-" + strconv.Itoa(rand.Intn(100000))
}

func TimeParseFormat(value string) (string, error) {
	layout := "2006-01-02 15:04:05.999999999-07:00"
	t, err := time.Parse(layout, value)
	if err != nil {
		return "", err
	}
	return t.Format("Monday, January 2, 2006 15:04 MST"), nil
}

func TimeAgo(timestamp string) (string, error) {
	// Parse the given timestamp using the appropriate layout
	layout := "2006-01-02 15:04:05.999999999-07:00" // Adjusted for your timestamp format
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		return "", err
	}

	// Calculate the difference from now
	duration := time.Since(t)

	// Get the absolute value of the duration
	seconds := int(duration.Seconds())
	if seconds < 0 {
		return "In the future", nil
	}

	switch {
	case seconds < 60:
		return fmt.Sprintf("%d seconds ago", seconds), nil
	case seconds < 3600:
		minutes := seconds / 60
		return fmt.Sprintf("%d minutes ago", minutes), nil
	case seconds < 86400:
		hours := seconds / 3600
		return fmt.Sprintf("%d hours ago", hours), nil
	default:
		days := seconds / 86400
		return fmt.Sprintf("%d days ago", days), nil
	}
}
