package main

import (
	"reflect"

	"github.com/olivere/elastic"
)

type Alarm struct {
	EquipId uint32 `json:"equipid"`
	UnitId  uint32 `json:"unitid"`
	Gprs    string `json:"gprs"`
	H       int    `json:"h"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Z       int    `json:"z"`
	Date    int64  `json:"date"` //unixtime
}

type EquipAntennaAlarm struct {
	Equip   Equip   `json:"equip"`
	Antenna Antenna `json:"antenna"`
	Alarm   Alarm   `json:"alarm"`
}

var alarm_type = []string{"无", "下顷", "滚降", "方位", "高度"}

func alarm_default_type() string {
	return alarm_type[0]
}

func alarm_init() Alarm {
	return Alarm{Date: date_now()}
}

//获取下位机的告警
func alarms_get_equip(equipid uint32, from, count int) []EquipAntennaAlarm {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	var q = elastic.NewTermQuery("equipid", equipid)
	//  var q = elastic.NewTermFilter("equipid", equipid)

	result, err := client.Search().Index(es_index).Type("alarm").Query(q).Sort("date", false).From(from).Size(count).Do()
	panic_error(err)
	var v []EquipAntennaAlarm
	var ta = reflect.TypeOf(Alarm{})
	for _, item := range result.Each(ta) {
		if a, ok := item.(Alarm); ok {
			v = append(v, alarm_fill(a))
		}
	}
	return v
}

//获取某根天线的告警
func alarms_get_antenna(equipid, unitid uint32, from, count int) []EquipAntennaAlarm {
	client, err := elastic.NewClient(elastic.SetURL(es_url), elastic.SetSniff(false))
	panic_error(err)
	f := elastic.NewAndFilter()
	f.Add(elastic.NewTermFilter("equipid", equipid))
	f.Add(elastic.NewTermFilter("unitid", unitid))
	result, err := client.Search().Index(es_index).Type("alarm").Source(f.Source()).Sort("date", false).From(from).Size(count).Do()
	panic_error(err)
	var v []EquipAntennaAlarm
	var ta = reflect.TypeOf(Alarm{})
	for _, item := range result.Each(ta) {
		if a, ok := item.(Alarm); ok {
			v = append(v, alarm_fill(a))
		}
	}
	return v
}

//根据告警信息中的equipid填充下位机和天线数据
func alarm_fill(a Alarm) EquipAntennaAlarm {
	equip, antenna := equip_get_id(a.EquipId), antenna_get_id(a.EquipId, a.UnitId)
	return EquipAntennaAlarm{equip, antenna, a}
}

func alarm_update(a Alarm) EquipAntennaAlarm {
	antenna_set_alarm(a.H+a.X+a.Y+a.Z, a.EquipId, a.UnitId)
	return alarm_fill(a)
}
