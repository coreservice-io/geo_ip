package lib

import (
	"errors"
	"math/big"
	"net"
	"path/filepath"

	"github.com/coreservice-io/geo_ip/data"
	"github.com/coreservice-io/package_client"
)

const NEW_SEARCHER_IPV4 = true

type SORT_ISP_IP struct {
	Start_ip       string
	Start_ip_score *big.Int
	Asn            string
	Is_datacenter  bool
	Isp            string
}

type SORT_COUNTRY_IP struct {
	Start_ip       string
	Start_ip_score *big.Int
	Country_code   string
	Region         string
	City           string
	Latitude       float64
	Longitude      float64
}

type GeoIpClient struct {
	country_ipv4_searcher *CountrySearcher
	country_ipv6_searcher *CountrySearcher
	isp_ipv4_searcher     *IspSearcher
	isp_ipv6_searcher     *IspSearcher
	ipv4_searcher         *Searcher
	pc                    *package_client.PackageClient
}

func GetEmptyCountryIP() *SORT_COUNTRY_IP {
	return &SORT_COUNTRY_IP{
		"0.0.0.0",
		big.NewInt(0),
		"ZZ",
		"",
		"",
		0.000000,
		0.000000,
	}
}

// / the second int is 32 for ipv4 or 128 for ipv6
func IpToBigInt(ip net.IP) (*big.Int, error) {
	val := &big.Int{}
	val.SetBytes([]byte(ip))
	if len(ip) == net.IPv4len {
		return val, nil
	} else if len(ip) == net.IPv6len {
		return val, nil
	} else {
		return nil, errors.New("ip format error")
	}
}

func ip_convert_num_ipv4(ip string) (*big.Int, error) {
	ipv := net.ParseIP(ip)
	if ipv.To4() == nil {
		return nil, errors.New(ip + " is not ipv4")
	}
	return IpToBigInt(ipv)
}

func ip_convert_num_ipv6(ip string) (*big.Int, error) {
	ipv := net.ParseIP(ip)
	if ipv.To16() == nil {
		return nil, errors.New(ip + " is not ipv6")
	}
	return IpToBigInt(ipv)
}

func (geoip_c *GeoIpClient) ReloadCsv(datafolder string,
	logger func(log_str string), err_logger func(err_log_str string)) error {

	country_ipv4_file_abs := filepath.Join(datafolder, "country_ipv4.csv")
	country_ipv6_file_abs := filepath.Join(datafolder, "country_ipv6.csv")
	isp_ipv4_file_abs := filepath.Join(datafolder, "isp_ipv4.csv")
	isp_ipv6_file_abs := filepath.Join(datafolder, "isp_ipv6.csv")

	////
	if !NEW_SEARCHER_IPV4 {

		country_ipv4_searcher := NewCountrySearcher()

		if err := country_ipv4_searcher.LoadFile(country_ipv4_file_abs, ip_convert_num_ipv4); err != nil {
			return err
		} else {
			geoip_c.country_ipv4_searcher = country_ipv4_searcher
		}
	} else {
		ipv4_searcher := NewSearch()

		if err := ipv4_searcher.LoadFile(country_ipv4_file_abs, ip_convert_num_ipv4); err != nil {
			return err
		} else {
			geoip_c.ipv4_searcher = ipv4_searcher
		}
	}

	///
	country_ipv6_searcher := NewCountrySearcher()

	if err := country_ipv6_searcher.LoadFile(country_ipv6_file_abs, ip_convert_num_ipv6); err != nil {
		return err
	} else {
		geoip_c.country_ipv6_searcher = country_ipv6_searcher
	}

	///
	isp_ipv4_searcher := NewIspSearcher()

	if err := isp_ipv4_searcher.LoadFile(isp_ipv4_file_abs, ip_convert_num_ipv4); err != nil {
		return err
	} else {
		geoip_c.isp_ipv4_searcher = isp_ipv4_searcher
	}
	///
	isp_ipv6_searcher := NewIspSearcher()

	if err := isp_ipv6_searcher.LoadFile(isp_ipv6_file_abs, ip_convert_num_ipv6); err != nil {
		return err
	} else {
		geoip_c.isp_ipv6_searcher = isp_ipv6_searcher
	}

	return nil
}

