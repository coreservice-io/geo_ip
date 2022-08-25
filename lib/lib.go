package lib

import (
	"bufio"
	"errors"
	"math/big"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/coreservice-io/geo_ip/data"
)

type SORT_GEO_IP struct {
	Start_ip       string
	Start_ip_score *big.Int
	Country_code   string
	Region         string
	Latitude       float64
	Longitude      float64
	Is_datacenter  bool
	Asn            string
}

type GeoIpClient struct {
	geo_ip_list []SORT_GEO_IP
}

/// the second int is 32 for ipv4 or 128 for ipv6
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

func NewClient(ip_geo_file_abs string) (GeoIpInterface, error) {

	client := &GeoIpClient{}

	///////////////////////////////////////////////////////////////
	ip_asn_d_f, err := os.Open(ip_geo_file_abs)
	if err != nil {
		return nil, err
	}
	defer ip_asn_d_f.Close()

	ip_asn_d_scanner := bufio.NewScanner(ip_asn_d_f)

	for ip_asn_d_scanner.Scan() {

		line := ip_asn_d_scanner.Text()
		line_split_array := strings.Split(line, ",")

		record := SORT_GEO_IP{
			Asn:            line_split_array[0],
			Start_ip:       line_split_array[1],
			Start_ip_score: nil,
			Is_datacenter:  false,
			Country_code:   line_split_array[3],
			Region:         line_split_array[4],
			Latitude:       0,
			Longitude:      0,
		}

		/////////////

		ipint, err := IpToBigInt(net.ParseIP(line_split_array[1]))
		if err != nil {
			return nil, err
		}

		record.Start_ip_score = ipint

		/////////////

		if line_split_array[2] == "1" {
			record.Is_datacenter = true
		} else {
			record.Is_datacenter = false
		}

		/////////////
		lati, err := strconv.ParseFloat(line_split_array[5], 64)
		if err != nil {
			return nil, err
		}

		longi, err := strconv.ParseFloat(line_split_array[6], 64)
		if err != nil {
			return nil, err
		}

		record.Latitude = lati
		record.Longitude = longi

		client.geo_ip_list = append(client.geo_ip_list, record)

	}

	if len(client.geo_ip_list) == 0 {
		return nil, errors.New("geo_ip_list len :0 ")
	}

	//////// sort  start ip desc ///////////////////
	sort.SliceStable(client.geo_ip_list, func(i, j int) bool {
		return client.geo_ip_list[i].Start_ip_score.Cmp(client.geo_ip_list[j].Start_ip_score) == 1
	})

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

	target_ip_score, err := IpToBigInt(net.ParseIP(target_ip))
	if err != nil {
		return nil, err
	}

	result := GeoInfo{
		Ip:            target_ip,
		Latitude:      0,
		Longitude:     0,
		Country_code:  data.NA,
		Region:        data.NA,
		Asn:           data.NA,
		Is_datacenter: false,
	}

	//////////////
	index := sort.Search(len(i.geo_ip_list), func(j int) bool {
		return i.geo_ip_list[j].Start_ip_score.Cmp(target_ip_score) <= 0
	})

	if index >= 0 && index < len(i.geo_ip_list) {
		result.Asn = i.geo_ip_list[index].Asn
		result.Is_datacenter = i.geo_ip_list[index].Is_datacenter
		result.Latitude = i.geo_ip_list[index].Latitude
		result.Longitude = i.geo_ip_list[index].Longitude
		result.Country_code = i.geo_ip_list[index].Country_code
		result.Region = i.geo_ip_list[index].Region
	}

	return &result, nil

}
