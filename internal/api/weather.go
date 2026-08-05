package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	openMeteoForecastURL  = "https://api.open-meteo.com/v1/forecast"
	openMeteoGeocodingURL = "https://geocoding-api.open-meteo.com/v1/search"
	weatherCacheDuration  = 10 * time.Minute
)

type weatherCache struct {
	mu        sync.Mutex
	key       string
	expiresAt time.Time
	value     currentWeatherResponse
}

type weatherLocation struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Country   string  `json:"country,omitempty"`
	Admin1    string  `json:"admin1,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone,omitempty"`
}

type currentWeatherResponse struct {
	Location            string  `json:"location"`
	Country             string  `json:"country,omitempty"`
	Temperature         float64 `json:"temperature"`
	ApparentTemperature float64 `json:"apparentTemperature"`
	Precipitation       float64 `json:"precipitation"`
	WeatherCode         int     `json:"weatherCode"`
	WindSpeed           float64 `json:"windSpeed"`
	WindDirection       int     `json:"windDirection"`
	WindGusts           float64 `json:"windGusts"`
	UpdatedAt           string  `json:"updatedAt"`
	Provider            string  `json:"provider"`
}

func (s *Server) searchWeatherLocations(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := strings.TrimSpace(
		r.URL.Query().Get("q"),
	)

	if len([]rune(query)) < 2 {
		writeError(
			w,
			http.StatusBadRequest,
			"enter at least two characters",
		)
		return
	}

	if len([]rune(query)) > 100 {
		writeError(
			w,
			http.StatusBadRequest,
			"location search is too long",
		)
		return
	}

	locations, err := s.openMeteoLocations(
		r.Context(),
		query,
		8,
	)
	if err != nil {
		s.writeWeatherUpstreamError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, locations)
}

func (s *Server) currentWeather(
	w http.ResponseWriter,
	r *http.Request,
) {
	weather := s.store.Get().Modules.Weather

	if !weather.Enabled {
		writeError(
			w,
			http.StatusNotFound,
			"weather widget is disabled",
		)
		return
	}

	if weather.Latitude == nil ||
		weather.Longitude == nil {
		writeError(
			w,
			http.StatusUnprocessableEntity,
			"weather location has not been selected",
		)
		return
	}

	cacheKey := fmt.Sprintf(
		"%s|%.6f|%.6f|%s|%s",
		weather.Provider,
		*weather.Latitude,
		*weather.Longitude,
		weather.Location,
		weather.Country,
	)

	if cached, ok := s.cachedWeather(cacheKey); ok {
		writeJSON(w, http.StatusOK, cached)
		return
	}

	current, err := s.fetchOpenMeteoWeather(
		r.Context(),
		*weather.Latitude,
		*weather.Longitude,
		weather.Location,
		weather.Country,
	)
	if err != nil {
		s.writeWeatherUpstreamError(w, err)
		return
	}

	s.storeCachedWeather(cacheKey, current)
	writeJSON(w, http.StatusOK, current)
}

func (s *Server) openMeteoLocations(
	ctx context.Context,
	query string,
	count int,
) ([]weatherLocation, error) {
	values := url.Values{}
	values.Set("name", query)
	values.Set("count", strconv.Itoa(count))
	values.Set("language", "en")
	values.Set("format", "json")

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		openMeteoGeocodingURL+"?"+values.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set(
		"User-Agent",
		"Hublet/1.0 (+https://github.com/pzdera/hublet)",
	)

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"geocoding service returned HTTP %d",
			response.StatusCode,
		)
	}

	var payload struct {
		Results []weatherLocation `json:"results"`
	}

	if err := decodeWeatherJSON(
		response.Body,
		&payload,
	); err != nil {
		return nil, err
	}

	if payload.Results == nil {
		payload.Results = []weatherLocation{}
	}

	return payload.Results, nil
}

func (s *Server) fetchOpenMeteoWeather(
	ctx context.Context,
	latitude float64,
	longitude float64,
	location string,
	country string,
) (currentWeatherResponse, error) {
	values := url.Values{}
	values.Set("latitude", strconv.FormatFloat(latitude, 'f', 6, 64))
	values.Set("longitude", strconv.FormatFloat(longitude, 'f', 6, 64))
	values.Set(
		"current",
		"temperature_2m,apparent_temperature,precipitation,weather_code,wind_speed_10m,wind_direction_10m,wind_gusts_10m",
	)
	values.Set("temperature_unit", "celsius")
	values.Set("wind_speed_unit", "kmh")
	values.Set("precipitation_unit", "mm")
	values.Set("timezone", "auto")

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		openMeteoForecastURL+"?"+values.Encode(),
		nil,
	)
	if err != nil {
		return currentWeatherResponse{}, err
	}

	request.Header.Set(
		"User-Agent",
		"Hublet/1.0 (+https://github.com/pzdera/hublet)",
	)

	response, err := s.httpClient.Do(request)
	if err != nil {
		return currentWeatherResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return currentWeatherResponse{}, fmt.Errorf(
			"weather service returned HTTP %d",
			response.StatusCode,
		)
	}

	var payload struct {
		Current struct {
			Time                string  `json:"time"`
			Temperature         float64 `json:"temperature_2m"`
			ApparentTemperature float64 `json:"apparent_temperature"`
			Precipitation       float64 `json:"precipitation"`
			WeatherCode         int     `json:"weather_code"`
			WindSpeed           float64 `json:"wind_speed_10m"`
			WindDirection       int     `json:"wind_direction_10m"`
			WindGusts           float64 `json:"wind_gusts_10m"`
		} `json:"current"`
	}

	if err := decodeWeatherJSON(
		response.Body,
		&payload,
	); err != nil {
		return currentWeatherResponse{}, err
	}

	return currentWeatherResponse{
		Location:            location,
		Country:             country,
		Temperature:         payload.Current.Temperature,
		ApparentTemperature: payload.Current.ApparentTemperature,
		Precipitation:       payload.Current.Precipitation,
		WeatherCode:         payload.Current.WeatherCode,
		WindSpeed:           payload.Current.WindSpeed,
		WindDirection:       payload.Current.WindDirection,
		WindGusts:           payload.Current.WindGusts,
		UpdatedAt:           payload.Current.Time,
		Provider:            "open-meteo",
	}, nil
}

func (s *Server) cachedWeather(
	key string,
) (currentWeatherResponse, bool) {
	s.weatherCache.mu.Lock()
	defer s.weatherCache.mu.Unlock()

	if s.weatherCache.key != key ||
		time.Now().After(s.weatherCache.expiresAt) {
		return currentWeatherResponse{}, false
	}

	return s.weatherCache.value, true
}

func (s *Server) storeCachedWeather(
	key string,
	value currentWeatherResponse,
) {
	s.weatherCache.mu.Lock()
	defer s.weatherCache.mu.Unlock()

	s.weatherCache.key = key
	s.weatherCache.value = value
	s.weatherCache.expiresAt = time.Now().Add(
		weatherCacheDuration,
	)
}

func (s *Server) writeWeatherUpstreamError(
	w http.ResponseWriter,
	err error,
) {
	status := http.StatusBadGateway

	if errors.Is(err, context.DeadlineExceeded) {
		status = http.StatusGatewayTimeout
	}

	writeError(
		w,
		status,
		"weather service is temporarily unavailable",
	)
}

func decodeWeatherJSON(
	reader io.Reader,
	target any,
) error {
	return json.NewDecoder(
		io.LimitReader(reader, 1<<20),
	).Decode(target)
}
