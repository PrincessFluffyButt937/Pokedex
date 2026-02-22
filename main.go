package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PrincessFluffyButt937/Pokedex/internal/pokecache"
)

func main() {
	Start_REPL()
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func Start_REPL() {
	scanner := bufio.NewScanner(os.Stdin)
	con := CFG{
		Next:  "https://pokeapi.co/api/v2/location-area/",
		cache: pokecache.NewCache(5 * time.Second),
		Repo:  make(map[string]Pokemon),
	}
	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := cleanInput(scanner.Text())
		if len(text) == 0 {
			continue
		}
		if len(text) > 1 {
			con.Arg = text[1]
		} else {
			con.Arg = ""
		}
		commmand := text[0]
		cmd, exists := cmd_req[commmand]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		cmd.callback(&con)
		//fmt.Printf("Your command was: %v\n", text[0])

	}

}
