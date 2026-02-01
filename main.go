package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	//Start_REPL()
	var temp CFG // delete this!!!!
	commandMap(&temp)
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func Start_REPL() {
	scanner := bufio.NewScanner(os.Stdin)
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
		var temp CFG // delete this!!!!
		cmd.callback(&temp)
		//fmt.Printf("Your command was: %v\n", text[0])

	}

}
