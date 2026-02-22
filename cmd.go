package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/PrincessFluffyButt937/Pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*CFG) error
}

type CFG struct {
	Next     string
	Previous string
	Arg      string
	cache    *pokecache.Cache
	Repo     map[string]Pokemon
}

type AreaList struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationEcounters struct {
	Encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name   string `json:"name"`
	Exp    int    `json:"base_experience"`
	Weight int    `json:"weight"`
	Height int    `json:"height"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

var cmd_req = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Lists next 20 locations",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Lists previous 20 locations",
		callback:    commandMapb,
	},
	"explore": {
		name:        "explore",
		description: "explore <location_name>",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "catch <pokemom_name>",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "inspect <pokemom_name>",
		callback:    commandInspect,
	},
}

func commandExit(con *CFG) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(con *CFG) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("explore <location_name>: Explores pokemon encouters in said area")
	fmt.Println("catch <pokemom_name>: Attempts to catch a pokemon.")
	fmt.Println("inspect <pokemom_name>: Inspects pokemon you have already caught.")
	fmt.Println("map: Lists next 20 locations")
	fmt.Println("mapb: Lists previous 20 locations")

	//fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	//for c := range cmd {
	//	fmt.Printf("%s: %s\n", cmd[c].name, cmd[c].description)
	return nil
}

func commandMap(con *CFG) error {
	base_url := con.Next
	cached, exists := con.cache.Get(base_url)
	//cache logic
	if exists {
		var cached_data AreaList
		if err := json.Unmarshal(cached, &cached_data); err != nil {
			fmt.Printf("Error - Unmarshaling cached response: %v", err)
			return err
		}
		con.Previous = con.Next
		con.Next = cached_data.Next

		for _, area := range cached_data.Results {
			fmt.Println(area.Name)
		}
		con.cache.Add(base_url, cached)
		//debug prints
		//fmt.Println("------------------")
		//fmt.Println("map cache")
		//fmt.Println("------------------")
		return nil
	}
	//Get request logic
	res, err := http.Get(base_url)
	if err != nil {
		fmt.Printf("Error - https GET request failed: %v", err)
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var data AreaList
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("Error - Unmarshaling response: %v", err)
		return err
	}
	con.Previous = con.Next
	con.Next = data.Next

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}
	con.cache.Add(base_url, body)
	return nil
}

func commandMapb(con *CFG) error {
	if con.Previous == "" {
		fmt.Println("There is no going back.")
		return nil
	}
	base_url := con.Previous
	cached, exists := con.cache.Get(base_url)
	//cache logic
	if exists {
		var cached_data AreaList
		if err := json.Unmarshal(cached, &cached_data); err != nil {
			fmt.Printf("Error - Unmarshaling cached response: %v", err)
			return err
		}
		con.Next = con.Previous
		con.Previous = cached_data.Previous

		for _, area := range cached_data.Results {
			fmt.Println(area.Name)
		}
		con.cache.Add(base_url, cached)
		//debug prints
		//fmt.Println("------------------")
		//fmt.Println("mapb cache")
		//fmt.Println("------------------")
		return nil
	}
	//Get request logic

	res, err := http.Get(base_url)
	if err != nil {
		fmt.Printf("Error - https GET request failed: %v", err)
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var data AreaList
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("Error - Unmarshaling response: %v", err)
		return err
	}
	con.Next = con.Previous
	con.Previous = data.Previous

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}
	con.cache.Add(base_url, body)
	return nil
}

func commandExplore(con *CFG) error {
	if con.Arg == "" {
		fmt.Println("No area selected")
		return nil
	}
	base_url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", con.Arg)

	cached, exists := con.cache.Get(base_url)
	//cache logic
	if exists {
		var cached_data LocationEcounters
		if err := json.Unmarshal(cached, &cached_data); err != nil {
			fmt.Printf("Error - Unmarshaling cached response: %v", err)
			return err
		}
		fmt.Printf("Exploring %s...\n", con.Arg)
		fmt.Println("Found Pokemon:")
		for _, encounter := range cached_data.Encounters {
			fmt.Printf(" - %s\n", encounter.Pokemon.Name)
		}
		con.cache.Add(base_url, cached)
		//debug prints
		//fmt.Println("------------------")
		//fmt.Println("explore cache")
		//fmt.Println("------------------")
		return nil
	}
	//Get request logic

	res, err := http.Get(base_url)
	if err != nil {
		fmt.Printf("Error - https GET request failed: %v", err)
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var data LocationEcounters

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("Error - Unmarshaling response: %v", err)
		return err
	}

	fmt.Printf("Exploring %s...\n", con.Arg)
	fmt.Println("Found Pokemon:")
	for _, encounter := range data.Encounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	con.cache.Add(base_url, body)
	return nil
}

func commandCatch(con *CFG) error {
	if con.Arg == "" {
		fmt.Println("No pokemon selected")
		return nil
	}
	base_url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", con.Arg)
	cached, exists := con.cache.Get(base_url)
	// cache logic
	if exists {
		var pokemon Pokemon
		if err := json.Unmarshal(cached, &pokemon); err != nil {
			fmt.Printf("Error - Unmarshaling response: %v", err)
			return err
		}
		con.cache.Add(base_url, cached)
		fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
		required := pokemon.Exp / 3
		if required < rand.IntN(pokemon.Exp) {
			con.Repo[pokemon.Name] = pokemon
			fmt.Printf("%s was caught!\n", pokemon.Name)
			return nil
		}
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil

	}
	res, err := http.Get(base_url)
	if err != nil {
		fmt.Printf("Error - https GET request failed: %v", err)
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		fmt.Printf("Error - Unmarshaling response: %v", err)
		return err
	}
	con.cache.Add(base_url, body)
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	required := pokemon.Exp / 3
	if required < rand.IntN(pokemon.Exp) {
		con.Repo[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
		return nil
	}
	fmt.Printf("%s escaped!\n", pokemon.Name)
	return nil
}

func commandInspect(con *CFG) error {
	if len(con.Repo) == 0 {
		fmt.Println("You have not caught any pokemons.")
	}
	if con.Arg == "" {
		fmt.Println("No pokemon name was selected.")
	} else {
		pokemon, exists := con.Repo[con.Arg]
		if !exists {
			fmt.Println("You havent caught this pokemon.")
			return nil
		}
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, pok_stat := range pokemon.Stats {
			fmt.Printf("  -%s: %v\n", pok_stat.Stat.Name, pok_stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, pok_type := range pokemon.Types {
			fmt.Printf("  - %s\n", pok_type.Type.Name)
		}
	}
	return nil
}
