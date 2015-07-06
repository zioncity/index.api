package main

import (
	"reflect"
	"strconv"

	"github.com/olivere/elastic"
)

type Equip struct {
	EquipId   uint32 `json:"equipid"`
	Gprs      string `json:"gprs"`
	Name      string `json:"name,omitempty"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Duration  uint32 `json:"duration"`  // Hour
	Lat       uint32 `json:"lat"`       // 129.11 * cmd //cmd=1000
	Lng       uint32 `json:"lng"`       // 36.87 *cmd
	Activated uint16 `json:"activated"` // 0 : unactivated, 1: activated
	Online    uint16 `json:"online"`    // 0:offline, 1:online
	Update    int64  `json:"update"`    // unixtime the last update time
	//	Antennas  []Antenna `json:"antenna,omitmepty"`
}

func equip_init() Equip {
	return Equip{Duration: 24}
}

func equip_get_id(equipid uint32) Equip {
	if v, ok := _id2equips[equipid]; ok {
		return *v
	}
	return equip_init()

}

func equip_drop_id(equipid uint32) {
	delete(_id2equips, equipid)

	equip_es_drop(equipid)
}
func equip_es_drop(equipid uint32) {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	_, err = client.Delete().Index(es_index).Id(strconv.Itoa(int(equipid))).Do()
	panic_error(err)
}
func equip_add(e Equip) Equip {
	orig := equip_get_id(e.EquipId)
	if e.EquipId != 0 && orig.EquipId == e.EquipId {
		orig.Name = select_string(orig.Name, e.Name)
		orig.Province = select_string(orig.Province, e.Province)
		orig.City = select_string(orig.City, e.City)
		e = orig
	}

	_id2equips[e.EquipId] = &e
	equip_es_upsert(e)
	antenna_upsert(e.EquipId, e.Gprs)
	return e
}
func select_string(p1, alt string) string {
	if alt != "" {
		p1 = alt
	}
	return p1
}
func equip_activate(e Equip) Equip {
	orig := equip_get_id(e.EquipId)
	if orig.EquipId != 0 && orig.EquipId == e.EquipId {
		//	orig.Gprs = e.Gprs
		e = orig
	}
	e.Activated = 1
	_id2equips[e.EquipId] = &e
	equip_es_upsert(e)
	antenna_upsert(e.EquipId, e.Gprs)
	return e
}

func equip_es_upsert(e Equip) {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	_, err = client.Index().Index(es_index).Type("equip").Id(strconv.Itoa(int(e.EquipId))).BodyJson(&e).Do()
	panic_error(err)
}

func equips_get_geo(province, city string) (ret []Equip) {
	for _, v := range _id2equips {
		if v.City == city && v.Province == province {
			ret = append(ret, *v)
		}
	}
	return
}

func equip_es_load() (ret []Equip) {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)

	for from, count := 0, 1000; count >= 1000; {
		result, err := client.Search().Index(es_index).Type("equip").From(from).Size(count).Do()
		panic_error(err)
		var v []Equip
		var te = reflect.TypeOf(Equip{})
		for _, item := range result.Each(te) {
			if e, ok := item.(Equip); ok {
				v = append(v, e)
			}
		}
		from, count = from+len(v), len(v)
		ret = append(ret, v...)
	}
	return ret
}
func equip_all() {
	var ret = equip_es_load()
	for _, equip := range ret {
		_id2equips[equip.EquipId] = &equip
	}
}
