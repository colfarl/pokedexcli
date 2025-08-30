package pokeapi

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Locations  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
