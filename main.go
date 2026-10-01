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
				delim = strings.Trim(delim, "'\"")
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

func printResult(res Result) {
	for _, item := range res.Output {
		fmt.Printf("CMD %s\n", item.Cmd)
		fmt.Println("BODY:")
		for _, line := range item.Body {
			fmt.Println(line)
		}
		fmt.Println("END")
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	var input strings.Builder
	for scanner.Scan() {
		input.WriteString(scanner.Text())
		input.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	res, err := Execute(input.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR %s\n", err)
		os.Exit(1)
	}
	printResult(res)
}
