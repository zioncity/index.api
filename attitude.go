package main

import (
	"math/rand"
	"time"

	"github.com/olivere/elastic"
)

type Attitude struct {
	EquipId uint32 `json:"equipid"`
	Gprs    string `json:"gprs"`
	UnitId  uint32 `json:"unitid"`
	H       int    `json:"h"`      //cm centimetre
	X       int    `json:"x"`      //degrees * 100000
	Y       int    `json:"y"`      //degrees * 100000
	Z       int    `json:"z"`      //degrees * 100000
	Update  int64  `json:"update"` // unixtime
}

type EquipAntennaAttitude struct {
	Equip    Equip    `json:"equip"`
	Antenna  Antenna  `json:"antenna"`
	Attitude Attitude `json:"attitude"`
}

func attitudes_fresh_antenna(atts []Attitude) {
	for _, a := range atts {
		a.Update = time.Now().Unix()
		antenna_set_attitude(a)
	}
}

//天线姿态上报，分批数据
func attitudes_update(atts []Attitude) []EquipAntennaAttitude {
	attitudes_fresh_antenna(atts)
	attitudes_insert(atts)
	v := make([]EquipAntennaAttitude, len(atts))
	for i := 0; i < len(atts); i++ {
		atts[i].Update = time.Now().Unix()
		if atts[i].Gprs == "" {
			atts[i].Gprs = random_no(16)
		}
		if atts[i].EquipId == 0 {
			atts[i].EquipId = rand.Uint32()
		}
		v[i].Antenna = antenna_init()
		v[i].Equip = random_equip()
		v[i].Attitude = atts[i]
	}
	return v
}
func attitudes_insert(atts []Attitude) {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	for _, attitude := range atts {
		_, err = client.Index().Index(es_index).Type("attitude").BodyJson(&attitude).Do()
		panic_error(err)
	}
}

func attitudes_get(equipid, unitid uint32, from, count int) []EquipAntennaAttitude {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	f := elastic.NewAndFilter()
	f.Add(elastic.NewTermFilter("equipid", equipid))
	f.Add(elastic.NewTermFilter("unitid", unitid))
	//	var q = elastic.NewTermFilter("equipid", equipid).Source()

	result, err := client.Search().Index(es_index).Type("attitude").Source(f.Source()).Sort("update", false).From(from).Size(count).Do()

	return nil
}

func attitudes_get_equip(equipid uint32) (v []EquipAntennaAttitude) {
	antennas := antennas_get_equip(equipid)

	return nil
}
