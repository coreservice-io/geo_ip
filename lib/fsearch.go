package lib

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Searcher struct {
	country_ip_map map[uint32]([]SORT_COUNTRY_IP)
}

func NewSearch() *Searcher {

	return &Searcher{}
}

func BytesToUint(b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)

	var x uint32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return uint32(x)
}

func (s *Searcher) LoadFile(path string,
	ip_convert func(ip string) (*big.Int, error)) error {

	///////////////////// country ipv4 //////////////////////////////////////////
	country_ip_f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer country_ip_f.Close()

	country_ip_scanner := bufio.NewScanner(country_ip_f)

	country_ip_map := make(map[uint32]([]SORT_COUNTRY_IP))

	var last_record *SORT_COUNTRY_IP
	var last_bucket_idx uint32
	last_record = GetEmptyCountryIP()
	last_bucket_idx = 0
	country_ip_map[0] = append(country_ip_map[0], *last_record)

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

		bucket_idx, _ := ExtractBucketIdx(record.Start_ip)

		if last_bucket_idx != bucket_idx {

			for idx := last_bucket_idx + 1; idx <= bucket_idx; idx++ {
				// fmt.Printf("fill %d use %d %s\n", idx, last_bucket_idx, last_record.Start_ip)
				country_ip_map[idx] = append(country_ip_map[idx], *last_record)
			}
		}

		country_ip_map[bucket_idx] = append(country_ip_map[bucket_idx], *record)

		last_bucket_idx = bucket_idx
		last_record = record
	}

	//////// sort  start ip desc ///////////////////
	for _, country_ip_list := range country_ip_map {
		sort.SliceStable(country_ip_list, func(i, j int) bool {
			return country_ip_list[i].Start_ip_score.Cmp(country_ip_list[j].Start_ip_score) == 1
		})
	}

	s.country_ip_map = country_ip_map
	return nil
}

func (s *Searcher) Search(target_ip string, target_ip_score *big.Int) *SORT_COUNTRY_IP {

	idx, _ := ExtractBucketIdx(target_ip)

	if group, ok := s.country_ip_map[idx]; !ok {
		return nil
	} else {
		country_index := sort.Search(len(group), func(j int) bool {
			return group[j].Start_ip_score.Cmp(target_ip_score) <= 0
		})

		if country_index >= 0 && country_index < len(group) {
			return &(group[country_index])
		}

		return nil
	}
}

var shiftIndex = []int{24, 16, 8, 0}

func CheckIP(ip string) (uint32, error) {
	var ps = strings.Split(strings.TrimSpace(ip), ".")
	if len(ps) != 4 {
		return 0, fmt.Errorf("invalid ip address `%s`", ip)
	}

	var val = uint32(0)
	for i, s := range ps {
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("the %dth part `%s` is not an integer", i, s)
		}

		if d < 0 || d > 255 {
			return 0, fmt.Errorf("the %dth part `%s` should be an integer bettween 0 and 255", i, s)
		}

		val |= uint32(d) << shiftIndex[i]
	}

	// convert the ip to integer
	return val, nil
}

func ExtractBucketIdx(ip string) (uint32, error) {
	var ps = strings.Split(strings.TrimSpace(ip), ".")
	if len(ps) != 4 {
		return 0, fmt.Errorf("invalid ip address `%s`", ip)
	}

	d0, _ := strconv.Atoi(ps[0])
	d1, _ := strconv.Atoi(ps[1])

	var val = uint32(0)
	val |= uint32(d0) << 8
	val |= uint32(d1)

	// convert the ip to integer
	return val, nil
}
