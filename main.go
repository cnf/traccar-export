package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/tkrajina/gpxgo/gpx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)


func main() {
	ParseConfig()
	
	from := viper.GetTime("from")
	to := viper.GetTime("to")
	TC := NewTraccar(viper.GetString("url"), viper.GetString("token"))

	dev, err := TC.GetDeviceByID(viper.GetInt("device"))
	if err != nil {
		log.Fatal(err)
	}
	id := dev.ID

	trips, err := TC.GetTrips(id, from, to)
	if err != nil {
		log.Fatal(err)
	}

	mygpx := new(gpx.GPX)
	mygpx.AddMissingTime()
	if dev.Name != "" {
		mygpx.Description = fmt.Sprintf("Trip from %s to %s with %s", from.Format(time.DateOnly), to.Format(time.DateOnly), dev.Name)
	}
	if viper.IsSet("gpx.title") {
		mygpx.Name = viper.GetString("gpx.title")
	} else {
		mygpx.Name = "Traccar Export"
	}
	
	if viper.IsSet("gpx.author") {
		mygpx.AuthorName = viper.GetString("gpx.author")
	}
	mygpx.Creator = "Traccar Exporter"
	t := time.Now()
	mygpx.Time = &t
	// mygpx.Bounds().

	var current time.Time
	tracks := make(map[string]*gpx.GPXTrack)
	for _, trip := range trips {
		fmt.Println(trip.StartTime.Format(time.DateOnly))
		route, err := TC.GetRoute(trip.DeviceId, trip.StartTime, trip.EndTime)
		if err != nil {
			fmt.Println(err)
			return
		}
		// for

		var segment = gpx.GPXTrackSegment{}
		for _, point := range route {
			segment.AppendPoint(&gpx.GPXPoint{
				Timestamp: point.FixTime,
				Point: gpx.Point{
					Latitude:  point.Latitude,
					Longitude: point.Longtitude,
					Elevation: *gpx.NewNullableFloat64(point.Altitude),
				},
			})
		}

		if current.YearDay() != trip.StartTime.YearDay() {
			current = trip.StartTime
			c := cases.Title(language.Und)
			var DevType string
			if dev.Category != "" {
				DevType = c.String(dev.Category)
			} else {
				DevType = "Unknown"
			}
			tracks[trip.StartTime.Format(time.DateOnly)] = &gpx.GPXTrack{
				Name: trip.StartTime.Format(time.DateOnly), 
				Type: DevType,
				Source: "Traccar",
			}

		}
		tracks[trip.StartTime.Format(time.DateOnly)].AppendSegment(&segment)
	}
	for _, t := range tracks {
		mygpx.AppendTrack(t)
	}

	stops, err := TC.GetStops(id, from, to)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, stop := range stops {
		var p = gpx.GPXPoint{
			Timestamp: stop.StartTime,
			Point: gpx.Point{
				Latitude:  stop.Latitude,
				Longitude: stop.Longtitude,
				// Elevation: *gpx.NewNullableFloat64(stop.Altitude),
			},
			Symbol: "Information",
		}
		mygpx.AppendWaypoint(&p)

	}

	mygpx.Bounds()

	xmlBytes, err := mygpx.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("test.gpx", xmlBytes, 0666); err != nil {
		log.Fatal(err)
	}
}
