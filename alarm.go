package main

import "math/rand"

type Alarm struct {
	EquipId uint32 `json:"equipid"`
	Gprs    string `json:"gprs"`
	UnitId  uint32 `json:"unitid"`
	Typ     string `json:"typ"`  //
	Date    int64  `json:"date"` //unixtime
}

type EquipAntennaAlarm struct {
	Equip   Equip   `json:"equip"`
	Antenna Antenna `json:"antenna"`
	Alarm   Alarm   `json:"alarm"`
}

var alarm_type = []string{"下顷", "滚降", "方位", "高度"}

func random_alarm_type() string {

	i := rand.Intn(len(alarm_type))
	return alarm_type[i]
}

func random_alarms(n int) []EquipAntennaAlarm {
	v := make([]EquipAntennaAlarm, n)
	for i := 0; i < n; i++ {
		v[i].Alarm = random_alarm()
		v[i].Equip = random_equip()
		v[i].Antenna = random_antenna()
	}
	return v
}

func random_alarm() Alarm {
	return Alarm{EquipId: rand.Uint32(), Gprs: random_no(16),
		UnitId: rand.Uint32(), Typ: random_alarm_type(), Date: random_date()}
}
