package lib

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/coreservice-io/geo_ip/data"
	"github.com/coreservice-io/package_client"
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
	City           string
	Latitude       float64
	Longitude      float64
}

type GeoIpClient struct {
	country_ipv4_list []SORT_COUNTRY_IP
	country_ipv6_list []SORT_COUNTRY_IP
	isp_ipv4_list     []SORT_ISP_IP
	isp_ipv6_list     []SORT_ISP_IP
	pc                *package_client.PackageClient
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

func line_parser_ip(line string, lineno int) (*SORT_COUNTRY_IP, error) {
	line_f := strings.ReplaceAll(line, "\\,", " ")
	line_split_array := strings.Split(line_f, ",")

	if _, exist := data.CountryList[line_split_array[1]]; !exist {
		return nil, fmt.Errorf("parser line err '%s'", line)
	}

	record := &SORT_COUNTRY_IP{
		Start_ip:       line_split_array[0],
		Start_ip_score: nil,
		Country_code:   line_split_array[1],
		Region:         line_split_array[2],
		City:           line_split_array[3],
	}

	if line_split_array[4] == "" || line_split_array[5] == "" {
		record.Latitude = 0
		record.Longitude = 0
	} else {
		if lati, err := strconv.ParseFloat(line_split_array[4], 64); err != nil {
			return nil, fmt.Errorf("parser line err '%s':%d. Err: %s", line, lineno, err.Error())
		} else {
			record.Latitude = lati
		}

		if longti, err := strconv.ParseFloat(line_split_array[5], 64); err != nil {
			return nil, fmt.Errorf("parser line err '%s':%d. Err: %s", line, lineno, err.Error())
		} else {
			record.Longitude = longti
		}
	}

	return record, nil
}

func line_parser_isp(line string, lineno int) (*SORT_ISP_IP, error) {
	line_split_array := strings.Split(line, ",")

	record := &SORT_ISP_IP{
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

	return record, nil
}

func init_country(country_abs_file string,
	ip_convert func(ip string) (*big.Int, error)) ([]SORT_COUNTRY_IP, error) {

	///////////////////// country ipv4 //////////////////////////////////////////
	country_ip_f, err := os.Open(country_abs_file)
	if err != nil {
		return nil, err
	}
	defer country_ip_f.Close()

	country_ip_list := []SORT_COUNTRY_IP{}

	country_ip_scanner := bufio.NewScanner(country_ip_f)

	line_no := 0
	for country_ip_scanner.Scan() {
		line_no = line_no + 1
		line := country_ip_scanner.Text()

		record, perr := line_parser_ip(line, line_no)
		if perr != nil {
			return nil, perr
		}

		/////////////
		ipint, err := ip_convert(record.Start_ip)
		if err != nil {
			return nil, err
		}
		record.Start_ip_score = ipint
		country_ip_list = append(country_ip_list, *record)
	}

	//////// sort  start ip desc ///////////////////
	sort.SliceStable(country_ip_list, func(i, j int) bool {
		return country_ip_list[i].Start_ip_score.Cmp(country_ip_list[j].Start_ip_score) == 1
	})

	return country_ip_list, nil
}

func init_isp(isp_abs_file string,
	ip_convert func(ip string) (*big.Int, error)) ([]SORT_ISP_IP, error) {

	///////////////////// country ipv4 //////////////////////////////////////////
	isp_ip_f, err := os.Open(isp_abs_file)
	if err != nil {
		return nil, err
	}
	defer isp_ip_f.Close()

	isp_ip_list := []SORT_ISP_IP{}

	isp_ip_scanner := bufio.NewScanner(isp_ip_f)

	line_no := 0
	for isp_ip_scanner.Scan() {
		line_no = line_no + 1
		line := isp_ip_scanner.Text()

		record, perr := line_parser_isp(line, line_no)
		if perr != nil {
			return nil, perr
		}

		/////////////
		ipint, err := ip_convert(record.Start_ip)
		if err != nil {
			return nil, err
		}
		record.Start_ip_score = ipint
		isp_ip_list = append(isp_ip_list, *record)
	}

	//////// sort  start ip desc ///////////////////
	sort.SliceStable(isp_ip_list, func(i, j int) bool {
		return isp_ip_list[i].Start_ip_score.Cmp(isp_ip_list[j].Start_ip_score) == 1
	})

	return isp_ip_list, nil
}

func (geoip_c *GeoIpClient) ReloadCsv(datafolder string,
	logger func(log_str string), err_logger func(err_log_str string)) error {

	country_ipv4_file_abs := filepath.Join(datafolder, "country_ipv4.csv")
	country_ipv6_file_abs := filepath.Join(datafolder, "country_ipv6.csv")
	isp_ipv4_file_abs := filepath.Join(datafolder, "isp_ipv4.csv")
	isp_ipv6_file_abs := filepath.Join(datafolder, "isp_ipv6.csv")

	////
	if country_ip_list, err := init_country(country_ipv4_file_abs, ip_convert_num_ipv4); err != nil {
		return err
	} else {
		if len(country_ip_list) == 0 {
			return errors.New("country_ipv4 len :0 ")
		}
		geoip_c.country_ipv4_list = country_ip_list
	}
	///
	if country_ip_list, err := init_country(country_ipv6_file_abs, ip_convert_num_ipv6); err != nil {
		return err
	} else {
		if len(country_ip_list) == 0 {
			return errors.New("country_ipv6 len :0 ")
		}
		geoip_c.country_ipv6_list = country_ip_list
	}

	///
	if isp_ip_list, err := init_isp(isp_ipv4_file_abs, ip_convert_num_ipv4); err != nil {
		return err
	} else {
		if len(isp_ip_list) == 0 {
			return errors.New("isp_ipv4 len :0 ")
		}
		geoip_c.isp_ipv4_list = isp_ip_list
	}
	///
	if isp_ip_list, err := init_isp(isp_ipv6_file_abs, ip_convert_num_ipv6); err != nil {
		return err
	} else {
		if len(isp_ip_list) == 0 {
			return errors.New("isp_ipv6 len :0 ")
		}
		geoip_c.isp_ipv6_list = isp_ip_list
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
	search_country := geoip_c.country_ipv4_list
	search_isp := geoip_c.isp_ipv4_list

	if ip_type == "ipv6" {
		search_country = geoip_c.country_ipv6_list
		search_isp = geoip_c.isp_ipv6_list
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

	country_index := sort.Search(len(search_country), func(j int) bool {
		return search_country[j].Start_ip_score.Cmp(target_ip_score) <= 0
	})

	if country_index >= 0 && country_index < len(search_country) {
		fillGeoInfo(&result, &search_country[country_index])
	}

	//
	isp_index := sort.Search(len(search_isp), func(j int) bool {
		return search_isp[j].Start_ip_score.Cmp(target_ip_score) <= 0
	})

	if isp_index >= 0 && isp_index < len(search_isp) {
		fillIspInfo(&result, &search_isp[isp_index])
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
