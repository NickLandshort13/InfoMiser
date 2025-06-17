package handlers

import (
	"net/url"
	"regexp"
	"strings"
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
