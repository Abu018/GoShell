package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Execute(input string) (string, error) {
	trimmedInput := strings.TrimSpace(input)
	trimmedInputLength := len(trimmedInput)
	var result []string
	temp := ""

	if trimmedInputLength == 0 {
		return "", fmt.Errorf("syntax error: empty command in pipeline")
	}
	if trimmedInput[0] == '|' || trimmedInput[len(trimmedInput)-1] == '|' {
		return "", fmt.Errorf("syntax error: empty command in pipeline")
	}
	for i := 0; i < trimmedInputLength; i++ {
		if trimmedInput[i] == '|' && i+1 < len(trimmedInput) && trimmedInput[i+1] == '|' {
			return "", fmt.Errorf("syntax error: empty command in pipeline")
		} else if trimmedInput[i] == ' ' {
			result = append(result, temp)
			temp = ""
		} else {
			if i > 1 && trimmedInput[i-1] == ' ' {
				result = append(result, " ")
			}
			temp += string(trimmedInput[i])
		}
	}
	if temp != "" {
		result = append(result, temp)
	}
	return strings.Join(result, ""), nil
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
			fmt.Printf("ERR %s", err)

			continue
		}
		fmt.Println(result)
		// for i, word := range result {
		// 	if i > 0 {
		// 		fmt.Print(" ")
		// 	}
		// 	fmt.Printf(word)
		// }
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
