package main

import "time"

// import

// https://github.com/tkrajina/gpxgo/blob/master/test_files/gpx1.0_with_all_fields.gpx
type GPX struct {
	Name        string `xml:"name"`
	Description string `xml:"desc"`
	Bounds      GPXBounds
	Waypoints   []GPXwpt
	Tracks      []GPXtrk
}

type GPXBounds struct{}

type GPXtrk struct {
	Name     string `xml:"name"`
	Number   int    `xml:"number"`
	Segments []GPXtrkseg
}

type GPXtrkseg struct {
	Points []GPXtrkpt
}

type GPXtrkpt struct {
	lat  float32
	lon  float32
	ele  float32
	time time.Time
}

type GPXwpt struct {
	lat  float32
	lon  float32
	ele  float32
	time time.Time
	name string
	cmt  string
	desc string
	sym  string
}

type GPXextensions struct{}
