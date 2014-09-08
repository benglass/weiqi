// http://howistart.org/posts/go/1

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	params, e := loadParametersFile("./parameters.json")
	if e != nil {
		log.Printf("Error loading parameters file: %s\n", e)
		os.Exit(1)
	}
	log.Println(params.WundergroundApiKey)

	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		provider := openWeatherMap{}

		data, err := provider.temperature(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})
	http.ListenAndServe(":8080", nil)
}

func loadParametersFile(path string) (parameters, error) {
	file, e := os.Open(path)
	if e != nil {
		params := parameters{}
		return params, e
	}
	return loadParameters(file)
}


func loadParameters(r io.Reader) (parameters, error) {
	params := parameters{}
	if err := json.NewDecoder(r).Decode(&params); err != nil {
		return params, err
	}
	return params, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

type parameters struct {
	WundergroundApiKey string `json:"wunderground_api_key"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

type weatherProvider interface {
	temperature(city string) (float64, error)
}

type openWeatherMap struct{}

func (w openWeatherMap) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)

	return d.Main.Kelvin, nil
}
