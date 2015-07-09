package main

import "testing"

func TestAntennaUpsert(t *testing.T) {
	t.Skip()
	antenna_upsert(equip_id_test, id2gprs(equip_id_test))
}
func TestAntennaEnable(t *testing.T) {
	t.Skip()
	if x := antenna_enable(equip_id_test, 0, 1); x != 0 {
		t.Fail()
	}
	if x := antenna_enable(equip_id_test, 0, 0); x != 0 {
		t.Fail()
	}
}
func TestAntennaUpdate(t *testing.T) {
	t.Skip()
	x := antenna_init()
	antenna_update(x)
	antenna_update_bias(equip_id_test, 0, 1, 2, 3, 4)
}

func TestAntennaEsLoad(t *testing.T) {
	t.Skip()
	v := antennas_es_load()
	t.Log(len(v))
}
