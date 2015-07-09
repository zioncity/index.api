package main

import "testing"

func TestBaiduLocation(t *testing.T) {
  url := baidu_location("39.983424051248", "116.32298703399", 1)
  t.Log(url)
}

func TestBaiduLatLng(t *testing.T) {
  lat, lng := baidu_latlong("29.579084284809", "114.23075168433")
  t.Log(lat, lng)
}
