package models

type IPInfo struct {
	IP         string            `json:"ip"`
	Type       string            `json:"type"`
	Country    string            `json:"country"`
	Region     string            `json:"region"`
	City       string            `json:"city"`
	Lat        float64           `json:"latitude"`
	Lon        float64           `json:"longitude"`
	Connection map[string]string `json:"connection"`
	Timezone   map[string]string `json:"timezone"`
}

type SubdomainInfo struct {
	Domain     string
	Subdomains []string
}
