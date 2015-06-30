package main

type Profile struct {
	Duration int `json:"duration"` //hour  >0 and < 30*24
	H        int `json:"h"`        //cm  >-100 and < 100
	X        int `json:"x"`        //degree/1000 >-2000 and <2000
	Y        int `json:"y"`        //degree/1000 >-5000 and < 5000
	Z        int `json:"z"`        //degreee/1000 > -5000 and < 5000
}

var _profile Profile = Profile{24, 1 * cm, 2 * cmd, 5 * cmd, 5 * cmd}

func validate_profile(p Profile) Profile {
	if p.Duration <= 0 || p.Duration > 30*24 {
		p.Duration = _profile.Duration
	}
	if p.H < -100 || p.H > 100 || p.H == 0 {
		p.H = _profile.H
	}
	if p.X < -360*cmd || p.X > 360*cmd || p.X == 0 {
		p.X = _profile.X
	}
	if p.Y < -360*cmd || p.Y > 360*cmd || p.Y == 0 {
		p.Y = _profile.Y
	}
	if p.Z < -360*cmd || p.Z > 360*cmd || p.Z == 0 {
		p.Z = _profile.Z
	}
	return p
}