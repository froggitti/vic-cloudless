package vtr

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"
)

// free no-sign-up weather for everywhere on earth, theoretically
// a user who will never go above, like, 5 requests per second

type WeatherCondition string

const (
	Cloudy        WeatherCondition = "Cloudy"
	Windy         WeatherCondition = "Windy"
	Thunderstorms WeatherCondition = "Thunderstorms"
	Sunny         WeatherCondition = "Sunny"
	Clear         WeatherCondition = "Clear"
	Snowy         WeatherCondition = "Snowy"
	Cold          WeatherCondition = "Cold"
)

func getWeather(location string) (tempC, tempF string, weather WeatherCondition, err error) {
	geoURL := "https://nominatim.openstreetmap.org/search?" + url.Values{
		"format": {"json"},
		"q":      {location},
	}.Encode()
	req1, _ := http.NewRequest("GET", geoURL, nil)
	res1, err := http.DefaultClient.Do(req1)
	if err != nil {
		return "", "", Cold, err
	}
	defer res1.Body.Close()

	var geo []struct{ Lat, Lon string }
	if err = json.NewDecoder(res1.Body).Decode(&geo); err != nil {
		return "", "", Cold, err
	}
	if len(geo) == 0 {
		return "", "", Cold, fmt.Errorf("no geocode for %q", location)
	}
	lat, lon := geo[0].Lat, geo[0].Lon
	oURL := "https://api.open-meteo.com/v1/forecast?" + url.Values{
		"latitude":        {lat},
		"longitude":       {lon},
		"current_weather": {"true"},
		"daily":           {"sunrise,sunset"},
		"timezone":        {"auto"},
	}.Encode()
	oRes, err := http.Get(oURL)
	if err != nil {
		return "", "", Cold, err
	}
	defer oRes.Body.Close()

	var om struct {
		UTCOffsetSeconds int `json:"utc_offset_seconds"`
		CurrentWeather   struct {
			Temperature float64 `json:"temperature"`
			WeatherCode int     `json:"weathercode"`
			Time        string  `json:"time"`
		} `json:"current_weather"`
		Daily struct {
			Sunrise []string `json:"sunrise"`
			Sunset  []string `json:"sunset"`
		} `json:"daily"`
	}
	if err = json.NewDecoder(oRes.Body).Decode(&om); err != nil {
		return "", "", Cold, err
	}

	c := om.CurrentWeather.Temperature
	f := c*9.0/5.0 + 32.0
	tempC = fmt.Sprint(math.Round(c))
	tempF = fmt.Sprint(math.Round(f))

	loc := time.FixedZone("local", om.UTCOffsetSeconds)
	layout := "2006-01-02T15:04"
	ct, err := time.ParseInLocation(layout, om.CurrentWeather.Time, loc)
	if err != nil {
		return tempC, tempF, Cold, fmt.Errorf("time parse current: %w", err)
	}
	sunrise, err := time.ParseInLocation(layout, om.Daily.Sunrise[0], loc)
	if err != nil {
		return tempC, tempF, Cold, fmt.Errorf("time parse sunrise: %w", err)
	}
	sunset, err := time.ParseInLocation(layout, om.Daily.Sunset[0], loc)
	if err != nil {
		return tempC, tempF, Cold, fmt.Errorf("time parse sunset: %w", err)
	}

	code := om.CurrentWeather.WeatherCode
	switch {
	case code >= 95 && code <= 99:
		weather = Thunderstorms
	case code == 71 || code == 73 || code == 75 || code == 77 || code == 85 || code == 86:
		weather = Snowy
	case code == 1 || code == 2:
		weather = Sunny
	case code == 0:
		if ct.After(sunrise) && ct.Before(sunset) {
			weather = Sunny
		} else {
			weather = Clear
		}
	case code == 3 || code == 45 || code == 48 || (code >= 61 && code <= 82):
		weather = Cloudy
	default:
		weather = Cloudy
	}
	if c <= 0 {
		weather = Cold
	}

	return tempC, tempF, weather, nil
}
