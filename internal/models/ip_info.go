package models

type IPInfo struct {
	IP            string  `json:"ip"`
	Success       bool    `json:"success"`
	Type          string  `json:"type"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	City          string  `json:"city"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	IsEU          bool    `json:"is_eu"`
	Postal        string  `json:"postal"`
	CallingCode   string  `json:"calling_code"`
	Capital       string  `json:"capital"`
	Borders       string  `json:"borders"`

	Flag struct {
		Img          string `json:"img"`
		Emoji        string `json:"emoji"`
		EmojiUnicode string `json:"emoji_unicode"`
	} `json:"flag"`

	Connection struct {
		ASN    int    `json:"asn"`
		Org    string `json:"org"`
		ISP    string `json:"isp"`
		Domain string `json:"domain"`
	} `json:"connection"`

	Timezone struct {
		ID           string `json:"id"`
		Abbreviation string `json:"abbr"`
		IsDst        bool   `json:"is_dst"`
		Offset       int    `json:"offset"`
		UTC          string `json:"utc"`
		CurrentTime  string `json:"current_time"`
	} `json:"timezone"`
}
