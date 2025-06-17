package handlers

import (
	"net"
	"net/http"
	"time"

	"infomiser/internal/models"
	"infomiser/internal/services"
)

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Domain     string
		IP         string
		IPWhois    models.IPWhois
		Subdomains models.Subdomains
		Shodan     models.Shodan
		Pwned      models.Pwned
	}{
		Domain: "",
		IP:     "",
		IPWhois: models.IPWhois{
			IP:      "",
			Country: "",
			Region:  "",
			City:    "",
			Lat:     0,
			Lon:     0,
			ISP:     "",
			Domain:  "",
		},
		Subdomains: models.Subdomains{
			Domain:     "",
			Subdomains: nil,
		},
		Shodan: models.Shodan{
			IP:    "",
			Ports: nil,
		},
		Pwned: models.Pwned{
			Domain:   "",
			Breaches: nil,
		},
	}

	h.templates.ExecuteTemplate(w, "index.html", data)
}

func (h *Handlers) Lookup(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	host := cleanInput(q)

	if host == "" || !isValidQuery(host) {
		h.templates.ExecuteTemplate(w, "results-multi.html", nil)
		return
	}

	var ip string
	if isDomain(host) {
		tmp, err := services.GetIPFromDomain(host)
		if err == nil && tmp != "" {
			ip = tmp
		}
	} else if net.ParseIP(host) != nil {
		ip = host
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
	}()

	ipWhois, _ := services.GetIPWhois(ip)
	subdomains, _ := services.FetchHackerTargetSubdomains(host)
	shodan, _ := services.QueryShodan(ip)
	pwned, _ := services.CheckBreaches(host)
	data := struct {
		Domain     string
		IP         string
		IPWhois    models.IPWhois
		Subdomains models.Subdomains
		Shodan     models.Shodan
		Pwned      models.Pwned
	}{
		Domain:     host,
		IP:         ip,
		IPWhois:    ipWhois,
		Subdomains: models.Subdomains{Domain: host, Subdomains: subdomains},
		Shodan:     shodan,
		Pwned:      pwned,
	}

	h.templates.ExecuteTemplate(w, "results-multi", data)
}
