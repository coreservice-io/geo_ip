package lib

import (
	"bufio"
	"errors"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/coreservice-io/geo_ip/data"
)

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
	Latitude       float64
	Longitude      float64
}

type GeoIpClient struct {
	country_ipv4_list []SORT_COUNTRY_IP
	country_ipv6_list []SORT_COUNTRY_IP
	isp_ipv4_list     []SORT_ISP_IP
	isp_ipv6_list     []SORT_ISP_IP
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

// iptype ="ipv4" or "ipv6"
func (geoip_c *GeoIpClient) init_country(country_abs_file string, ip_type string) error {

	if ip_type != "ipv4" && ip_type != "ipv6" {
		return errors.New("ip_type error ,only 'ipv4' or 'ipv6' allowed")
	}

	///////////////////// country ipv4 //////////////////////////////////////////
	country_ip_f, err := os.Open(country_abs_file)
	if err != nil {
		return err
	}
	defer country_ip_f.Close()

	country_ipv4_list := []SORT_COUNTRY_IP{}
	country_ipv6_list := []SORT_COUNTRY_IP{}

	country_ip_scanner := bufio.NewScanner(country_ip_f)

	for country_ip_scanner.Scan() {

		line := country_ip_scanner.Text()
		line_split_array := strings.Split(line, ",")

		record := SORT_COUNTRY_IP{
			Start_ip:       line_split_array[0],
			Start_ip_score: nil,
			Country_code:   line_split_array[1],
			Region:         line_split_array[2],
		}

		if lati, err := strconv.ParseFloat(line_split_array[3], 64); err != nil {
			return err
		} else {
			record.Latitude = lati
		}

		if longti, err := strconv.ParseFloat(line_split_array[4], 64); err != nil {
			return err
		} else {
			record.Longitude = longti
		}

		/////////////
		if ip_type == "ipv4" {
			ipv4_ := net.ParseIP(line_split_array[0])
			if ipv4_.To4() == nil {
				return errors.New(line_split_array[0] + "is not ipv4 [country]")
			}
			ipint, err := IpToBigInt(ipv4_)
			if err != nil {
				return err
			}
			record.Start_ip_score = ipint
			country_ipv4_list = append(country_ipv4_list, record)
		} else {

			ipv6_ := net.ParseIP(line_split_array[0])
			if ipv6_.To16() == nil {
				return errors.New(line_split_array[0] + "is not ipv6 [country]")
			}
			ipint, err := IpToBigInt(ipv6_)
			if err != nil {
				return err
			}
			record.Start_ip_score = ipint
			country_ipv6_list = append(country_ipv6_list, record)
		}
	}

	if ip_type == "ipv4" {
		//////// sort  start ip desc ///////////////////
		sort.SliceStable(country_ipv4_list, func(i, j int) bool {
			return country_ipv4_list[i].Start_ip_score.Cmp(country_ipv4_list[j].Start_ip_score) == 1
		})
		////
		if len(country_ipv4_list) == 0 {
			return errors.New("country_ipv4_list len :0 ")
		}

		geoip_c.country_ipv4_list = country_ipv4_list

	} else {
		//////// sort  start ip desc ///////////////////
		sort.SliceStable(country_ipv6_list, func(i, j int) bool {
			return country_ipv6_list[i].Start_ip_score.Cmp(country_ipv6_list[j].Start_ip_score) == 1
		})
		if len(country_ipv6_list) == 0 {
			return errors.New("country_ipv6_list len :0 ")
		}

		geoip_c.country_ipv6_list = country_ipv6_list
	}

	return nil
}

// iptype ="ipv4" or "ipv6"
func (geoip_c *GeoIpClient) init_isp(isp_abs_file string, ip_type string) error {

	if ip_type != "ipv4" && ip_type != "ipv6" {
		return errors.New("ip_type error ,only 'ipv4' or 'ipv6' allowed")
	}

	///////////////////// country ipv4 //////////////////////////////////////////
	isp_ip_f, err := os.Open(isp_abs_file)
	if err != nil {
		return err
	}
	defer isp_ip_f.Close()

	isp_ipv4_list := []SORT_ISP_IP{}
	isp_ipv6_list := []SORT_ISP_IP{}

	isp_ip_scanner := bufio.NewScanner(isp_ip_f)

	for isp_ip_scanner.Scan() {

		line := isp_ip_scanner.Text()
		line_split_array := strings.Split(line, ",")

		record := SORT_ISP_IP{
			Start_ip:       line_split_array[0],
			Start_ip_score: nil,
			Asn:            line_split_array[1],
			Isp:            line_split_array[3],
		}

		if strings.Trim(line_split_array[2], " ") == "1" {
			record.Is_datacenter = true
		} else {
			record.Is_datacenter = false
		}

		/////////////
		if ip_type == "ipv4" {
			ipv4_ := net.ParseIP(record.Start_ip)
			if ipv4_.To4() == nil {
				return errors.New(record.Start_ip + " is not ipv4 [isp]")
			}
			ipint, err := IpToBigInt(ipv4_)
			if err != nil {
				return err
			}
			record.Start_ip_score = ipint
			isp_ipv4_list = append(isp_ipv4_list, record)
		} else {
			ipv6_ := net.ParseIP(record.Start_ip)
			if ipv6_.To16() == nil {
				return errors.New(record.Start_ip + " is not ipv6 [isp]")
			}
			ipint, err := IpToBigInt(ipv6_)
			if err != nil {
				return err
			}
			record.Start_ip_score = ipint
			isp_ipv6_list = append(isp_ipv6_list, record)
		}
	}

	if ip_type == "ipv4" {
		//////// sort  start ip desc ///////////////////
		sort.SliceStable(isp_ipv4_list, func(i, j int) bool {
			return isp_ipv4_list[i].Start_ip_score.Cmp(isp_ipv4_list[j].Start_ip_score) == 1
		})
		if len(isp_ipv4_list) == 0 {
			return errors.New("isp_ipv4_list len :0 ")
		}
		geoip_c.isp_ipv4_list = isp_ipv4_list
	} else {
		//////// sort  start ip desc ///////////////////
		sort.SliceStable(isp_ipv6_list, func(i, j int) bool {
			return isp_ipv6_list[i].Start_ip_score.Cmp(isp_ipv6_list[j].Start_ip_score) == 1
		})
		if len(isp_ipv6_list) == 0 {
			return errors.New("isp_ipv6_list len :0 ")
		}
		geoip_c.isp_ipv6_list = isp_ipv6_list
	}

	return nil
}

func (gip_client *GeoIpClient) ReloadCsv(datafolder string) error {

	country_ipv4_file_abs := filepath.Join(datafolder, "country_ipv4.csv")
	country_ipv6_file_abs := filepath.Join(datafolder, "country_ipv6.csv")
	isp_ipv4_file_abs := filepath.Join(datafolder, "isp_ipv4.csv")
	isp_ipv6_file_abs := filepath.Join(datafolder, "isp_ipv6.csv")

	////
	err := gip_client.init_country(country_ipv4_file_abs, "ipv4")
	if err != nil {
		return err
	}
	///
	err = gip_client.init_country(country_ipv6_file_abs, "ipv6")
	if err != nil {
		return err
	}
	///
	err = gip_client.init_isp(isp_ipv4_file_abs, "ipv4")
	if err != nil {
		return err
	}
	///
	err = gip_client.init_isp(isp_ipv6_file_abs, "ipv6")
	if err != nil {
		return err
	}

	//fmt.Println(len(gip_client.country_ipv4_list))

	return nil
}

func NewClient(datafolder string, init_update bool, logger func(string)) (GeoIpInterface, error) {

	client := &GeoIpClient{}
	load_err := client.ReloadCsv(datafolder)
	if load_err != nil {
		logger("load_err:" + load_err.Error())
		return nil, load_err
	}
	///
	StartAutoUpdate(false, init_update, datafolder, func() {
		client.ReloadCsv(datafolder)
	}, logger)

	////////////////////////
	return client, nil
}

func (i *GeoIpClient) GetInfo(target_ip string) (*GeoInfo, error) {

	//pre check ip
	isLan, err := data.IsLanIp(target_ip)
	if err != nil {
		return nil, err
	}
	if isLan {
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

	result := GeoInfo{
		Ip:             target_ip,
		Latitude:       0,
		Longitude:      0,
		Country_code:   data.NA,
		Country_name:   data.NA,
		Continent_code: data.NA,
		Continent_name: data.NA,
		Region:         data.NA,
		Asn:            data.NA,
		Isp:            data.NA,
		Is_datacenter:  false,
	}

	//////////////
	search_country := i.country_ipv4_list
	search_isp := i.isp_ipv4_list

	if ip_type == "ipv6" {
		search_country = i.country_ipv6_list
		search_isp = i.isp_ipv6_list
	}
	////

	// fmt.Println("ip_type", ip_type)
	// fmt.Println("target_ip_score", target_ip_score)
	// fmt.Println("search_country", len(search_country))

	country_index := sort.Search(len(search_country), func(j int) bool {
		return search_country[j].Start_ip_score.Cmp(target_ip_score) <= 0
	})

	//fmt.Println(country_index)

	if country_index >= 0 && country_index < len(search_country) {
		result.Latitude = search_country[country_index].Latitude
		result.Longitude = search_country[country_index].Longitude
		result.Country_code = search_country[country_index].Country_code
		result.Region = search_country[country_index].Region

		if val, ok := data.CountryList[result.Country_code]; ok {
			result.Continent_code = val.ContinentCode
			result.Continent_name = val.ContinentName
			result.Country_name = val.CountryName
		}
	}

	//
	isp_index := sort.Search(len(search_isp), func(j int) bool {
		return search_isp[j].Start_ip_score.Cmp(target_ip_score) <= 0
	})

	if isp_index >= 0 && isp_index < len(search_isp) {
		result.Asn = search_isp[isp_index].Asn
		result.Isp = search_isp[isp_index].Isp
		result.Is_datacenter = search_isp[isp_index].Is_datacenter
	}

	return &result, nil
}
