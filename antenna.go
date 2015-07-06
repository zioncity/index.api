package main

import (
	"reflect"
	"strconv"

	"github.com/olivere/elastic"
)

const ()

var antenna_type = []string{"reserve", "移动", "联通", "电信"}

type Antenna struct {
	EquipId     uint32 `json:"equipid"`
	Gprs        string `json:"gprs"`
	UnitId      uint32 `json:"unitid"`  //[0,18)
	Lng         int    `json:"lng"`     // 129.23 *100000
	Lat         int    `json:"lat"`     // 32.11 * 100000
	Disable     uint32 `json:"disable"` //1:disable, 0:not disable
	Network     int    `json:"network"` //0:reserved 1:移动， 2:联通, 3:电信
	BaseId      uint32 `json:"baseid"`
	Basename    string `json:"basename"`
	CellId      uint16 `json:"cellid"`
	BaseManu    string `json:"basemanu"`
	AntennaManu string `json:"antennamanu"`
	AntennaTyp  string `json:"antennatyp"`
	AlertStart  int64  `json:"alterstart"` //unixtime, 0 for no alert
	Alarm       int    `json:"alarm"`      // has alarm
	H           uint32 `json:"h"`          //cm 高度
	X           uint32 `json:"x"`          //cmd下倾角
	Y           uint32 `json:"y"`          //cmd方向角
	Z           uint32 `json:"z"`          //横滚角
	Update      int64  `json:"update"`     //unixtime the last update
	attitude    Attitude
}

type EquipAntenna struct {
	Equip   Equip   `json:"equip"`
	Antenna Antenna `json:"antenna"`
}

func antenna_init() Antenna {
	return Antenna{}
}

//unitid < unit_count
func antenna_get_id(equipid uint32, unitid uint32) Antenna {
	var val []*Antenna
	var ok bool
	if val, ok = _id2antennas[equipid]; ok {
		return *val[unitid]
	}
	return antenna_init()
}

func antennas_get_equip(equipid uint32) []EquipAntenna {
	equip := equip_get_id(equipid)
	var v []EquipAntenna
	if ats, ok := _id2antennas[equipid]; ok {
		for _, at := range ats {
			v = append(v, EquipAntenna{equip, *at})
		}
	}
	return v
}
func antenna_set_alarm(alarm int, equipid, unitid uint32) {
	if v, ok := _id2antennas[equipid]; ok {
		v[unitid].Alarm = alarm
	}
}
func antenna_set_attitude(a Attitude) {
	if v, ok := _id2antennas[a.EquipId]; ok {
		v[a.UnitId].attitude = a
	}
}
func antenna_initialize_equip(ats []*Antenna, gprs string, equipid uint32) {
	for i := 0; i < len(ats); i++ {
		ats[i] = &Antenna{}
		ats[i].UnitId = uint32(i)
		ats[i].EquipId = equipid
		ats[i].Gprs = gprs
	}
}

func antenna_upsert(equipid uint32, gprs string) {
	if _, ok := _id2antennas[equipid]; ok {
		return
	}
	ats := make([]*Antenna, unit_count)
	antenna_initialize_equip(ats, gprs, equipid)

	_id2antennas[equipid] = ats
	for _, a := range ats {
		antenna_es_upsert(*a)
	}
}

func antenna_disable(equipid, unitid, disable uint32) int {
	ret := -1
	if unitid >= unit_count {
		return ret
	}
	if v, ok := _id2antennas[equipid]; ok {
		v[unitid].Disable = disable
		antenna_es_upsert(*v[unitid])
		ret = 0
	}

	return ret
}
func antenna_es_upsert(a Antenna) {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	_, err = client.Index().Index(es_index).Type("antenna").Id(strconv.Itoa(int(a.EquipId*100 + a.UnitId))).BodyJson(&a).Do()
	panic_error(err)
}
func antennas_es_load() (ret []Antenna) {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	for from, count := 0, 1000; count >= 1000; {
		result, err := client.Search().Index(es_index).Type("antenna").From(from).Size(count).Do()
		panic_error(err)
		var v []Antenna
		var ta = reflect.TypeOf(Antenna{})
		for _, item := range result.Each(ta) {
			if a, ok := item.(Antenna); ok {
				v = append(v, a)
			}
		}
		from, count = from+len(v), len(v)
		ret = append(ret, v...)
	}
	return ret
}

//var _id2antennas map[uint32][]*Antenna
func antennas_all() {
	ats := antennas_es_load()
	for _, at := range ats {
		if x, ok := _id2antennas[at.EquipId]; ok {
			x[at.UnitId] = &at
		} else {
			units := make([]*Antenna, unit_count)
			antenna_initialize_equip(units, at.Gprs, at.EquipId)
			_id2antennas[at.EquipId] = units
			units[at.UnitId] = &at
		}
	}
}
