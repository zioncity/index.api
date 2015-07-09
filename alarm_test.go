package main

import (
	"testing"
	"time"

	"code.google.com/p/go/src/pkg/strconv"
)

func TestAlarmGetEquip(t *testing.T) {
	t.Skip()
	x := alarms_get_equip(equip_id_test, 0, 10)
	t.Log(x)
}

func TestAlarmGetAntenna(t *testing.T) {
	t.Skip()
	x := alarms_get_antenna(equip_id_test, 0, 0, 10)
	t.Log(x, "what")
}

const equip_id_test = 18610071876

func TestAlarmUpdate(t *testing.T) {
	t.Skip()
	a := Alarm{equip_id_test, 0, strconv.Itoa(equip_id_test), 1, 2, 3, 4, time.Now().Unix()}
	//  a := alarm_init()
	x := alarm_update(a)
	t.Log(x)
}
