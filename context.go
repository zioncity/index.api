package main

import (
	"math/rand"
	"strconv"
	"time"
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

func random_date() int64 {
	return time.Now().Add(-time.Hour * time.Duration(rand.Intn(24))).Unix()
}

var _profile Profile = Profile{24, 1 * cm, 2 * cmd, 5 * cmd, 5 * cmd}

var _id2equips map[uint32]Equip
var _id2antennas map[uint32][]Antenna
var _gprs2equipid map[string]uint32
var _id2gprs map[uint32]string

var equip_id_seed uint32 = 10000000

func get_equip_byid(equipid uint32) Equip {
	if v, ok := _id2equips[equipid]; ok {
		return v
	}
	return random_equip()
	/*
		gprs := strconv.Itoa(int(equipid))
		_gprs2equipid[gprs] = equipid
		_id2equips[equipid] = random_equip()
	*/
}
func upsert_equip(e Equip) Equip {
	orig := get_equip_byid(e.EquipId)
	if e.EquipId != 0 && orig.EquipId == e.EquipId {
		return orig
	}
	orig = get_equip_bygprs(e.Gprs)
	if e.Gprs != "" && orig.Gprs == e.Gprs {
		return orig
	}
	if e.EquipId == 0 {
		e.EquipId = new_equipid()
	}
	if e.Gprs == "" {
		e.Gprs = strconv.Itoa(int(e.EquipId))
	}
	_id2equips[e.EquipId] = e
	post_index_equip(e)
	upsert_antenna(e.EquipId, e.Gprs)
	return e
}
func post_index_equip(e Equip) {

}
func post_index_antennas([]Antenna) {

}
func upsert_antenna(equipid uint32, gprs string) {
	ats := make([]Antenna, unit_count)
	initialize_antennas(ats, gprs, equipid)
	post_index_antennas(ats)
	_id2antennas[equipid] = ats
}
func get_equip_bygprs(gprs string) Equip {
	var id uint32
	var ok bool
	if id, ok = _gprs2equipid[gprs]; ok {
		return get_equip_byid(id)
	}
	return random_equip()
	/*	id = new_equipid()
		_gprs2equipid[gprs] = id
		_id2equips[id] = random_equip()
		return _id2equips[id]
	*/
}
func get_antenna_byid(equipid uint32, unitid int) Antenna {
	if unitid >= unit_count {
		return random_antenna()
	}
	var val []Antenna
	var ok bool
	if val, ok = _id2antennas[equipid]; ok {
		return val[unitid]
	}
	return random_antenna()
	/*	gprs := strconv.Itoa(int(equipid))
		_gprs2equipid[gprs] = equipid
		val = make([]Antenna, unit_count)
		initialize_antennas(val, gprs, equipid)
		_id2antennas[equipid] = val
		return val[unitid]
	*/
}
func initialize_antennas(ats []Antenna, gprs string, equipid uint32) {
	for i := 0; i < len(ats); i++ {
		ats[i].UnitId = uint32(i)
		ats[i].EquipId = equipid
		ats[i].Gprs = gprs
	}
}
func get_antenna_bygprs(gprs string, unitid int) Antenna {
	var id uint32
	var ok bool
	if id, ok = _gprs2equipid[gprs]; ok {
		return get_antenna_byid(id, unitid)
	}
	return random_antenna()
	/*	id = new_equipid()
		_gprs2equipid[gprs] = id
		ats := make([]Antenna, unit_count)
		initialize_antennas(ats, gprs, id)
		_id2antennas[id] = ats
		return ats[unitid]
	*/
}

func new_equipid() uint32 {
	equip_id_seed = equip_id_seed + 1
	return equip_id_seed
}
