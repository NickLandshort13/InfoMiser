package services

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func FetchHackerTargetSubdomains(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	lines := strings.Split(string(body), "\n")

	var subs []string
	for _, line := range lines {
		parts := strings.SplitN(line, ",", 2)
		if len(parts) > 0 && strings.Contains(parts[0], ".") {
			subs = append(subs, parts[0])
		}
	}

	return subs, nil
}
