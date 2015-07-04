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
	r.ParseForm()
	var val Equip
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&val))
	val, err := equip_upsert(val)
	panic_error(err)
	panic_error(json.NewEncoder(w).Encode(&val))
}
func handle_equip_activate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var val Equip
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&val))
	val, err := equip_activate(val)
	panic_error(err)
	panic_error(json.NewEncoder(w).Encode(&val))
}
func handle_equip_drop(w http.ResponseWriter, r *http.Request) {
	equipid, _ := strconv.Atoi(r.FormValue("equipid"))
	gprs := r.FormValue("gprs")
	if equipid != 0 {
		equip_drop_id(uint32(equipid))
	} else if gprs != "" {
		equip_drop_id(equip_get_gprs(gprs).EquipId)
	}
	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_equip_edit(w http.ResponseWriter, r *http.Request) {
	var v Equip
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&v))
	v, err := equip_upsert(v)
	panic_error(err)
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_equip_show(w http.ResponseWriter, r *http.Request) {
	equipid := atoui32(r.FormValue("equipid"))
	gprs := r.FormValue("gprs")
	v := equip_get_id(equipid)
	if gprs != "" {
		v = equip_get_gprs(gprs)
	}

	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_equip_attitude(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	equipid := atoui32(r.FormValue("equipid"))
	v := attitudes_get_equip(equipid)
	//	v := random_attitudes(18)
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_equips_show(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	prov, city := r.FormValue("province"), r.FormValue("city")
	v := equips_get_geo(prov, city)
	//v := random_equips(32)
	panic_error(json.NewEncoder(w).Encode(&v))
}

type errorx []error

func (self errorx) Error() string {
	var v string
	for _, e := range self {
		v = v + e.Error()
	}
	return v
}
func handle_equips_batch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var v []Equip
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&v))
	errs := errorx{}
	for _, e := range v {
		if _, err := equip_upsert(e); err != nil {
			errs = append(errs, err)
		}
	}
	panic_error(json.NewEncoder(w).Encode(errs.Error()))
}

func handle_alarms_show(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	eqid, gprs := atoui32(r.FormValue("equipid")), r.FormValue("gprs")
	from, count := atoui32(r.FormValue("from")), atoui32(r.FormValue("count"))
	if count == 0 {
		count = 10
	}
	equip := equip_get_id(eqid)
	if gprs != "" {
		equip = equip_get_gprs(gprs)
	}
	v := alarms_get_equip(equip.EquipId, int(from), int(count))
	//	v := random_alarms(32)
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_alarm_add(w http.ResponseWriter, r *http.Request) {
	var alarm Alarm
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&alarm))
	alarm.Date = time.Now().Unix()
	v := alarm_update(alarm)
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_antenna_enable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	disable, equipid, unitid := atoui32(r.FormValue("disable")), atoui32(r.FormValue("equipid")), atoui32(r.FormValue("unitid"))
	//	gprs := r.FormValue("gprs")
	//	equip := equip_get_id(equipid)
	//	antenna := antenna_get_id(equipid, unitid)
	v := antenna_disable(equipid, unitid, disable)
	//	antenna.Disable = disable
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_antennas_show(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	equipid := atoui32(r.FormValue("equipid"))
	v := antennas_get_equip(equipid)
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_attitude_append(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	atts := []Attitude{} //EquipId/Gprs HXYZ
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&atts))
	v := attitudes_update(atts)
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_attitudes_show(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	equipid, unitid := atoui32(r.FormValue("equpid")), atoui32(r.FormValue("unitid"))
	v := attitudes_get(equipid, unitid)
	//	v := random_attitudes(32)
	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_profile_edit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var p Profile
	panic_error(dec.Decode(&p))
	_profile = profile_update(p)
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
