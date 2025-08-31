# pokedexcli (Go)

A tiny terminal Pokedex that talks to the public **PokeAPI**.  
Browse location areas, list encounterable Pokémon, try to catch them, and inspect your Pokédex — all from a REPL-style prompt.

## Quick start

    git clone https://github.com/colfarl/pokedexcli
    cd pokedexcli
    go run .
    # or: go build && ./pokedexcli

> Requires Go 1.20+ (recent Go versions work fine).

## Commands

- `help` — show available commands  
- `map` — list the next 20 location areas  
- `mapb` — list the previous 20 location areas  
- `explore <area-name>` — list Pokémon that can appear in the given location area  
- `catch <pokemon>` — attempt to catch a Pokémon by name  
- `inspect <pokemon>` — show details for a caught Pokémon  
- `pokedex` — list all Pokémon you’ve caught  
- `exit` — quit

### Examples

    Pokedex > map
    Pokedex > explore canalave-city-area
     - starly
     - bidoof
    Pokedex > catch starly
    Throwing a Pokeball at starly...
    starly was caught!
    Pokedex > pokedex
     - starly
    Pokedex > inspect starly
    # (prints stats/height/weight/etc.)

## How it works

- **CLI loop**: `main.go` reads a line, tokenizes it, then dispatches to a command function with `(*settings, []string)`.
- **HTTP wrapper**: `internal/pokeapi` performs `GET` requests and unmarshals JSON into typed structs.
- **Response cache**: `internal/pokecache` stores raw JSON by URL with a TTL (default **5s**) and a background ticker-based reaper to purge stale entries.
- **Catch logic**: a simple RNG check versus the Pokémon’s `BaseExperience` (fun, not canon).

## Notes

- Area names for `explore` are PokeAPI slugs (discover via `map`/`mapb`).
- If you hit rate limits or slow networks, increase the cache TTL in `settings` (e.g., `pokecache.NewCache(15)`).

## TODO (nice-to-haves)

Update the CLI to support the "up" arrow to cycle through previous commands
Simulate battles between pokemon
Add more unit tests
Refactor your code to organize it better and make it more testable
Keep pokemon in a "party" and allow them to level up
Allow for pokemon that are caught to evolve after a set amount of time
Persist a user's Pokedex to disk so they can save progress between sessions
Use the PokeAPI to make exploration more interesting. For example, rather than typing the names of areas, maybe you are given choices of areas and just type "left" or "right"
Random encounters with wild pokemon
Adding support for different types of balls (Pokeballs, Great Balls, Ultra Balls, etc), which have different chances of catching pokemon

