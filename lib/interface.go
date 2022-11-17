package lib

type GeoInfo struct {
	Ip string

	Region        string
	Isp           string
	Latitude      float64
	Longitude     float64
	Is_datacenter bool
	Asn           string

	Country_code   string
	Country_name   string
	Continent_code string
	Continent_name string
}

type (
	GeoIpInterface interface {
		GetInfo(ip string) (*GeoInfo, error)
	}
)
