package main

import (
	"math/rand"
	"time"

	"github.com/olivere/elastic"
)

var letters = []rune("1234567890")

func random_no(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	cmd        = 1000 //degree / 1000
	cm         = 100
	unit_count = 18
)

var equip_id_seed uint32 = 0x80000000

func date_now() int64 {
	return time.Now().Unix()
	//return time.Now().Add(-time.Hour * time.Duration(rand.Intn(24))).Unix()
}

var _profile Profile = Profile{24, 1 * cm, 2 * cmd, 5 * cmd, 5 * cmd}

var _id2equips map[int64]*Equip
var _id2antennas map[int64][]*Antenna
var _gprs2equipid map[string]uint32
var _id2gprs map[uint32]string

func es_upsert(typ string, v interface{}) {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	_, err = client.Index().Index(es_index).Type(typ).BodyJson(v).Do()
	panic_error(err)
}
