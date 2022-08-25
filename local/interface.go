package local

type GeoLocalInfo struct {
	Ip            string
	Country_code  string
	Region        string
	Latitude      float64
	Longitude     float64
	Is_datacenter bool
	Asn           string
}

type (
	GeoIpLocalI interface {
		GetLocalInfo(ip string) (*GeoLocalInfo, error)
	}
)
