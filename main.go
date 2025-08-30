package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))	
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for ; ; {
		fmt.Print("Pokedex > ")	
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		word := cleanedInput[0]
		fmt.Printf("Your command was: %v\n", word)
	}
}