func NewClient(update_key string, current_version string, datafolder string, ignore_data_exist bool,
	logger func(log_str string), err_logger func(err_log_str string)) (GeoIpInterface, error) {

	client := &GeoIpClient{}
	if !ignore_data_exist {
		load_err := client.ReloadCsv(datafolder, logger, err_logger)
		if load_err != nil {
			logger("load_err:" + load_err.Error())
			return nil, load_err
		}
	}
	///
	pc, err := StartAutoUpdate(update_key, current_version, false, datafolder, func() {
		client.ReloadCsv(datafolder, logger, err_logger)
	}, logger, err_logger)

	if err != nil {
		logger("StartAutoUpdate err:" + err.Error())
	}

	client.pc = pc
	////////////////////////
	return client, nil
}

func (geoip_c *GeoIpClient) Upgrade(ignore_version bool) error {
	return geoip_c.pc.Update(ignore_version)
}

func (geoip_c *GeoIpClient) GetInfo(target_ip string) (*GeoInfo, error) {

	// pre check ip
	if isLan, err := data.IsLanIp(target_ip); err != nil {
		return nil, err
	} else if isLan {
		return nil, errors.New("is lan ip")
	}

	ip_type := ""
	target_net_ip := net.ParseIP(target_ip)

	if target_net_ip.To4() != nil {
		ip_type = "ipv4"
	} else if target_net_ip.To16() != nil {
		ip_type = "ipv6"
	} else {
		return nil, errors.New("ip format error:" + target_ip)
	}

	target_ip_score, err := IpToBigInt(target_net_ip)
	if err != nil {
		return nil, err
	}

	//////////////
	search_country := geoip_c.country_ipv4_searcher
	search_isp := geoip_c.isp_ipv4_searcher

	if ip_type == "ipv6" {
		search_country = geoip_c.country_ipv6_searcher
		search_isp = geoip_c.isp_ipv6_searcher
	}
	////
	result := GeoInfo{
		Ip:             target_ip,
		Latitude:       0,
		Longitude:      0,
		Country_code:   data.NA,
		Country_name:   data.NA,
		Continent_code: data.NA,
		Continent_name: data.NA,
		Region:         data.NA,
		City:           data.NA,
		Asn:            data.NA,
		Isp:            data.NA,
		Is_datacenter:  false,
	}

	var country_info *SORT_COUNTRY_IP
	if NEW_SEARCHER_IPV4 && ip_type == "ipv4" {
		ipv4_searcher := geoip_c.ipv4_searcher
		country_info = ipv4_searcher.Search(target_ip, target_ip_score)
	} else {
		country_info = search_country.Search(target_ip_score)
	}
	isp_info := search_isp.Search(target_ip_score)

	if country_info != nil {
		fillGeoInfo(&result, country_info)
	}

	//
	if isp_info != nil {
		fillIspInfo(&result, isp_info)
	}

	return &result, nil
}

func fillGeoInfo(result *GeoInfo, info *SORT_COUNTRY_IP) {
	result.Latitude = info.Latitude
	result.Longitude = info.Longitude
	result.Country_code = info.Country_code
	result.Region = info.Region
	result.City = info.City

	if val, ok := data.CountryList[result.Country_code]; ok {
		result.Continent_code = val.ContinentCode
		result.Continent_name = val.ContinentName
		result.Country_name = val.CountryName
	}
}

func fillIspInfo(result *GeoInfo, info *SORT_ISP_IP) {
	result.Asn = info.Asn
	result.Isp = info.Isp
	result.Is_datacenter = info.Is_datacenter
}
