package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

const (
	ak   = "YTB7UKOXCsL9mlqUIT8tKRTR"
	sk   = "fbBFzGF0A1SBtGAvuE1UlaMpvqyLd0R7"
	base = "http://api.map.baidu.com"
)

//      "lng": 116.32298703399,
//      "lat": 39.983424051248

func baidu_location(lat, lng string, isbaidu int) string {
	params := url.Values{}
	if isbaidu == 0 {
		blat, blng := baidu_latlong(lat, lng)
		params.Add("location", strconv.FormatFloat(blat, 'f', 11, 64)+","+strconv.FormatFloat(blng, 'f', 11, 64))
	} else {
		params.Add("location", lat+","+lng)
	}
	params.Add("output", "json")
	uri := baidu_sn("/geocoder/v2/", params)
	return uri
}

func baidu_sn(path string, params url.Values) string {
	params.Add("ak", ak)
	x := params.Encode()
	x = url.QueryEscape(path + "?" + x + sk)
	sn := md5.Sum([]byte(x))
	params.Add("sn", hex.EncodeToString(sn[:]))
	return base + path + "?" + params.Encode()
}

//{"status":0,"result":[{"x":114.23075168433,"y":29.579084284809}]}
func baidu_latlong(gpslat, gpslng string) (float64, float64) {
	params := url.Values{}
	params.Add("coords", gpslng+","+gpslat)
	params.Add("output", "json")

	uri := baidu_sn("/geoconv/v1/", params)
	resp, err := http.Get(uri)
	panic_error(err)

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic_error(errors.New("status not ok"))
	}
	var v struct {
		Status int `json:"status"`
		Result []struct {
			Long float64 `json:"x"`
			Lat  float64 `json:"y"`
		} `json:"result"`
	}
	err = json.NewDecoder(resp.Body).Decode(&v)
	panic_error(err)
	if v.Status != 0 || len(v.Result) != 1 {
		panic_error(errors.New("invalid baidu status"))
	}
	return v.Result[0].Lat, v.Result[0].Long
	//return strconv.FormatFloat(v.Result[0].Lat, 'f', 11, 64), strconv.FormatFloat(v.Result[0].Long, 'f', 11, 64), nil
}
