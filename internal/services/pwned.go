package services

import (
	"encoding/json"
	"fmt"
	"infomiser/internal/models"
	"io"
	"net/http"
	"net/url"
)

func CheckBreaches(domain string) (models.Pwned, error) {
	var result models.Pwned
	if domain == "" {
		return result, nil
	}

	reqUrl := fmt.Sprintf("https://haveibeenpwned.com/api/v3/breachedaccount/%s", url.QueryEscape(domain))

	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Set("User-Agent", "InfoMiser")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return result, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var breaches []map[string]interface{}
	json.Unmarshal(body, &breaches)

	names := make([]string, 0)
	for _, b := range breaches {
		if name, ok := b["Name"].(string); ok {
			names = append(names, name)
		}
	}

	result.Domain = domain
	result.Breaches = names

	return result, nil
}
