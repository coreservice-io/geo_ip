package lib

type GeoInfo struct {
	Ip string

	Country_code   string
	Country_name   string
	Continent_code string
	Continent_name string
	Region         string
	Latitude       float64
	Longitude      float64

	Asn           string
	Isp           string
	Is_datacenter bool
}

type (
	GeoIpInterface interface {
		GetInfo(ip string) (*GeoInfo, error)
		Upgrade() error
	}
)
