package main

import "math/rand"

type Equip struct {
	EquipId   uint32 `json:"equipid"`
	Gprs      string `json:"gprs"`
	Name      string `json:"name,omitempty"`
	Province  string `json:"province"`
	City      string
	Lat       uint32 // 129.11 * cmd //cmd=1000
	Lng       uint32 // 36.87 *cmd
	Activated uint16 // 0 : unactivated, 1: activated
	Online    uint16 // 0:offline, 1:online
	Update    int64  // unixtime the last update time
}

func random_equip() Equip {
	return Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
}

func random_equips(n int) []Equip {
	v := make([]Equip, n)
	for i := 0; i < n; i++ {
		v[i] = Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	}
	return v
}
