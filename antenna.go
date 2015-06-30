package main

import "math/rand"

type Antenna struct {
	EquipId uint32
	Gprs    string
	UnitId  uint32
	Lng     int // 129.23 *100000
	Lat     int // 32.11 * 100000
	Disable int
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
