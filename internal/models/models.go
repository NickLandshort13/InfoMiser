package models

type IPWhois struct {
	IP      string
	Country string
	Region  string
	City    string
	Lat     float64
	Lon     float64
	ISP     string
	Domain  string
}

type Subdomains struct {
	Domain     string
	Subdomains []string
}

type SSLScan struct {
	Host      string
	Valid     bool
	Protocols []string
	Issuer    string
}

type Shodan struct {
	IP    string
	Ports []int
}

type Pwned struct {
	Domain   string
	Breaches []string
}
