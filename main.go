package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	EmptyArr         = make([]rune, 1, 1)
	Empty            = EmptyArr[0]
	ErrMaxIterations = errors.New("Exceeded 1000 iterations")
	NumColsToPrint   = 5
)

type LetterPosition struct {
	Row int
	Col int
}

type Configuration struct {
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Words  []string `json:"words"`
}

func PrintPuzzle(board [][]rune, words []string) {
	for _, row := range board {
		for _, letter := range row {
			fmt.Print(string(letter), " ")
		}
		fmt.Println("")
	}
	fmt.Println("")
	sort.Strings(words)
	div := len(words) / NumColsToPrint
	rem := len(words) % NumColsToPrint
	numRows := div
	if rem != 0 {
		numRows = numRows + 1
	}
	for row := 0; row < numRows; row++ {
		remLeft := rem
		idx := row
		for col := 0; col < NumColsToPrint; col++ {
			if row == (numRows-1) && remLeft == 0 {
				break
			}
			fmt.Print(strings.ToUpper(words[idx]))
			if len(words[idx]) >= 8 {
				fmt.Print("\t")
			} else {
				fmt.Print("\t\t")
			}
			idx = idx + div
			if remLeft > 0 {
				idx = idx + 1
				remLeft = remLeft - 1
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func PlaceWords(config Configuration) (puzzle [][]rune, err error) {
	puzzle = make([][]rune, config.Height, config.Height)
	positions := map[rune][]LetterPosition{}
	for i := 0; i < int(config.Height); i++ {
		puzzle[i] = make([]rune, config.Width, config.Width)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, word := range config.Words {
		// We need to pick a direction, find valid range, and roll a spot.
		// TODO: Allow words to cross?

		word = strings.ToUpper(word)
		success := false
		iterations := 0
		for !success {
			direction := r.Intn(12)
			colDir := 0
			rowDir := 0
			row := 0
			col := 0
			switch direction {
			case 0: //horizontal forwards
				row = r.Intn(config.Height)
				col = r.Intn(config.Width - len(word) + 1)
				colDir = 1
			case 1: //horizontal backwards
				row = r.Intn(config.Height)
				col = r.Intn(config.Width-len(word)+1) + len(word) - 1
				colDir = -1
			case 2: //vertical forwards
				row = r.Intn(config.Height-len(word)+1) + len(word) - 1
				col = r.Intn(config.Width)
				rowDir = -1
			case 3: //vertical backwards
				row = r.Intn(config.Height - len(word) + 1)
				col = r.Intn(config.Width)
				rowDir = 1
			case 4:
				fallthrough
			case 5: //forward slash forwards
				row = r.Intn(config.Height-len(word)+1) + len(word) - 1
				col = r.Intn(config.Width - len(word) + 1)
				rowDir = -1
				colDir = 1
			case 6:
				fallthrough
			case 7: //forward slash backwards
				row = r.Intn(config.Height - len(word) + 1)
				col = r.Intn(config.Width-len(word)+1) + len(word) - 1
				rowDir = 1
				colDir = -1
			case 8:
				fallthrough
			case 9: //backslash forwards
				row = r.Intn(config.Height - len(word) + 1)
				col = r.Intn(config.Width - len(word) + 1)
				rowDir = 1
				colDir = 1
			case 10:
				fallthrough
			case 11: //backslash backwards
				row = r.Intn(config.Height-len(word)+1) + len(word) - 1
				col = r.Intn(config.Width-len(word)+1) + len(word) - 1
				rowDir = -1
				colDir = -1
			default:
				fmt.Println("Something strange occurred")
				return puzzle, errors.New("Rand broke")
			}
			clear := true
			for i, letter := range word {
				check := puzzle[row+rowDir*i][col+colDir*i]
				if check != Empty && check != letter {
					clear = false
					break
				}
			}

			if clear {
				for i, letter := range word {
					puzzle[row+rowDir*i][col+colDir*i] = letter
					positions[letter] = append(positions[letter], LetterPosition{
						Row: row + rowDir*i,
						Col: col + colDir*i,
					})
				}
				success = true
			}

			iterations = iterations + 1
			if iterations > 1000 {
				return puzzle, ErrMaxIterations
			}
		}

	}
	return puzzle, nil
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
	puzzle := [][]rune{}
	for !success {
		puzzle, err = PlaceWords(config)
		if err != nil && err != ErrMaxIterations {
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

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < config.Height; i++ {
		for j := 0; j < config.Width; j++ {
			if puzzle[i][j] == Empty {
				value := r.Intn(26) + 65
				r, _ := utf8.DecodeRune([]byte{byte(value)})
				puzzle[i][j] = r
			}
		}
	}

	PrintPuzzle(puzzle, config.Words)
}
