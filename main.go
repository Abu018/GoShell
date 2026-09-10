package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func Execute(input string) ([]string, error) {
	var result []string
	var temp strings.Builder
	inToken := false
	inSingle := false
	inDouble := false

	flush := func() {
		if inToken {
			result = append(result, temp.String())
			temp.Reset()
			inToken = false
		}
	}

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if inSingle {
			if ch == '\'' {
				inSingle = false
			} else {
				temp.WriteByte(ch)
			}
			continue
		}

		if inDouble {
			if ch == '"' {
				inDouble = false
				continue
			}
			if ch == '\\' && i+1 < len(input) {
				n := input[i+1]
				if n == '"' || n == '\\' || n == '$' || n == '`' {
					temp.WriteByte(n)
					i++
					continue
				}
			}
			temp.WriteByte(ch)
			continue
		}

		if ch == '\\' {
			inToken = true
			if i+1 < len(input) {
				temp.WriteByte(input[i+1])
				i++
			}
			continue
		}

		if ch == '\'' {
			inSingle = true
			inToken = true
			continue
		}
		if ch == '"' {
			inDouble = true
			inToken = true
			continue
		}

		if unicode.IsSpace(rune(ch)) {
			flush()
			continue
		}

		inToken = true
		temp.WriteByte(ch)
	}

	if inSingle || inDouble {
		return nil, fmt.Errorf("unterminated quote")
	}
	flush()
	return result, nil
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for {
		if !sc.Scan() {
			break
		}
		result, err := Execute(sc.Text())
		if err != nil {
			fmt.Print("ERR unterminated quote")
			continue
		}
		for i, word := range result {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Printf("[%s]", word)
		}
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
