package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type handler func(w http.ResponseWriter, r *http.Request)

func main() {
	http.Handle("/equip/add", handler(handle_equip_add))
	http.Handle("/equip/activate", handler(handle_equip_activate))
	http.Handle("/equip/drop", handler(handle_equip_drop))
	http.Handle("/equip/edit", handler(handle_equip_edit))
	http.Handle("/equip/show", handler(handle_equip_show))
	http.Handle("/equip/attitude", handler(handle_equip_attitude))
	http.Handle("/equips/show", handler(handle_equips_show))
	http.Handle("/equips/batch", handler(handle_equips_batch))
	http.Handle("/alarms/show", handler(handle_alarms_show))
	http.Handle("/alarm/add", handler(handle_alarm_add))
	http.Handle("/antennas/show", handler(handle_antennas_show))
	http.Handle("/antenna/enable", handler(handle_antenna_enable))
	http.Handle("/attitude/append", handler(handle_attitude_append))
	http.Handle("/attitudes/show", handler(handle_attitudes_show))
	http.Handle("/profile/edit", handler(handle_profile_edit))
	http.Handle("/profile/show", handler(handle_profile_show))
	http.Handle("/geo", handler(handle_baidu_geo)) // ?lat=xxx&lng=xxx&baidu=0
	http.ListenAndServe(":9202", nil)
}

//post.body = equip
func handle_equip_add(w http.ResponseWriter, r *http.Request) {
	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_equip_activate(w http.ResponseWriter, r *http.Request) {
	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_equip_drop(w http.ResponseWriter, r *http.Request) {
	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_equip_edit(w http.ResponseWriter, r *http.Request) {
	var v Equip
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&v))
	if v.EquipId == 0 {
		v.EquipId = rand.Uint32()
	}
	if v.Gprs == "" {
		v.Gprs = random_no(16)
	}
	if v.Name == "" {
		v.Name = v.Gprs
	}
	//	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_equip_show(w http.ResponseWriter, r *http.Request) {
	equipid, _ := strconv.Atoi(r.FormValue("equipid"))
	gprs := r.FormValue("gprs")
	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	if equipid != 0 {
		v.EquipId = uint32(equipid)
	}
	if gprs != "" {
		v.Gprs = gprs
	}
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_equip_attitude(w http.ResponseWriter, r *http.Request) {
	v := random_attitudes(18)
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_equips_show(w http.ResponseWriter, r *http.Request) {
	v := random_equips(32)
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_equips_batch(w http.ResponseWriter, r *http.Request) {
	v := random_equips(24)
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_alarms_show(w http.ResponseWriter, r *http.Request) {
	v := random_alarms(32)
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_alarm_add(w http.ResponseWriter, r *http.Request) {
	var alarm Alarm
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&alarm))
	alarm.Date = time.Now().Unix()
	var v EquipAntennaAlarm
	v.Equip = random_equip()
	v.Antenna = random_antenna()
	v.Alarm = alarm
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_antenna_enable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	disable, _ := strconv.Atoi(r.FormValue("disable"))
	equipid, _ := strconv.Atoi(r.FormValue("equipid"))
	unitid, _ := strconv.Atoi(r.FormValue("unitid"))
	gprs := r.FormValue("gprs")
	v := EquipAntenna{random_equip(), random_antenna()}
	v.Antenna.EquipId = uint32(equipid)
	v.Equip.EquipId = uint32(equipid)
	v.Antenna.UnitId = uint32(unitid)
	v.Antenna.Disable = disable
	if gprs != "" {
		v.Equip.Gprs = gprs
		v.Antenna.Gprs = gprs
	}
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_antennas_show(w http.ResponseWriter, r *http.Request) {
	v := random_antennas(32)
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_attitude_append(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	atts := []Attitude{} //EquipId/Gprs HXYZ
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&atts))
	v := validate_attitudes(atts)
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_attitudes_show(w http.ResponseWriter, r *http.Request) {
	v := random_attitudes(32)
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_profile_edit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var p Profile
	panic_error(dec.Decode(&p))
	_profile = validate_profile(p)
	panic_error(json.NewEncoder(w).Encode(&_profile))
}

func handle_profile_show(w http.ResponseWriter, r *http.Request) {
	panic_error(json.NewEncoder(w).Encode(&_profile))
}

func handle_notfound(w http.ResponseWriter, r *http.Request) {
	panic_error(json.NewEncoder(w).Encode("hello fixed"))
}

func handle_baidu_geo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	lat, lng := r.FormValue("lat"), r.FormValue("lng")
	is_baidu, _ := strconv.Atoi(r.FormValue("baidu"))
	uri := baidu_location(lat, lng, is_baidu)
	w.Header().Del("Content-Type")
	w.Header().Set("Location", uri)
	w.WriteHeader(http.StatusFound)
}
func (imp handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()
	imp(w, r)
}
