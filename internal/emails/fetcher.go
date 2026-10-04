package emails

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/reztheperson/boba-bash-shipifyer/internal/keyring"
)

type RSVPItem struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	HeardAbout string `json:"heard_about"`
}

type SubmissionItem struct {
	ID     int `json:"id"`
	RsvpID int `json:"rsvp_id"`
}

type InertiaPayload struct {
	Component string `json:"component"`
	Props     struct {
		RSVPs       []RSVPItem       `json:"rsvps"`
		Submissions []SubmissionItem `json:"submissions"`
	} `json:"props"`
}

func Fetch() []string {
	var emails []string

	targetURL := fmt.Sprintf("https://bash.hackclub.com/organize/%v", keyring.Secrets.PanelOrgNumber)
	sessionID := keyring.Secrets.PanelOrgCookie

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		log.Printf("[fetcher] Failed creating request: %v", err)
		return emails
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Cookie", fmt.Sprintf("_session_id=%s", sessionID))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[fetcher] Network request failed: %v", err)
		return emails
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther {
		log.Printf("[fetcher] Session expired or unauthorized (Redirected to: %s)", resp.Header.Get("Location"))
		return emails
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[fetcher] Failed reading response: %v", err)
		return emails
	}

	re := regexp.MustCompile(`data-page="([^"]+)"`)
	matches := re.FindSubmatch(bodyBytes)
	if len(matches) < 2 {
		log.Printf("[fetcher] data-page attribute missing (Status: %d)", resp.StatusCode)
		return emails
	}

	unescaped := html.UnescapeString(string(matches[1]))

	var payload InertiaPayload
	if err := json.Unmarshal([]byte(unescaped), &payload); err != nil {
		log.Printf("[fetcher] Failed parsing Inertia JSON: %v", err)
		return emails
	}

	// Map rsvp_id -> normalized email
	rsvpEmailMap := make(map[int]string)
	for _, item := range payload.Props.RSVPs {
		normalized := strings.ToLower(strings.TrimSpace(item.Email))
		if normalized != "" {
			rsvpEmailMap[item.ID] = normalized
		}
	}

	// Identify emails linked to existing submissions
	submittedEmails := make(map[string]bool)
	for _, sub := range payload.Props.Submissions {
		if email, found := rsvpEmailMap[sub.RsvpID]; found {
			submittedEmails[email] = true
		}
	}

	// Collect unique emails that do not have a submission
	seen := make(map[string]bool)
	for _, item := range payload.Props.RSVPs {
		email := strings.ToLower(strings.TrimSpace(item.Email))
		if email == "" || seen[email] || submittedEmails[email] {
			continue
		}

		seen[email] = true
		emails = append(emails, email)
	}

	log.Printf("[fetcher] Done. RSVPs: %d, Submissions: %d -> Returning %d unsubmitted unique emails.",
		len(payload.Props.RSVPs), len(payload.Props.Submissions), len(emails))

	return emails
}