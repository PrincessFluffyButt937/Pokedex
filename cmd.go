package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*CFG) error
}

type CFG struct {
	Next     string
	Previous string
}

type AreaList struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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
	fmt.Println("map: Lists next 20 locations")
	fmt.Println("mapb: Lists previous 20 locations")
	//fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	//for c := range cmd {
	//	fmt.Printf("%s: %s\n", cmd[c].name, cmd[c].description)
	return nil
}

func commandMap(con *CFG) error {
	base_url := con.Next
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
	return nil
}

func commandMapb(con *CFG) error {
	if con.Previous == "" {
		fmt.Println("There is no going back.")
		return nil
	}
	base_url := con.Previous
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
	return nil
}
