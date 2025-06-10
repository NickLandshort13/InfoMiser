package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type IPInfo struct {
	IP            string  `json:"ip"`
	Success       bool    `json:"success"`
	Type          string  `json:"type"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	City          string  `json:"city"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	IsEU          bool    `json:"is_eu"`
	Postal        string  `json:"postal"`
	CallingCode   string  `json:"calling_code"`
	Capital       string  `json:"capital"`
	Borders       string  `json:"borders"`

	Flag struct {
		Img          string `json:"img"`
		Emoji        string `json:"emoji"`
		EmojiUnicode string `json:"emoji_unicode"`
	} `json:"flag"`

	Connection struct {
		ASN    int    `json:"asn"`
		Org    string `json:"org"`
		ISP    string `json:"isp"`
		Domain string `json:"domain"`
	} `json:"connection"`

	Timezone struct {
		ID           string `json:"id"`
		Abbreviation string `json:"abbr"`
		IsDst        bool   `json:"is_dst"`
		Offset       int    `json:"offset"`
		UTC          string `json:"utc"`
		CurrentTime  string `json:"current_time"`
	} `json:"timezone"`
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "index.html", map[string]string{
		"Title": "InfoMiser — OSINT Lookup",
	})
}

func (h *Handlers) Lookup(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	q = strings.TrimSpace(q)
	if q == "" {
		h.templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	resp, err := http.Get("https://ipwho.is/" + url.PathEscape(q))
	if err != nil || resp.StatusCode != 200 {
		h.templates.ExecuteTemplate(w, "error.html", nil)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	if len(body) == 0 || body[0] != '{' {
		h.templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	var info IPInfo
	if err := json.Unmarshal(body, &info); err != nil {
		h.templates.ExecuteTemplate(w, "error.html", nil)
		return
	}
	fmt.Println(url.PathEscape(q))
	h.templates.ExecuteTemplate(w, "results.html", info)
}
