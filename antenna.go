package main

import "math/rand"

type Antenna struct {
	EquipId     uint32 `json:"equipid"`
	Gprs        string `json:"gprs"`
	UnitId      uint32 `json:"unitid"`  //[0,18)
	Lng         int    `json:"lng"`     // 129.23 *100000
	Lat         int    `json:"lat"`     // 32.11 * 100000
	Disable     int    `json:"disable"` //1:disable, 0:not disable
	Network     int    `json:"network"` //0:reserved 1:移动， 2:联通, 3:电信
	BaseId      uint32 `json:"baseid"`
	Basename    string `json:"basename"`
	CellId      uint16 `json:"cellid"`
	BaseManu    string `json:"basemanu"`
	AntennaManu string `json:"antennamanu"`
	AntennaTyp  string `json:"antennatyp"`
	AlertStart  int64  `json:"alterstart"` //unixtime, 0 for no alert
	Alert       int
	H           uint32 `json:"h"` //cm
	X           uint32 `json:"x"` //cmd
	Y           uint32 `json:"y"` //cmd
	Z           uint32 `json:"z"`
	Update      int64  `json:"update"` //unixtime the last update
}
type EquipAntenna struct {
	Equip   Equip   `json:"equip"`
	Antenna Antenna `json:"antenna"`
}

func random_antenna() Antenna {
	return Antenna{EquipId: rand.Uint32(), Gprs: random_no(16), UnitId: uint32(rand.Intn(18))}
}

func random_antennas(n int) []EquipAntenna {
	v := make([]EquipAntenna, n)
	for i := 0; i < n; i++ {
		v[i].Antenna = random_antenna()
		v[i].Equip = random_equip()
	}
	return v
}
