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
	Email string `json:"email"`
}

type InertiaPayload struct {
	Component string `json:"component"`
	Props     struct {
		RSVPs []RSVPItem `json:"rsvps"`
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

	seen := make(map[string]bool)
	for _, item := range payload.Props.RSVPs {
		email := strings.ToLower(strings.TrimSpace(item.Email))
		if email != "" && !seen[email] {
			seen[email] = true
			emails = append(emails, email)
		}
	}

	return emails
}