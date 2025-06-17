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

type Shodan struct {
	IP    string
	Ports []int
}

type Pwned struct {
	Domain   string
	Breaches []PwnedBreach
	Error    string
}

type PwnedBreach struct {
	Name        string
	Title       string
	BreachDate  string
	Description string
	PwnCount    int
	DataClasses []string
	LogoPath    string
	IsVerified  bool
}
