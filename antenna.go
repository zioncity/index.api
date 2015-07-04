package main

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
	H           uint32 `json:"h"`          //cm
	X           uint32 `json:"x"`          //cmd
	Y           uint32 `json:"y"`          //cmd
	Z           uint32 `json:"z"`
	Update      int64  `json:"update"` //unixtime the last update
}
type EquipAntenna struct {
	Equip   Equip   `json:"equip"`
	Antenna Antenna `json:"antenna"`
}

func antenna_init() Antenna {
	return Antenna{}
}

func antenna_get_id(equipid uint32, unitid uint32) Antenna {
	if unitid >= unit_count {
		return antenna_init()
	}
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
	if unitid >= unit_count {
		return
	}
	if v, ok := _id2antennas[equipid]; ok {
		v[unitid].Alarm = alarm
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
	ats := make([]*Antenna, unit_count)
	antenna_initialize_equip(ats, gprs, equipid)
	post_index_antennas(ats)
	_id2antennas[equipid] = ats
}

func antenna_disable(equipid, unitid, disable uint32) int {
	ret := -1
	if unitid >= unit_count {
		return ret
	}
	if v, ok := _id2antennas[equipid]; ok {
		v[unitid].Disable = disable
		ret = 0
	}
	return ret
}
