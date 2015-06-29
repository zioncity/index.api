package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.Handle("/equip/add", handler(handle_equip_add))
	http.Handle("/equip/activate", handler(handle_equip_activate))
	http.Handle("/equip/drop", handler(handle_equip_drop))
	http.Handle("/equip/edit", handler(handle_equip_edit))
	http.Handle("/equips/show", handler(handle_equips_show))
	http.Handle("/equips/batch", handler(handle_equips_batch))
	http.Handle("/alarms/show", handler(handle_alarms_show))
	http.Handle("/alarm/add", handler(handle_alarm_add))
	http.Handle("/antennas/show", handler(handle_antennas_show))
	http.Handle("/antenna/enable", handler(handle_antenna_enable))
	http.Handle("/profile/edit", handler(handle_profile_edit))
	http.Handle("/profile/show", handler(handle_profile_show))

	http.ListenAndServe(":9202", nil)
}

type Equip struct {
	EquipId uint32 `json:"equipid"`
	Gprs    string `json:"gprs"`
	Name    string `json:"name,omitempty"`
}

const fake_gprs = "08618610041897"

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}
func gprs2id(gprs []byte) (v uint32) {
	if len(gprs) >= 4 {
		binary.Read(bytes.NewReader(gprs[len(gprs)-4:]), binary.BigEndian, &v)
	}
	return v
}

var letters = []rune("1234567890")

func random_no(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//post.body = equip
func handle_equip_add(w http.ResponseWriter, r *http.Request) {
	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	check_error(json.NewEncoder(w).Encode(&v))
}
func handle_equip_activate(w http.ResponseWriter, r *http.Request) {
	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	check_error(json.NewEncoder(w).Encode(&v))
}
func handle_equip_drop(w http.ResponseWriter, r *http.Request) {
	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	check_error(json.NewEncoder(w).Encode(&v))
}
func handle_equip_edit(w http.ResponseWriter, r *http.Request) {
	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	check_error(json.NewEncoder(w).Encode(&v))
}
func random_equips(n int) []Equip {
	v := make([]Equip, n)
	for i := 0; i < n; i++ {
		v[i] = Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	}
	return v
}
func handle_equips_show(w http.ResponseWriter, r *http.Request) {
	//	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	v := random_equips(32)
	check_error(json.NewEncoder(w).Encode(&v))
}
func handle_equips_batch(w http.ResponseWriter, r *http.Request) {
	v := random_equips(24)
	check_error(json.NewEncoder(w).Encode(&v))
}

type Alarm struct {
	EquipId uint32 `json:"equipid"`
	Gprs    string `json:"gprs"`
	UnitId  uint32 `json:"unitid"`
	Typ     string `json:"typ"`  //
	Date    int64  `json:"date"` //unixtime
}

type Attitude struct {
	EquipId uint32 `json:"equipid"`
	Gprs    string `json:"gprs"`
	UnitId  uint32 `json:"unitid"`
	H       int    `json:"h"`    //cm centimetre
	X       int    `json:"x"`    //degrees * 100000
	Y       int    `json:"y"`    //degrees * 100000
	Z       int    `json:"z"`    //degrees * 100000
	Date    int64  `json:"date"` // unixtime
}

type Antenna struct {
	EquipId uint32
	Gprs    string
	UnitId  uint32
	Lng     int // 129.23 *100000
	Lat     int // 32.11 * 100000
}

func random_alarm_type() string {
	var alarm_type = []string{"下顷", "滚降", "方位", "高度"}
	i := rand.Intn(len(alarm_type))
	return alarm_type(i)
}
func random_date() int64 {
	return time.Now().Add(-rand.Intn(24) * time.Hour).Unix()
}
func random_alarms(n int) []Alarm {
	v := make([]Alarm, n)
	for i := 0; i < n; i++ {
		v[i] = Alarm{EquipId: rand.Uint32(), Gprs: random_no(16),
			Uintid: rand.Uint32(), Typ: random_alarm_type(), Date: random_date()}
	}
	return v
}

func handle_alarms_show(w http.ResponseWriter, r *http.Request) {
	v := random_alarms(32)
	check_error(json.NewEncoder(w).Encode(&v))
}
func handle_alarm_add(w http.ResponseWriter, r *http.Request) {
	v := random_alarms(2)
	check_error(json.NewEncoder(w).Encode(&v))
}
func handle_antenna_enable(w http.ResponseWriter, r *http.Request) {
	v := Equip{EquipId: rand.Uint32(), Gprs: random_no(16)}
	check_error(json.NewEncoder(w).Encode(&v))
}
func random_antennas(n int) []Antenna {
	v := make([]Antenna, n)
	for i := 0; i < n; i++ {
		v[i] = Antenna{EquipId: rand.Uint32(), Gprs: random_no(16), UnitId: uint32(rand.Intn(18))}
	}
	return v
}
func handle_antennas_show(w http.ResponseWriter, r *http.Request) {
	v := random_antennas(32)
	check_error(json.NewEncoder(w).Encode(&v))
}

type Profile struct {
	Duration time.Duration `json:"duration"`
	H        int           `json:"h"`
	X        int           `json:"x"`
	Y        int           `json:"y"`
	Z        int           `json:"z"`
}

func handle_profile_edit(w http.ResponseWriter, r *http.Request) {
	v := Profile{time.Hour * 12, 1, 2, 5, 5}
	check_error(json.NewEncoder(w).Encode(&v))
}

func handle_profile_show(w http.ResponseWriter, r *http.Request) {
	v := Profile{time.Hour * 12, 1, 2, 5, 5}
	check_error(json.NewEncoder(w).Encode("hello fixed"))
}

func handle_notfound(w http.ResponseWriter, r *http.Request) {
	check_error(json.NewEncoder(w).Encode("hello fixed"))
}

type handler func(w http.ResponseWriter, r *http.Request)

func (imp handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()
	imp(w, r)
}
