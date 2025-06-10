package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"infomiser/internal/models"
)

var client = &http.Client{Timeout: 10 * time.Second}

func (h *Handlers) Lookup(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	host := cleanInput(q)

	if host == "" || !isValidQuery(host) {
		w.Write([]byte(""))
		return
	}

	var ip string
	var err error

	if isDomain(host) {
		ip, err = getIPFromDomain(host)
		if err != nil {
			ip = ""
		}
	} else if net.ParseIP(host) != nil {
		ip = host
	} else {
		w.Write([]byte(""))
		return
	}

	var ipInfo models.IPInfo
	if ip != "" {
		resp, err := client.Get("https://ipwho.is/" + url.PathEscape(ip))
		if err == nil && resp.StatusCode == 200 {
			body, _ := io.ReadAll(resp.Body)
			if len(body) > 0 && body[0] == '{' {
				json.Unmarshal(body, &ipInfo)
			}
		}
	}

	var subdomains []string
	if isDomain(host) {
		subs, err := fetchHackerTargetSubdomains(host)
		if err == nil && len(subs) > 0 {
			subdomains = subs
		}
	}

	data := struct {
		HasIP         bool
		IP            string
		IPInfo        models.IPInfo
		HasSubdomains bool
		Domain        string
		Subdomains    []string
	}{
		HasIP:         ip != "",
		IP:            ip,
		IPInfo:        ipInfo,
		HasSubdomains: len(subdomains) > 0,
		Domain:        host,
		Subdomains:    subdomains,
	}

	h.templates.ExecuteTemplate(w, "results-multi", data)
}
func fetchHackerTargetSubdomains(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	lines := strings.Split(string(body), "\n")

	var subs []string
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) > 0 && strings.Contains(parts[0], ".") {
			subs = append(subs, parts[0])
		}
	}

	return removeDuplicates(subs), nil
}

func removeDuplicates(list []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, item := range list {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
func cleanInput(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	if u, err := url.Parse(input); err == nil && u.Host != "" {
		return u.Host
	}
	if strings.Contains(input, "/") {
		parts := strings.SplitN(input, "/", 2)
		input = parts[0]
	}
	return input
}

func isDomain(host string) bool {
	if host == "" {
		return false
	}
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(host)
}

func isValidQuery(s string) bool {
	if s == "" {
		return false
	}
	if len(s) > 100 {
		return false
	}
	return true
}

func getIPFromDomain(domain string) (string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}
	for _, ip := range ips {
		if v4 := ip.To4(); v4 != nil {
			return v4.String(), nil
		}
	}
	return "", nil
}
