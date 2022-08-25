package lib

type GeoInfo struct {
	Ip            string
	Country_code  string
	Region        string
	Latitude      float64
	Longitude     float64
	Is_datacenter bool
	Asn           string
}

type (
	GeoIpInterface interface {
		GetInfo(ip string) (*GeoInfo, error)
	}
)
