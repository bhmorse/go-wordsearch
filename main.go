package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bhmorse/go-wordsearch/wordsearch"
	"os"
)

type Configuration struct {
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Words  []string `json:"words"`
}

func main() {
	path := flag.String("config", "", "json config file")
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		fmt.Println("Error opening words file:", err)
		return
	}
	defer f.Close()

	config := Configuration{}
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file", err)
		return
	}

	tries := 0
	success := false
	var puzzle *wordsearch.Puzzle
	for !success {
		puzzle, err = wordsearch.NewPuzzle(config.Width, config.Height, config.Words)
		if err != nil && err != wordsearch.ErrMaxIterations {
			fmt.Println("Error placing words", err)
			return
		} else if err == nil {
			success = true
		}
		tries = tries + 1

		if tries > 100000 {
			fmt.Println("Cannot create a valid puzzle after 100000 tries")
			return
		}
	}

	puzzle.Print(5)
}
