package main

import (
	"bufio"
	"math/rand"
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/colfarl/pokedexcli/internal"
	pokecache "github.com/colfarl/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*settings, []string) error
}

type settings struct {
	nextLocationURL string
	prevLocationURL string
	cache		    *pokecache.Cache
	pokedex			map[string]pokeapi.Pokemon
}

var config = settings{
	nextLocationURL : "https://pokeapi.co/api/v2/location-area/",
	prevLocationURL : "",
	cache : pokecache.NewCache(5),
	pokedex: make(map[string]pokeapi.Pokemon),
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "explore the next 20 locations on the map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "explore the previous 20 locations on the map",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "lists pokemon in provided area: explore <area-name>",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "attempt to catch pokemon provided: catch <pokemon>",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "get a detailed look of a pokemon you have caught: inspect <pokemon>",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "see all pokemon you have caught: pokedex",
			callback:    commandPokedex,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))	
}

func commandExit(config *settings, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(config *settings, params []string) error {

	if config.nextLocationURL == "" {
		fmt.Printf("No more locations this way!\n")
		return fmt.Errorf("reached end of locations")
	}
	
	url := config.nextLocationURL
	val, exists := config.cache.Get(url)
	
	var locations pokeapi.LocationArea

	if exists {
		pokeapi.ResponseToLocationArea(val, &locations)
	} else {
		res, err := pokeapi.GetRequest(url)
		if err != nil {
			return err
		}
		pokeapi.ResponseToLocationArea(res, &locations)
		config.cache.Add(url, res)
	}

	config.nextLocationURL = locations.Next
	config.prevLocationURL = locations.Previous
	
	fmt.Println()
	for _, s := range locations.Locations {
		fmt.Println(s.Name)
	}
	fmt.Println()

	return nil
}

func commandMapB(config *settings, params []string) error {

	if config.prevLocationURL == "" {
		fmt.Printf("No more locations this way!\n")
		return fmt.Errorf("reached beginning of locations")
	}

	url := config.prevLocationURL
	val, exists := config.cache.Get(url)
	
	var locations pokeapi.LocationArea

	if exists {
		pokeapi.ResponseToLocationArea(val, &locations)
	} else {
		res, err := pokeapi.GetRequest(url)
		if err != nil {
			return err
		}
		pokeapi.ResponseToLocationArea(res, &locations)
		config.cache.Add(url, res)
	}

	config.nextLocationURL = locations.Next
	config.prevLocationURL = locations.Previous
	
	fmt.Println()
	for _, s := range locations.Locations {
		fmt.Println(s.Name)
	}
	fmt.Println()

	return nil
}

func commandHelp(config *settings, params []string) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandExplore(config *settings, params []string) error {
	if len(params) == 0 {
		return fmt.Errorf("usage: explore <area-name>")
	}
	
	area := params[0]
	url := "https://pokeapi.co/api/v2/location-area/" + area
	val, exists := config.cache.Get(url)		

	var locationInfo pokeapi.LocationInformation

	if exists {
		pokeapi.ResponseToLocationInformation(val, &locationInfo)
	} else {
		res, err := pokeapi.GetRequest(url)
		if err != nil {
			return err
		}
		pokeapi.ResponseToLocationInformation(res, &locationInfo)
		config.cache.Add(url, res)
	}	
	
	fmt.Println()
	for _, s := range locationInfo.PokemonEncounters {
		fmt.Println("-", s.Pokemon.Name)
	}
	fmt.Println()
	return nil
}

func commandCatch(config *settings, params []string) error {
	if len(params) == 0 {
		return fmt.Errorf("usage: catch <pokemon-name>")
	}
	
	name := params[0]
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	val, exists := config.cache.Get(url)
	var pokemon pokeapi.Pokemon

	if exists {
		pokeapi.ResponseToPokemon(val, &pokemon)
	} else {
		res, err := pokeapi.GetRequest(url)
		if err != nil {
			return err
		}
		pokeapi.ResponseToPokemon(res, &pokemon)
		config.cache.Add(url, res)
	}	
	
	fmt.Println()
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	catchRoll := rand.Intn(700)	
	
	if(catchRoll >= pokemon.BaseExperience){
		fmt.Printf("%s was caught!\n", name)
		config.pokedex[name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	fmt.Println()

	return nil
}

func commandInspect(config *settings, params []string) error {
	if len(params) == 0 {
		return fmt.Errorf("usage: inspect <pokemon-name>")
	}	
	
	name := params[0]
	pokemon, exists := config.pokedex[name]
	if !exists {
		fmt.Println("you have not caught that pokemon")
	} else {
		pokeapi.PokePrint(pokemon)
	}
	return nil
}

func commandPokedex(config *settings, params []string) error {

	if len(params) > 0 {
		return fmt.Errorf("usage: pokedex")
	}	
	
	if len(config.pokedex) == 0 {
		fmt.Println("you don't have any pokemon yet")
		return nil
	} 
	
	fmt.Println()
	for key, _ := range config.pokedex {
		fmt.Printf(" - %s\n", key)	
	}
	fmt.Println()
	return nil
}

func main() {	
	scanner := bufio.NewScanner(os.Stdin)
	
	for ; ; {
		fmt.Print("Pokedex > ")	
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		
		if len(cleanedInput) == 0 {
			continue
		}  

		command := cleanedInput[0]
		cmd, exists := getCommands()[command]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(&config, cleanedInput[1:]); err != nil {
			fmt.Println(err)
		}
	}
}
