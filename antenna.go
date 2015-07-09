package main

import (
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/olivere/elastic"
)

const ()

var antenna_type = []string{"reserve", "移动", "联通", "电信"}

type Antenna struct {
	EquipId     int64  `json:"equipid"`
	Gprs        string `json:"gprs"`
	UnitId      uint32 `json:"unitid"`  //[0,18)
	Lng         int    `json:"lng"`     // 129.23 *100000
	Lat         int    `json:"lat"`     // 32.11 * 100000
	Enable      int    `json:"enable"`  //0:disable, 1:enable
	Network     int    `json:"network"` //0:reserved 1:移动， 2:联通, 3:电信
	BaseId      int    `json:"baseid"`
	BaseName    string `json:"basename"`
	CellId      int    `json:"cellid"`
	BaseManu    string `json:"basemanu"`
	AntennaManu string `json:"antennamanu"`
	AntennaTyp  string `json:"antennatyp"`
	AlertStart  int64  `json:"alterstart"` //unixtime, 0 for no alert
	Alarm       int    `json:"alarm"`      // has alarm
	H           int    `json:"h"`          //cm 高度
	X           int    `json:"x"`          //cmd下倾角
	Y           int    `json:"y"`          //cmd方向角
	Z           int    `json:"z"`          //横滚角
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
func antenna_get_id(equipid int64, unitid uint32) Antenna {
	if val, ok := _id2antennas[equipid]; ok {
		return *val[unitid]
	}
	return antenna_init()
}

func antennas_get_equip(equipid int64) []EquipAntenna {
	equip := equip_get_id(equipid)
	var v []EquipAntenna
	if ats, ok := _id2antennas[equipid]; ok {
		for _, at := range ats {
			v = append(v, EquipAntenna{equip, *at})
		}
	}
	return v
}
func antenna_set_alarm(alarm int, equipid int64, unitid uint32) (ret int) {
	if v, ok := _id2antennas[equipid]; ok {
		ret = v[unitid].Alarm
		v[unitid].Alarm = alarm
		antenna_es_upsert(*v[unitid])
	}
	return ret
}
func antenna_set_attitude(a Attitude) {
	if v, ok := _id2antennas[a.EquipId]; ok {
		v[a.UnitId].attitude = a
	}
}
func antenna_initialize_equip(ats []*Antenna, gprs string, equipid int64) {
	for i := 0; i < len(ats); i++ {
		ats[i] = &Antenna{}
		ats[i].UnitId = uint32(i)
		ats[i].EquipId = equipid
		ats[i].Gprs = gprs
	}
}

func antenna_upsert(equipid int64, gprs string) {
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

func antenna_enable(equipid int64, unitid uint32, enable int) int {
	log.Println("antenna-disable", equipid, unitid, enable)
	ret := -1

	if v, ok := _id2antennas[equipid]; ok {
		v[unitid].Enable = enable
		antenna_es_upsert(*v[unitid])
		ret = 0
	}

	return ret
}

func antenna_update(a Antenna) {
	log.Println("antenna-upadte", a.EquipId, a.UnitId)
	if v, ok := _id2antennas[a.EquipId]; ok {
		v[a.UnitId].Lat = select_int(v[a.UnitId].Lat, a.Lat)
		v[a.UnitId].Lng = select_int(v[a.UnitId].Lng, a.Lng)
		v[a.UnitId].BaseName = select_string(v[a.UnitId].BaseName, a.BaseName)
		v[a.UnitId].BaseId = a.BaseId
		v[a.UnitId].CellId = a.CellId
		v[a.UnitId].BaseManu = select_string(v[a.UnitId].BaseManu, a.BaseManu)
		v[a.UnitId].AntennaManu = select_string(v[a.UnitId].AntennaManu, a.AntennaManu)
		v[a.UnitId].AntennaTyp = select_string(v[a.UnitId].AntennaTyp, a.AntennaTyp)
		v[a.UnitId].H = a.H
		v[a.UnitId].X = a.X
		v[a.UnitId].Y = a.Y
		v[a.UnitId].Z = a.Z
		v[a.UnitId].Enable = 1
		v[a.UnitId].Update = time.Now().Unix()
		antenna_es_upsert(*v[a.UnitId])
	}
}

func antenna_update_bias(equipid int64, unitid uint32, h, x, y, z int) {
	log.Println("antenna-update-bias", equipid, unitid, h, x, y, z)
	if v, ok := _id2antennas[equipid]; ok {
		ant := v[unitid]
		ant.H, ant.X, ant.Y, ant.Z = h, x, y, z
		ant.Update = time.Now().Unix()
		antenna_es_upsert(*ant)
	}
}
func antenna_es_upsert(a Antenna) {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	_, err = client.Index().Index(es_index).Type("antenna").Id(strconv.Itoa(int(a.EquipId*100 + int64(a.UnitId)))).BodyJson(&a).Do()
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
	log.Println("antenna-es-load", len(ret))
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
