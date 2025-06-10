package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"infomiser/internal/models"
)

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

func getIPFromDomain(domain string) (string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}

	return "", fmt.Errorf("no IPv4 found")
}

func (h *Handlers) Lookup(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query().Get("q")
	host := cleanInput(q)

	if host == "" {
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
		resp, err := http.Get("https://ipwho.is/" + url.PathEscape(ip))
		if err == nil && resp.StatusCode == 200 {
			body, _ := io.ReadAll(resp.Body)
			if len(body) > 0 && body[0] == '{' {
				json.Unmarshal(body, &ipInfo)
			}
		}
	}

	var subdomains []string
	if isDomain(host) {
		subs, err := fetchCrtShSubdomains(host)
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

	h.templates.ExecuteTemplate(w, "results-multi.html", data)
	fmt.Println(data.HasIP)
}
