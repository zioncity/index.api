package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strconv"
	"time"
)

type handler func(w http.ResponseWriter, r *http.Request)

var addr = flag.String("addr", ":9202", "listen address like ip:port")

func main() {
	flag.Parse()
	equip_all()
	antennas_all()
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
	http.Handle("/antennas/batch", handler(handle_antennas_batch))
	http.Handle("/attitude/append", handler(handle_attitude_append))
	http.Handle("/attitudes/show", handler(handle_attitudes_show))
	http.Handle("/profile/edit", handler(handle_profile_edit))
	http.Handle("/profile/show", handler(handle_profile_show))
	http.Handle("/geo", handler(handle_baidu_geo)) // ?lat=xxx&lng=xxx&baidu=0
	http.ListenAndServe(*addr, nil)
}

//post.body = equip
func handle_equip_add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var val Equip
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&val))
	val = equip_add(val)

	panic_error(json.NewEncoder(w).Encode(&val))
}
func handle_equip_activate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var val Equip
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&val))
	val = equip_activate(val)

	panic_error(json.NewEncoder(w).Encode(&val))
}
func handle_equip_drop(w http.ResponseWriter, r *http.Request) {
	equipid := atoi(r.FormValue("equipid"))

	equip_drop_id(equipid)

}
func handle_equip_edit(w http.ResponseWriter, r *http.Request) {
	var v Equip
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&v))
	v = equip_add(v)

	panic_error(json.NewEncoder(w).Encode(&v))
}
func handle_equip_show(w http.ResponseWriter, r *http.Request) {
	equipid := atoi(r.FormValue("equipid"))

	v := equip_get_id(equipid)

	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_equip_attitude(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	equipid := atoi(r.FormValue("equipid"))
	from, count := atoi(r.FormValue("from")), atoi(r.FormValue("count"))
	if count == 0 {
		count = 10
	}
	v := attitudes_get_equip(equipid, int(from), int(count))
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

	for _, e := range v {
		equip_add(e)
	}

}
func handle_antennas_batch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var v []Antenna
	dec := json.NewDecoder(r.Body)
	panic_error(dec.Decode(&v))
	for _, a := range v {
		antenna_update(a)
	}
}
func handle_alarms_show(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	eqid := atoi(r.FormValue("equipid"))
	from, count := atoui32(r.FormValue("from")), atoui32(r.FormValue("count"))
	if count == 0 {
		count = 10
	}
	equip := equip_get_id(eqid)

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
	enable, equipid, unitid := atoui32(r.FormValue("enable")), atoi(r.FormValue("equipid")), atoui32(r.FormValue("unitid"))

	v := antenna_enable(equipid, unitid, int(enable))
	panic_error(json.NewEncoder(w).Encode(&v))
}

func handle_antennas_show(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	equipid := atoi(r.FormValue("equipid"))
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
	equipid, unitid := atoi(r.FormValue("equpid")), atoui32(r.FormValue("unitid"))
	from, count := atoui32(r.FormValue("from")), atoui32(r.FormValue("count"))
	if count == 0 {
		count = 10
	}
	v := attitudes_get(equipid, unitid, int(from), int(count))
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
