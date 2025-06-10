package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func isDomain(host string) bool {
	if host == "" {
		return false
	}

	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(host)
}

func fetchCrtShSubdomains(domain string) ([]string, error) {
	url := "https://crt.sh/?q=" + url.QueryEscape("%."+domain) + "&output=json"

	resp, err := http.Get(url)
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
		if ok {
			if strings.Contains(nameValue, ".") && !strings.HasPrefix(nameValue, "*") {
				subs[nameValue] = true
			}
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

	subdomains, err := fetchCrtShSubdomains(host)
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
