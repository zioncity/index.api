package main

import "math/rand"

type Equip struct {
	EquipId uint32 `json:"equipid"`
	Gprs    string `json:"gprs"`
	Name    string `json:"name,omitempty"`
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
