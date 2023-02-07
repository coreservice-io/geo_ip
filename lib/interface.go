package lib

type GeoInfo struct {
	Ip string

	Country_code   string
	Country_name   string
	Continent_code string
	Continent_name string
	Region         string
	City           string
	Latitude       float64
	Longitude      float64

	Asn           string
	Isp           string
	Is_datacenter bool
}

type (
	GeoIpInterface interface {
		GetInfo(ip string) (*GeoInfo, error)
		Upgrade(ignore_version bool) error //when ignore_version is true , version will be neglected,always redownload and upgrade to lastest version
	}
)
