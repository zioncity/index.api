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
}

var _profile Profile = Profile{24, 1 * cm, 2 * cmd, 5 * cmd, 5 * cmd}

var (
	_id2equips    = make(map[int64]*Equip)
	_id2antennas  = make(map[int64][]*Antenna)
	_gprs2equipid = make(map[string]uint32)
	_id2gprs      = make(map[uint32]string)
)

func init() {

}
func es_upsert(typ string, v interface{}) {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	_, err = client.Index().Index(es_index).Type(typ).BodyJson(v).Do()
	panic_error(err)
}
