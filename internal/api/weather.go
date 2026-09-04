package api

import (
	"fmt"
	"sync"
	"encoding/json"
	"math"
)


type WeatherData struct {
	CurrentData  CurrentWeather
	ForecastData Forecast
}

type CurrentWeather struct {
	Station struct {
		Name       string   `json:"name"`
		Latitude   float64  `json:"latitude"`
		Longitude  float64  `json:"longitude"`
	} `json:"station"`
	Measures struct {
		Temperature struct{
			Celcius string `json:"value"`
		} `json:"outside_temperature"`
		Humidity struct {
			Percent string     `json:"value"`
		} `json:"outside_humidity"`
		Wind struct {
			Speed string   `json:"value"`
		} `json:"wind_speed"`
	} `json:"observation"`
}

type Forecast struct {
	Day [3]DailyForecast
}

type DailyForecast struct {
	Data struct {
		Tmin int `json:"tmin"`
		Tmax int `json:"tmax"`
	} `json:"forecast"`
}


func getForecastPerDay(day int, cityCode string, client *Client) (*DailyForecast, error) {
	query := map[string]string{ "insee": cityCode }

	resp, err := client.APICall("forecast/daily/"+fmt.Sprint(day), query)
	if err != nil {
		return nil, fmt.Errorf("error while retreiving daily forecast data: %w", err)
	}

	var dailyForecast DailyForecast

	err = json.Unmarshal(resp, &dailyForecast)
	if err != nil {
		return nil, fmt.Errorf("error while reading JSON data: %w", err)
	}

	return &dailyForecast, nil
}

func getForecast(cityCode string, client *Client) (*Forecast, error) {
	var (
		wg2        sync.WaitGroup
		mut        sync.Mutex
		forecast   Forecast
		firstErr   error
		pointerErr error
	)

	for day := 0; day < 3; day++ {
		wg2.Add(1)
		go func(day int) {
			defer wg2.Done()
			dailyForecast, err := getForecastPerDay(day, cityCode, client)
			if err != nil {
				mut.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("error when retreiving forecast data: %w", err)
				}
				mut.Unlock()
				return
			}
			if dailyForecast == nil {
				mut.Lock()
				if pointerErr == nil {
					pointerErr = fmt.Errorf("null pointer received: %w", err)
				}
				mut.Unlock()
			}

			forecast.Day[day] = *dailyForecast
		}(day)
	}

	wg2.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	if pointerErr != nil {
		return nil, pointerErr
	}

	return &forecast, nil
}

func chooseBestStation(city *City, locations []CurrentWeather) (int, error) {
	var closestStation int = -1
	var minimumDistance float64 = math.Inf(1)
	var distance float64

	for i, location := range locations {
		if location.Measures.Temperature.Celcius == "" || location.Measures.Humidity.Percent == "" || location.Measures.Wind.Speed == "" {
			continue
		}
		distance = math.Pow(math.Abs(city.Latitude - location.Station.Latitude), 2) + math.Pow(math.Abs(city.Longitude - location.Station.Longitude), 2)
		distance = math.Sqrt(distance)
		if distance < minimumDistance {
			minimumDistance = distance
			closestStation = i
		}
	}

	if closestStation == -1 {
		return 0, fmt.Errorf("no station has enough data close to this point")
	}

	return closestStation, nil
}

func getCurrentWeather(city *City, client *Client) (*CurrentWeather, error) {
	query := map[string]string{
		"insee": city.Insee,
		"radius": "30",
	}

	resp, err := client.APICall("observations/around", query)
	if err != nil {
		return nil, fmt.Errorf("error while retreiving current weather data: %w", err)
	}

	var locations []CurrentWeather
	json.Unmarshal(resp, &locations)
	if len(locations) == 0 {
		return nil, fmt.Errorf("no weather station found")
	}

	closestStationIndex, err := chooseBestStation(city, locations)
	if err != nil {
		return nil, fmt.Errorf("no information available: %w", err)
	}

	return &locations[closestStationIndex], nil
}

func GetWeatherData(city *City, client *Client) (*WeatherData, error) {
	var (
		wg          sync.WaitGroup
		mu          sync.Mutex
		fullData    WeatherData
		errForecast error
		errCurrent  error
		errPointer  error
	)
	
	wg.Go(func() {
		var f *Forecast
		f, errForecast = getForecast(city.Insee, client)
		if f == nil {
			mu.Lock()
			errPointer = fmt.Errorf("null pointer returned")
			mu.Unlock()
		}
		if errForecast == nil {
			fullData.ForecastData = *f
		}
	})

	wg.Go(func() {
		var cw *CurrentWeather
		cw, errForecast = getCurrentWeather(city, client)
		if cw == nil {
			mu.Lock()
			errPointer = fmt.Errorf("null pointer returned")
			mu.Unlock()
		}
		if errForecast == nil {
			fullData.CurrentData = *cw
		}
	})

	wg.Wait()

	if errForecast != nil {
		return nil, fmt.Errorf("weather forecast cancelled: %w", errForecast)
	}
	if errCurrent != nil {
		return nil, fmt.Errorf("current weather retrieval cancelled: %w", errCurrent)
	}
	if errPointer != nil {
		return nil, fmt.Errorf("null pointer does not support dereferencing: %w", errPointer)
	}
	

	return &fullData, nil
}