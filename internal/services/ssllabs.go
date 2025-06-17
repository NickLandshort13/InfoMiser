package services

import (
	"encoding/json"
	"fmt"
	"infomiser/internal/models"
	"io"
	"net/http"
)

func AnalyzeSSL(domain string) (models.SSLScan, error) {
	var result models.SSLScan
	if domain == "" {
		return result, nil
	}

	url := fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s&fromCache=on", domain)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		result.Valid = false
		return result, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	status, ok := data["status"].(string)
	if !ok || status != "READY" {
		result.Valid = false
		return result, nil
	}

	certs, ok := data["certs"].([]interface{})
	if !ok || len(certs) == 0 {
		result.Valid = false
		return result, nil
	}

	if issuer, ok := certs[0].(map[string]interface{})["issuerRaw"].(string); ok {
		result.Issuer = issuer
	}

	if protos, ok := data["protocols"].([]interface{}); ok {
		for _, p := range protos {
			if proto, ok := p.(map[string]interface{})["version"].(string); ok {
				result.Protocols = append(result.Protocols, proto)
			}
		}
	}

	result.Host = domain
	result.Valid = true

	return result, nil
}
