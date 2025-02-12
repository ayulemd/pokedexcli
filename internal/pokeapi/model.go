package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Result `json:"results"`
}

type Config struct {
	Next     string
	Previous string
}

func GetLocationArea(url string) (LocationArea, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationArea{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return LocationArea{}, err
	}

	locationArea := LocationArea{}
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return LocationArea{}, nil
	}

	return locationArea, nil
}

func DisplayLocationAreas(areas LocationArea, config *Config) {
	config.Next = areas.Next
	config.Previous = areas.Previous

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
}
