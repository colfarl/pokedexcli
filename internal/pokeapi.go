// Package pokeapi serves as a wrapper around pokeapi
package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetLocations(url string) (LocationArea, error) {
	var locations LocationArea	
	res, err := http.Get(url)
	if err != nil {
		return LocationArea{}, err	
	}
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, err
	}

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return LocationArea{}, err
	}	

	return locations, nil
}
