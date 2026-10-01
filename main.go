package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Output struct {
	Cmd  string
	Body []string
	End  string
}
type Result struct {
	Output []Output
}

func Execute(input string) (Result, error) {
	var result Result // local instead of global
	var cmd string
	var temp []string
	inHeredoc := false
	delim := ""

	for _, line := range strings.Split(input, "\n") {
		// 2. split each line into words
		words := strings.Fields(line)
		if len(words) == 0 {
			continue
		}

		// 3. inside a heredoc: collect body until the delimiter line
		if inHeredoc {
			if len(words) == 1 && words[0] == delim {
				result.Output = append(result.Output, Output{
					Cmd:  cmd,
					Body: temp, // 4. use temp, not the empty Body
					End:  delim,
				})
				inHeredoc = false
				temp = nil
			} else {
				temp = append(temp, strings.Join(words, " "))
			}
			continue
		}

		opened := false
		for _, w := range words {
			if strings.HasPrefix(w, "<<") {
				delim = strings.TrimPrefix(strings.TrimPrefix(w, "<<"), "-")
				cmd = strings.Join(words, " ")
				inHeredoc = true
				opened = true
				break
			}
		}
		if !opened {
			result.Output = append(result.Output, Output{Cmd: strings.Join(words, " ")})
		}
	}

	if inHeredoc {
		return Result{}, fmt.Errorf("unterminated heredoc %q", delim)
	}
	return result, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		res, err := Execute(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERR %s\n", err)
			os.Exit(1)
		}
		fmt.Print(res)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
