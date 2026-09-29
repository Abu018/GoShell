package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var redirectionOperators = map[string]bool{
	"<":    true, // stdin from file
	"<<":   true, // here-document
	"<<<":  true, // here-string
	"<&":   true, // duplicate/close stdin
	">":    true, // stdout to file (truncate)
	">>":   true, // stdout to file (append)
	">|":   true, // stdout, force overwrite (noclobber override)
	">&":   true, // duplicate/close stdout
	"2>":   true, // stderr to file (truncate)
	"2>>":  true, // stderr to file (append)
	"2>&1": true, // stderr to wherever stdout points
	"2>&-": true, // close stderr
	"&>":   true, // stdout+stderr to file (truncate) - bash-ism
	"&>>":  true, // stdout+stderr to file (append) - bash-ism
	"<>":   true, // open file for read+write on stdin
	"|":    true, // pipe
	"|&":   true, // pipe stdout+stderr - bash-ism
}

type Redirect struct {
	FD     int
	Op     string
	Target string
}

type Result struct {
	Token     []string
	Redirects []Redirect
}

var result Result

func splitting_txt(input string) ([][]string, error) {
	var result [][]string
	temp := ""
	var token []string
	for ch := 0; ch < len(input); ch++ {
		if input[ch] != ' ' {
			temp += string(input[ch])
		} else if input[ch] == ' ' && !redirectionOperators[string(ch)] {
			token = append(token, temp)
			result = append(result, token)
			temp = ""
			token = []string{}
			continue
		}
	}
	if temp != "" {
		token = append(token, temp)
		result = append(result, token)
	}
	return result, nil
}

var count = 0

func tokens_txt(input [][]string) error {

	for _, wordGroup := range input {
		for _, word := range wordGroup {
			if strings.Contains(word, ".") {
				continue
			}
			if !redirectionOperators[word] {
				result.Token = append(result.Token, word)
			}
		}
	}
	fmt.Println(result.Token)

	return nil
}

func parseRedirect(tok string) (fd int, op, target string, ok bool) {
	ops := []string{">>", ">", "<"}
	j := 0
	for j < len(tok) && tok[j] >= '0' && tok[j] <= '9' {
		j++
	}
	rest := tok[j:]
	for _, o := range ops {
		if strings.HasPrefix(rest, o) {
			op = o
			break
		}
	}
	if op == "" {
		return 0, "", "", false
	}
	if j > 0 {
		fmt.Sscanf(tok[:j], "%d", &fd)
	} else if op == "<" {
		fd = 0
	} else {
		fd = 1
	}
	return fd, op, rest[len(op):], true
}

func Execute(input string) (Result, error) {
	var out Result
	fields := strings.Fields(input)

	for i := 0; i < len(fields); i++ {
		fd, op, target, ok := parseRedirect(fields[i])
		if !ok {
			out.Token = append(out.Token, fields[i])
			continue
		}
		if target == "" {
			if i+1 >= len(fields) {
				return Result{}, fmt.Errorf("missing target after %q", op)
			}
			i++
			target = fields[i]
		}
		out.Redirects = append(out.Redirects, Redirect{FD: fd, Op: op, Target: target})
	}
	return out, nil
}
func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		result, err := Execute(sc.Text())
		if err != nil {
			fmt.Printf("ERR %s", err)
			continue
		}
		fmt.Printf("argv=['%s']\n", strings.Join(result.Token, "', '"))
		for _, r := range result.Redirects {
			fmt.Printf("redir fd=%d op=%s target=%s\n", r.FD, r.Op, r.Target)
		}
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
