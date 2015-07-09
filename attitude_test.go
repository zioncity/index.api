package main

import "testing"

func attitude_init(id int64) Attitude {
	return Attitude{EquipId: id, UnitId: 0}
}
func TestAttitudeEsInsert(t *testing.T) {
	t.Skip()
	atts := []Attitude{attitude_init(equip_id_test)}
	attitudes_es_insert(atts)
}

func TestAttitudesGet(t *testing.T) {
	t.Skip()
	x := attitudes_get(equip_id_test, 0, 0, 10)
	t.Log(len(x))
}

func TestAttitudesGetEquip(t *testing.T) {
	t.Skip()
	x := attitudes_get_equip(equip_id_test, 0, 10)
	t.Log(len(x))
}
