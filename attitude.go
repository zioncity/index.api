package main

import (
	"math/rand"
	"time"
)

type Attitude struct {
	EquipId uint32 `json:"equipid"`
	Gprs    string `json:"gprs"`
	UnitId  uint32 `json:"unitid"`
	H       int    `json:"h"`      //cm centimetre
	X       int    `json:"x"`      //degrees * 100000
	Y       int    `json:"y"`      //degrees * 100000
	Z       int    `json:"z"`      //degrees * 100000
	Update  int64  `json:"update"` // unixtime
}

type EquipAntennaAttitude struct {
	Equip    Equip    `json:"equip"`
	Antenna  Antenna  `json:"antenna"`
	Attitude Attitude `json:"attitude"`
}

func attitudes_update(atts []Attitude) []EquipAntennaAttitude {
	v := make([]EquipAntennaAttitude, len(atts))
	for i := 0; i < len(atts); i++ {
		atts[i].Update = time.Now().Unix()
		if atts[i].Gprs == "" {
			atts[i].Gprs = random_no(16)
		}
		if atts[i].EquipId == 0 {
			atts[i].EquipId = rand.Uint32()
		}
		v[i].Antenna = antenna_init()
		v[i].Equip = random_equip()
		v[i].Attitude = atts[i]
	}
	return v
}
func random_attitude() Attitude {
	return Attitude{rand.Uint32(), random_no(16), rand.Uint32() % 18, rand.Intn(cm), rand.Intn(90 * cmd), rand.Intn(30 * cmd), rand.Intn(45 * cmd), time.Now().Unix()}
}
func random_attitudes(n int) []EquipAntennaAttitude {
	v := make([]EquipAntennaAttitude, n)
	for i := 0; i < n; i++ {
		v[i].Antenna = antenna_init()
		v[i].Equip = random_equip()
		v[i].Attitude = random_attitude()
	}
	return v
}

func attitudes_get(equipid, unitid uint32) []EquipAntennaAttitude {
	return nil
}

func attitudes_get_equip(equipid uint32) []EquipAntennaAttitude {
	return nil
}
