package lib

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/coreservice-io/geo_ip/data"
)

type CountrySearcher struct {
	country_ip_list []SORT_COUNTRY_IP
}

func NewCountrySearcher() *CountrySearcher {
	return &CountrySearcher{}
}

func (s *CountrySearcher) LoadFile(country_abs_file string,
	ip_convert func(ip string) (*big.Int, error)) error {

	///////////////////// country ipv4 //////////////////////////////////////////
	country_ip_f, err := os.Open(country_abs_file)
	if err != nil {
		return err
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
			return perr
		}

		if record == nil {
			continue
		}

		/////////////
		ipint, err := ip_convert(record.Start_ip)
		if err != nil {
			return err
		}
		record.Start_ip_score = ipint
		country_ip_list = append(country_ip_list, *record)
	}

	//////// sort  start ip desc ///////////////////
	sort.SliceStable(country_ip_list, func(i, j int) bool {
		return country_ip_list[i].Start_ip_score.Cmp(country_ip_list[j].Start_ip_score) == 1
	})

	if len(country_ip_list) == 0 {
		return errors.New("country_ipv4 len :0 ")
	}

	s.country_ip_list = country_ip_list
	return nil
}

func (s *CountrySearcher) Search(target_ip_score *big.Int) *SORT_COUNTRY_IP {

	c_len := len(s.country_ip_list)
	country_index := sort.Search(c_len, func(j int) bool {
		return s.country_ip_list[j].Start_ip_score.Cmp(target_ip_score) <= 0
	})

	if country_index >= 0 && country_index < c_len {
		return &(s.country_ip_list[country_index])
	}

	return nil
}

func line_parser_ip(line string, lineno int) (*SORT_COUNTRY_IP, error) {
	line_f := strings.ReplaceAll(line, "\\,", " ")
	line_split_array := strings.Split(line_f, ",")

	if line_split_array[1] == "" {
		return nil, nil
	}
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
