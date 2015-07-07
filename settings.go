package main

import "strconv"

var (
	es_url   = "http://"
	es_index = "fixed"
)

func atoui32(x string) uint32 {
	v, _ := strconv.Atoi(x)
	return uint32(v)
}

func gprs2id(g string) uint32 {
	return atoui32(g)
}

func id2gprs(id uint32) string {
	return strconv.Itoa(int(id))
}

func select_int(p1, p2 int) int {
	if p2 != 0 {
		p1 = p1
	}
	return p1
}
