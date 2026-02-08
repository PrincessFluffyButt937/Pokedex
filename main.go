package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Print(time.Now())
	//Start_REPL()
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func Start_REPL() {
	scanner := bufio.NewScanner(os.Stdin)
	var con CFG
	con.Next = "https://pokeapi.co/api/v2/location-area/"
	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := cleanInput(scanner.Text())
		if len(text) == 0 {
			continue
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
