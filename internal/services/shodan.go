package services

import (
	"encoding/json"
	"infomiser/internal/models"
	"io"
	"net/http"
)

func QueryShodan(ip string) (models.Shodan, error) {
	var result models.Shodan
	if ip == "" {
		return result, nil
	}

	resp, err := http.Get("https://internetdb.shodan.io/" + ip)
	if err != nil || resp.StatusCode != 200 {
		return result, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	if ports, ok := data["ports"].([]interface{}); ok {
		for _, p := range ports {
			if port, ok := p.(float64); ok {
				result.Ports = append(result.Ports, int(port))
			}
		}
	}

	result.IP = ip

	return result, nil
}
