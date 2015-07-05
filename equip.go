package main

import "strconv"

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
}

func random_equip() Equip {
	return Equip{Duration: 24}
}

func random_equips(n int) []Equip {
	v := make([]Equip, n)
	for i := 0; i < n; i++ {
		v[i] = random_equip()
	}
	return v
}

func equip_save(e Equip) error {
	return nil
}

func equip_get_id(equipid uint32) Equip {
	if v, ok := _id2equips[equipid]; ok {
		return *v
	}
	return random_equip()
	/*
		gprs := strconv.Itoa(int(equipid))
		_gprs2equipid[gprs] = equipid
		_id2equips[equipid] = random_equip()
	*/
}

func equip_drop_id(equipid uint32) Equip {
	if _, ok := _id2antennas[equipid]; ok {

	}
	return random_equip()
}
func equip_upsert(e Equip) (Equip, error) {
	orig := equip_get_id(e.EquipId)
	if e.EquipId != 0 && orig.EquipId == e.EquipId {
		return orig, nil
	}
	orig = equip_get_gprs(e.Gprs)
	if e.Gprs != "" && orig.Gprs == e.Gprs {
		return orig, nil
	}
	if e.EquipId == 0 {
		e.EquipId = equipid_new()
	}
	if e.Gprs == "" {
		e.Gprs = strconv.Itoa(int(e.EquipId))
	}
	_id2equips[e.EquipId] = &e
	post_index_equip(e)
	antenna_upsert(e.EquipId, e.Gprs)
	return e, nil
}
func equip_activate(e Equip) (Equip, error) {
	return random_equip(), nil
}
func post_index_equip(e Equip) {

}
func equip_get_gprs(gprs string) Equip {
	var id uint32
	var ok bool
	if id, ok = _gprs2equipid[gprs]; ok {
		return equip_get_id(id)
	}
	return random_equip()
	/*	id = new_equipid()
		_gprs2equipid[gprs] = id
		_id2equips[id] = random_equip()
		return _id2equips[id]
	*/
}

func equipid_new() uint32 {
	equip_id_seed = equip_id_seed + 1
	return equip_id_seed
}
func equips_get_geo(province, city string) []Equip {
	return nil
}
