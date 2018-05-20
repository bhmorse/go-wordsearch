package wordsearch

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	EmptyArr         = make([]rune, 1, 1)
	Empty            = EmptyArr[0]
	ErrMaxIterations = errors.New("Exceeded 1000 iterations")
	ErrWordTooLong   = errors.New("Word length exceeds width and height")
)

type LetterPosition struct {
	Row int
	Col int
}

type Puzzle struct {
	Height int
	Width  int
	Board  [][]rune
	Words  []string
}

func NewPuzzle(width, height int, words []string) (puzzle *Puzzle, err error) {
	puzzle = &Puzzle{
		Height: height,
		Width:  width,
		Words:  words,
		Board:  make([][]rune, height, height),
	}
	err = puzzle.InitializeAndFill()
	return puzzle, err

}

func (p *Puzzle) Print(numCols int) {
	for _, row := range p.Board {
		for _, letter := range row {
			fmt.Print(string(letter), " ")
		}
		fmt.Println("")
	}
	fmt.Println("")
	sort.Strings(p.Words)
	div := len(p.Words) / numCols
	rem := len(p.Words) % numCols
	numRows := div
	if rem != 0 {
		numRows = numRows + 1
	}
	for row := 0; row < numRows; row++ {
		remLeft := rem
		idx := row
		for col := 0; col < numCols; col++ {
			if row == (numRows-1) && remLeft == 0 {
				break
			}
			fmt.Print(strings.ToUpper(p.Words[idx]))
			if len(p.Words[idx]) >= 8 {
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

func (p *Puzzle) InitializeAndFill() (err error) {
	positions := map[rune][]LetterPosition{}
	for i := 0; i < int(p.Height); i++ {
		p.Board[i] = make([]rune, p.Width, p.Width)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, word := range p.Words {
		// We need to pick a direction, find valid range, and roll a spot.
		// TODO: Allow words to cross?

		if len(word) > p.Width && len(word) > p.Height {
			return ErrWordTooLong
		}

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
				row = r.Intn(p.Height)
				col = r.Intn(p.Width - len(word) + 1)
				colDir = 1
			case 1: //horizontal backwards
				row = r.Intn(p.Height)
				col = r.Intn(p.Width-len(word)+1) + len(word) - 1
				colDir = -1
			case 2: //vertical forwards
				row = r.Intn(p.Height-len(word)+1) + len(word) - 1
				col = r.Intn(p.Width)
				rowDir = -1
			case 3: //vertical backwards
				row = r.Intn(p.Height - len(word) + 1)
				col = r.Intn(p.Width)
				rowDir = 1
			case 4:
				fallthrough
			case 5: //forward slash forwards
				row = r.Intn(p.Height-len(word)+1) + len(word) - 1
				col = r.Intn(p.Width - len(word) + 1)
				rowDir = -1
				colDir = 1
			case 6:
				fallthrough
			case 7: //forward slash backwards
				row = r.Intn(p.Height - len(word) + 1)
				col = r.Intn(p.Width-len(word)+1) + len(word) - 1
				rowDir = 1
				colDir = -1
			case 8:
				fallthrough
			case 9: //backslash forwards
				row = r.Intn(p.Height - len(word) + 1)
				col = r.Intn(p.Width - len(word) + 1)
				rowDir = 1
				colDir = 1
			case 10:
				fallthrough
			case 11: //backslash backwards
				row = r.Intn(p.Height-len(word)+1) + len(word) - 1
				col = r.Intn(p.Width-len(word)+1) + len(word) - 1
				rowDir = -1
				colDir = -1
			default:
				fmt.Println("Something strange occurred")
				return errors.New("Rand broke")
			}
			clear := true
			for i, letter := range word {
				check := p.Board[row+rowDir*i][col+colDir*i]
				if check != Empty && check != letter {
					clear = false
					break
				}
			}

			if clear {
				for i, letter := range word {
					p.Board[row+rowDir*i][col+colDir*i] = letter
					positions[letter] = append(positions[letter], LetterPosition{
						Row: row + rowDir*i,
						Col: col + colDir*i,
					})
				}
				success = true
			}

			iterations = iterations + 1
			if iterations > 1000 {
				return ErrMaxIterations
			}
		}

	}

	for i := 0; i < p.Height; i++ {
		for j := 0; j < p.Width; j++ {
			if p.Board[i][j] == Empty {
				value := r.Intn(26) + 65
				r, _ := utf8.DecodeRune([]byte{byte(value)})
				p.Board[i][j] = r
			}
		}
	}

	return nil
}
