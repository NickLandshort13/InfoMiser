package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func isDomain(host string) bool {
	if host == "" {
		return false
	}

	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(host)
}

func fetchCrtShSubdomainsWithTimeout(domain string, timeout time.Duration) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", "https://crt.sh/?q="+url.QueryEscape("%."+domain)+"&output=json", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw []map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	subs := make(map[string]bool)

	for _, cert := range raw {
		nameValue, ok := cert["common_name"].(string)
		if ok && strings.Contains(nameValue, ".") && !strings.HasPrefix(nameValue, "*") {
			subs[nameValue] = true
		}

		altNames, ok := cert["name_value"].(string)
		if ok {
			for _, name := range strings.Split(altNames, "\n") {
				if strings.Contains(name, ".") && !strings.HasPrefix(name, "*") {
					subs[name] = true
				}
			}
		}
	}

	result := make([]string, 0, len(subs))
	for sub := range subs {
		result = append(result, sub)
	}

	return result, nil
}

func (h *Handlers) LookupDomain(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	host := cleanInput(q)

	if host == "" {
		w.Write([]byte(""))
		return
	}

	if !isDomain(host) {
		w.Write([]byte(""))
		return
	}

	subdomains, err := fetchCrtShSubdomainsWithTimeout(host, 3*time.Second)
	if err != nil || len(subdomains) == 0 {
		w.Write([]byte(""))
		return
	}

	data := struct {
		Domain     string
		Subdomains []string
	}{
		Domain:     host,
		Subdomains: subdomains,
	}

	h.templates.ExecuteTemplate(w, "results-domains.html", data)
}
