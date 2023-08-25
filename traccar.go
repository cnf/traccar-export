package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"

	"golang.org/x/net/publicsuffix"
)

const TSFORMAT = "2006-01-02T15:04:05Z"

type Device struct {
	ID int `json:"id"`
	GroupID int `json:"groupId"`
	Name string `json:"name"`
	Model string `json:"model"`
	Category string `json:"Category"`
}
type Base struct {
	DeviceId   int       `json:"deviceId"`
	DeviceName string    `json:"deviceName"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
}
type Stop struct {
	Base
	// DeviceId      int       `json:"deviceId"`
	// DeviceName    string    `json:"deviceName"`
	// StartTime     time.Time `json:"startTime"`
	// EndTime       time.Time `json:"endTime"`
	Duration      int     `json:"duration"`
	Latitude      float64 `json:"latitude"`
	Longtitude    float64 `json:"longitude"`
	PositionId    int     `json:"positionId"`
	Address       string  `json:"address"`
	SpentFuel     float32 `json:"spentFuel"`
	EngineHours   int     `json:"engineHours"`
	Distance      float32 `json:"distance"`
	AverageSpeed  float32 `json:"averageSpeed"`
	MaxSpeed      float32 `json:"maxSpeed"`
	StartOdometer float32 `json:"startOdometer"`
	EndOdometer   float32 `json:"endOdometer"`
}

type Trip struct {
	Base
	Duration        int     `json:"duration"`
	Latitude        float32 `json:"latitude"`
	Longtitude      float32 `json:"longitude"`
	PositionId      int     `json:"positionId"`
	Address         string  `json:"address"`
	SpentFuel       float32 `json:"spentFuel"`
	EngineHours     int     `json:"engineHours"`
	Distance        float32 `json:"distance"`
	AverageSpeed    float32 `json:"averageSpeed"`
	MaxSpeed        float32 `json:"maxSpeed"`
	StartOdometer   float64 `json:"startOdometer"`
	EndOdometer     float64 `json:"endOdometer"`
	StartPositionId int     `json:"startPositionId"`
	EndPositionId   int     `json:"endPositionId"`
	StartLat        float32 `json:"startLat"`
	StartLon        float32 `json:"startLon"`
	EndLat          float32 `json:"endLat"`
	EndLon          float32 `json:"endLon"`
	StartAddress    string  `json:"startAddress"`
	EndAddress      string  `json:"endAddress"`
	DriverUniqueId  string  `json:"driverUniqueId"`
	DriverName      string  `json:"driverName"`
}

type Attr struct {
	Alarm string `json:"alarm"`
	//  "priority": 0,
	//         "sat": 0,
	//         "event": x,
	//         "ignition": false,
	//         "motion": true,
	//         "io200": 2,
	//         "io69": 3,
	//         "power": 12.623000000000001,
	//         "battery": 4.051,
	//         "io68": 0,
	//         "odometer": 123,
	//         "distance": 0.0,
	//         "totalDistance": 123.64,
	//         "hours": 123456,
	//         "alarm": "lowPower"

}
type Position struct {
	// Base
	ID         int     `json:"id"`
	DeviceId   int     `json:"deviceId"`
	Attributes Attr    `json:"attributes"`
	Protocol   string  `json:"protocol"`
	Latitude   float64 `json:"latitude"`
	Longtitude float64 `json:"longitude"`
	Altitude   float64 `json:"altitude"`

	FixTime  time.Time `json:"fixTime"`
	Outdated bool      `json:"outdated"`
	Valid    bool      `json:"valid"`

	// "serverTime": "2023-08-03T10:19:37.000+00:00",
	// "deviceTime": "2023-08-03T10:18:23.000+00:00",
	// "fixTime": "2023-08-03T10:18:23.000+00:00",

	// "speed": 0.0,
	// "course": 0.0,
	// "address": null,
	// "accuracy": 0.0,
	// "network": null,
	// "geofenceIds": null
}

type Traccar struct {
	url   string
	token string
	jar http.CookieJar
}

func NewTraccar(url string, token string) Traccar {
	var TC Traccar
	TC.url = url
	TC.token = token
	return TC
}

func (t *Traccar) get(path string, params map[string]string, target interface{}) error {
	if t.jar == nil {
		t.session()
	}
	return t.httpdo("GET", path, params, target)
}

func (t *Traccar) httpdo(method string, path string, params map[string]string, target interface{}) error {
	url := fmt.Sprintf("%s%s", t.url, path)
	// method := "GET"
	client := &http.Client{Jar: t.jar}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// if res.StatusCode >=200 and res.StatusCode < 400 {
	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		fmt.Printf("Body: %s\n", res.Body)
		return fmt.Errorf("API returned %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return err
	}
	return nil
}

func (t *Traccar) session() error {
	fmt.Println("Getting session token")
	var err error
	params := map[string]string{"token": t.token}
	t.jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	return t.httpdo("GET", "/session", params, nil)
}

func (t *Traccar) GetDeviceByID(id int) (*Device, error) {
	params := map[string]string{
		"id": strconv.Itoa(id),
	}
	
	var dev []Device

	err := t.get("/devices", params, &dev)
	if err != nil {
		return nil, err
	}
	return &dev[0], nil
}

func (t *Traccar) GetStops(id int, from time.Time, to time.Time) ([]Stop, error) {
	params := map[string]string{
		"deviceId": strconv.Itoa(id),
		"from":     from.Format(TSFORMAT),
		"to":       to.Format(TSFORMAT),
	}

	var stops []Stop

	err := t.get("/reports/stops", params, &stops)
	if err != nil {
		return nil, err
	}
	return stops, nil
}

func (t *Traccar) GetTrips(id int, from time.Time, to time.Time) ([]Trip, error) {
	params := map[string]string{
		"deviceId": strconv.Itoa(id),
		"from":     from.Format(TSFORMAT),
		"to":       to.Format(TSFORMAT),
	}

	var trips []Trip

	err := t.get("/reports/trips", params, &trips)
	if err != nil {
		return nil, err
	}
	return trips, nil
}

func (t *Traccar) GetRoute(id int, from time.Time, to time.Time) ([]Position, error) {
	params := map[string]string{
		"deviceId": strconv.Itoa(id),
		"from":     from.Format(TSFORMAT),
		"to":       to.Format(TSFORMAT),
	}

	var route []Position

	err := t.get("/reports/route", params, &route)
	if err != nil {
		return nil, err
	}
	return route, nil
}
