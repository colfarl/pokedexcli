package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	pokeapi "github.com/colfarl/pokedexcli/internal"
	pokecache "github.com/colfarl/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*settings) error
}

type settings struct {
	nextLocationURL string
	prevLocationURL string
	cache		    *pokecache.Cache
}

var config = settings{
	nextLocationURL : "https://pokeapi.co/api/v2/location-area/",
	prevLocationURL : "",
	cache : pokecache.NewCache(5),
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
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))	
}

func commandExit(config *settings) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(config *settings) error {

	if config.nextLocationURL == "" {
		fmt.Printf("No more locations this way!\n")
		return fmt.Errorf("reached end of locations")
	}
	
	url := config.nextLocationURL
	val, exists := config.cache.Get(url)
	
	var locations pokeapi.LocationArea

	if exists {
		pokeapi.ResponseToStruct(val, &locations)
	} else {
		res, err := pokeapi.GetLocations(url)
		if err != nil {
			return err
		}
		pokeapi.ResponseToStruct(res, &locations)
		config.cache.Add(url, res)
	}

	config.nextLocationURL = locations.Next
	config.prevLocationURL = locations.Previous
	
	for _, s := range locations.Locations {
		fmt.Println(s.Name)
	}

	return nil
}

func commandMapB(config *settings) error {

	if config.prevLocationURL == "" {
		fmt.Printf("No more locations this way!\n")
		return fmt.Errorf("reached beginning of locations")
	}

	url := config.prevLocationURL
	val, exists := config.cache.Get(url)
	
	var locations pokeapi.LocationArea

	if exists {
		pokeapi.ResponseToStruct(val, &locations)
	} else {
		res, err := pokeapi.GetLocations(url)
		if err != nil {
			return err
		}
		pokeapi.ResponseToStruct(res, &locations)
		config.cache.Add(url, res)
	}

	config.nextLocationURL = locations.Next
	config.prevLocationURL = locations.Previous
	
	for _, s := range locations.Locations {
		fmt.Println(s.Name)
	}

	return nil
}

func commandHelp(config *settings) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
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

		if err := cmd.callback(&config); err != nil {
			fmt.Println(err)
		}
	}
}
