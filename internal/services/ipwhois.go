package services

import (
	"encoding/json"
	"infomiser/internal/models"
	"io"
	"net/http"
	"net/url"
)

type RawIPWhois struct {
	IP         string  `json:"ip"`
	Success    bool    `json:"success"`
	Type       string  `json:"type"`
	Country    string  `json:"country"`
	Region     string  `json:"region"`
	City       string  `json:"city"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Connection struct {
		Org    string `json:"org"`
		ISP    string `json:"isp"`
		Domain string `json:"domain"`
	} `json:"connection"`
	Timezone struct {
		ID string `json:"id"`
	} `json:"timezone"`
}

func GetIPWhois(ip string) (models.IPWhois, error) {
	var result models.IPWhois
	if ip == "" {
		return result, nil
	}

	resp, err := http.Get("https://ipwho.is/" + url.PathEscape(ip))
	if err != nil || resp.StatusCode != 200 {
		result.IP = ip
		return result, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var raw RawIPWhois
	json.Unmarshal(body, &raw)

	result.IP = raw.IP
	result.Country = raw.Country
	result.Region = raw.Region
	result.City = raw.City
	result.Lat = raw.Latitude
	result.Lon = raw.Longitude
	result.ISP = raw.Connection.ISP
	result.Domain = raw.Connection.Domain

	return result, nil
}
