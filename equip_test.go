package main

import "testing"

func TestEquipDrop(t *testing.T) {
  equip_es_drop(equip_id_test)
}
func TestEquipsLoad(t *testing.T) {
  x := equip_es_load()
  t.Log(len(x))
}
func TestEquipEsUpsert(t *testing.T) {
  equip_es_upsert(equip_init())
}
