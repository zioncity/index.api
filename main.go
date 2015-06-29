package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	/*	http.HandleFunc("/equip/add", handle_notfound)
		http.HandleFunc("/equip/activate", handle_notfound)
		http.HandleFunc("/equip/drop", handle_notfound)
		http.HandleFunc("/equip/edit", handle_notfound)
		http.HandleFunc("/equips/show", handle_notfound)
		http.HandleFunc("/equips/batch", handle_notfound)
		http.HandleFunc("/alarm/show", handle_notfound)
		http.HandleFunc("/antenna/show", handle_notfound)
		http.HandleFunc("/antenna/enable", handle_notfound)*/
	//	http.HandleFunc("/", handle_notfound)
	http.Handle("/", handler(handle_notfound))
	http.ListenAndServe(":9202", nil)
}

func handle_notfound(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	if err := enc.Encode("hello handler"); err != nil {
		panic(err)
	}
}

type handler func(w http.ResponseWriter, r *http.Request)

func (imp handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Del("Content-Length")
			http.Error(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()
	imp(w, r)
	//	w.Header().Set("Content-Length", strconv.Itoa(n))
}
