// Package pokeapi serves as a wrapper around pokeapi
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetRequest(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err	
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Unsuccessful request") 
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}	
	return body, nil
}

func ResponseToLocationArea(bytes []byte, dest *LocationArea)  error {
	err := json.Unmarshal(bytes, &dest);
	if err != nil {
		return err
	}
	return nil
}

func ResponseToLocationInformation(bytes []byte, dest *LocationInformation)  error {
	err := json.Unmarshal(bytes, &dest);
	if err != nil {
		return err
	}
	return nil
}

func ResponseToPokemon(bytes []byte, dest *Pokemon)  error {
	err := json.Unmarshal(bytes, &dest);
	if err != nil {
		return err
	}
	return nil
}

func PokePrint(pokemon Pokemon) {
	fmt.Println()
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)	
	fmt.Println()
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats{
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
		
	fmt.Println()
	fmt.Println("Type(s):")
	for _, typ := range pokemon.Types{
		fmt.Printf(" - %s\n", typ.Type.Name) 
	}

	fmt.Println()
	
}

