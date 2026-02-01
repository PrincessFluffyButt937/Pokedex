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
	//fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	//for c := range cmd {
	//	fmt.Printf("%s: %s\n", cmd[c].name, cmd[c].description)
	return nil
}

func commandMap(con *CFG) error {
	base_url := "https://pokeapi.co/api/v2/location-area"
	res, err := http.Get(base_url)
	if err != nil {
		fmt.Printf("Error - https GET request failed: %v", err)
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("Error - Unmarshaling response: %v", err)
		return err
	}
	fmt.Println(string(body))

	return nil
}
