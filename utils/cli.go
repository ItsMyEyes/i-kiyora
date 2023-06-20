package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func MakeInputCli(question string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(question + " ")
	scanner.Scan()
	input := scanner.Text()
	input = strings.TrimSuffix(input, "\n")
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
		return MakeInputCli(question)
	}
	return input
}

func NCenter(width int, s string) *bytes.Buffer {
	const half, space = 2, "\u0020"
	var b bytes.Buffer
	n := (width - utf8.RuneCountInString(s)) / half
	fmt.Fprintf(&b, "%s%s", strings.Repeat(space, n), s)
	return &b
}

func MakeLine() {
	var line string = "-"
	for i := 0; i < 70; i++ {
		line += "-"
	}
	fmt.Println("\n" + line + "\n")
}

func getPathSlash() string {
	if os.PathSeparator == '\\' {
		return "\\"
	}
	return "/"
}

func MakeDirectoryString(dir ...string) string {
	var result string
	for _, d := range dir {
		result += d + getPathSlash()
	}
	return result
}
