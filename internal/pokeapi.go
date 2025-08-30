// Package pokeapi serves as a wrapper around pokeapi
package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetLocations(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err	
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}	
	return body, nil
}

func ResponseToStruct(bytes []byte, dest *LocationArea)  error {
	err := json.Unmarshal(bytes, &dest);
	if err != nil {
		return err
	}
	return nil
}
