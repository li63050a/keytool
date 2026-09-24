package menu

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func Select(prompt string, options []string) int {
	return SelectWithDefault(prompt, options, 0)
}

func SelectWithDefault(prompt string, options []string, defaultIdx int) int {
	if len(options) == 0 {
		return -1
	}
	if defaultIdx < 0 || defaultIdx >= len(options) {
		defaultIdx = 0
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return selectFallback(prompt, options, defaultIdx)
	}
	defer term.Restore(fd, oldState)

	cursor := defaultIdx

	fmt.Print(prompt + "\r\n")
	for i, opt := range options {
		renderOption(opt, i == cursor)
		fmt.Print("\r\n")
	}

	buf := make([]byte, 8)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return -1
		}

		old := cursor

		if buf[0] == 0x1b {
			if n >= 3 && buf[1] == '[' {
				switch buf[2] {
				case 'A':
					if cursor > 0 {
						cursor--
					}
				case 'B':
					if cursor < len(options)-1 {
						cursor++
					}
				}
			} else if n == 1 {
				return -1
			}
		} else if n == 1 {
			switch buf[0] {
			case '\r', '\n':
				return cursor
			case 3:
				return -1
			case 'q', 'Q':
				return -1
			case 'k', 'K':
				if cursor > 0 {
					cursor--
				}
			case 'j', 'J':
				if cursor < len(options)-1 {
					cursor++
				}
			}
		}

		if old != cursor {
			redrawOptions(options, cursor)
		}
	}
}

func redrawOptions(options []string, cursor int) {
	fmt.Printf("\033[%dA", len(options))
	for i, opt := range options {
		fmt.Print("\r\033[2K")
		renderOption(opt, i == cursor)
		fmt.Print("\r\n")
	}
}

func renderOption(opt string, selected bool) {
	if selected {
		fmt.Printf(" \033[44;97m> %s \033[0m", opt)
	} else {
		fmt.Printf("   %s", opt)
	}
}

func selectFallback(prompt string, options []string, defaultIdx int) int {
	fmt.Println(prompt)
	for i, opt := range options {
		mark := " "
		if i == defaultIdx {
			mark = ">"
		}
		fmt.Printf(" %s [%d] %s\n", mark, i+1, opt)
	}
	s := ReadLine("请输入数字: ")
	var n int
	fmt.Sscanf(s, "%d", &n)
	if n < 1 || n > len(options) {
		return -1
	}
	return n - 1
}
