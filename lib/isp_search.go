package lib

import (
	"bufio"
	"errors"
	"math/big"
	"os"
	"sort"
	"strings"
)

type IspSearcher struct {
	isp_ip_list []SORT_ISP_IP
}

func NewIspSearcher() *IspSearcher {
	return &IspSearcher{}
}

func (s *IspSearcher) LoadFile(isp_abs_file string, ip_convert func(ip string) (*big.Int, error)) error {

	isp_ip_f, err := os.Open(isp_abs_file)
	if err != nil {
		return err
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
			return perr
		}

		/////////////
		ipint, err := ip_convert(record.Start_ip)
		if err != nil {
			return err
		}
		record.Start_ip_score = ipint
		isp_ip_list = append(isp_ip_list, *record)
	}

	//////// sort  start ip desc ///////////////////
	sort.SliceStable(isp_ip_list, func(i, j int) bool {
		return isp_ip_list[i].Start_ip_score.Cmp(isp_ip_list[j].Start_ip_score) == 1
	})

	if len(isp_ip_list) == 0 {
		return errors.New("isp_ipv6 len :0 ")
	}

	s.isp_ip_list = isp_ip_list
	return nil
}

func (s *IspSearcher) Search(target_ip_score *big.Int) *SORT_ISP_IP {

	country_index := sort.Search(len(s.isp_ip_list), func(j int) bool {
		return s.isp_ip_list[j].Start_ip_score.Cmp(target_ip_score) <= 0
	})

	if country_index >= 0 && country_index < len(s.isp_ip_list) {
		return &(s.isp_ip_list[country_index])
	}

	return nil
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
