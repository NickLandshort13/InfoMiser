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
		result.Error = "no domain"
		return result, nil
	}

	url := fmt.Sprintf("https://haveibeenpwned.com/api/v3/breaches?domain=%s", url.QueryEscape(domain))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "InfoMiser")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		result.Domain = domain
		result.Breaches = nil
		return result, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var breaches []map[string]interface{}
	json.Unmarshal(body, &breaches)

	pwnedBreaches := make([]models.PwnedBreach, 0)
	for _, b := range breaches {
		name, _ := b["Name"].(string)
		title, _ := b["Title"].(string)
		breachDate, _ := b["BreachDate"].(string)
		desc, _ := b["Description"].(string)
		count, _ := b["PwnCount"].(float64)
		dataClasses, _ := b["DataClasses"].([]interface{})
		logo, _ := b["LogoPath"].(string)
		isVerified, _ := b["IsVerified"].(bool)

		classes := make([]string, 0)
		for _, dc := range dataClasses {
			if s, ok := dc.(string); ok {
				classes = append(classes, s)
			}
		}

		pwnedBreaches = append(pwnedBreaches, models.PwnedBreach{
			Name:        name,
			Title:       title,
			BreachDate:  breachDate,
			Description: desc,
			PwnCount:    int(count),
			DataClasses: classes,
			LogoPath:    logo,
			IsVerified:  isVerified,
		})
	}

	result = models.Pwned{
		Domain:   domain,
		Breaches: pwnedBreaches,
		Error:    "",
	}

	return result, nil
}
