package main

import "strconv"

var (
	es_url   = "http://testbox02.chinacloudapp.cn:9200"
	es_index = "fixed"
)

func atoui32(x string) uint32 {
	v, _ := strconv.Atoi(x)
	return uint32(v)
}

func atoi(x string) int64 {
	v, _ := strconv.Atoi(x)
	return int64(v)
}
func gprs2id(g string) int64 {
	return atoi(g)
}

func id2gprs(id int64) string {
	return strconv.Itoa(int(id))
}

func select_int(p1, p2 int) int {
	if p2 != 0 {
		p1 = p1
	}
	return p1
}
