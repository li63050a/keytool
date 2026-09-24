package menu

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

func ReadLine(prompt string) string {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Print(prompt)
		var s string
		fmt.Scanln(&s)
		return s
	}
	defer term.Restore(fd, oldState)

	buf := []rune{}
	cursor := 0

	redraw := func() {
		fmt.Print("\r\033[2K")
		fmt.Print(prompt)
		fmt.Print(string(buf))
		if cursor < len(buf) {
			fmt.Printf("\033[%dD", len(buf)-cursor)
		}
	}
	redraw()

	tmp := make([]byte, 8)
	for {
		n, err := os.Stdin.Read(tmp)
		if err != nil || n == 0 {
			return ""
		}

		if tmp[0] == 0x1b {
			if n >= 3 && tmp[1] == '[' {
				switch tmp[2] {
				case 'C':
					if cursor < len(buf) {
						cursor++
						fmt.Print("\033[C")
					}
				case 'D':
					if cursor > 0 {
						cursor--
						fmt.Print("\033[D")
					}
				case 'H':
					cursor = 0
					redraw()
				case 'F':
					cursor = len(buf)
					redraw()
				case 'A', 'B':
				case '3':
					if cursor < len(buf) {
						buf = append(buf[:cursor], buf[cursor+1:]...)
						redraw()
					}
				}
			}
			continue
		}

		if n == 1 {
			switch tmp[0] {
			case '\r', '\n':
				fmt.Print("\r\n")
				return strings.TrimSpace(string(buf))
			case 3: // Ctrl+C：恢复终端 + 退出
				fmt.Print("\r\n")
				term.Restore(fd, oldState)
				fmt.Println("已取消。")
				os.Exit(0)
			case 127, 8:
				if cursor > 0 {
					buf = append(buf[:cursor-1], buf[cursor:]...)
					cursor--
					redraw()
				}
			case 21:
				buf = buf[:0]
				cursor = 0
				redraw()
			case 4:
				if len(buf) == 0 {
					fmt.Print("\r\n")
					return ""
				}
			default:
				if tmp[0] >= 32 {
					if tmp[0] < 0x80 {
						insertRune(&buf, &cursor, rune(tmp[0]))
						redraw()
					} else {
						size := 2
						if tmp[0] >= 0xf0 {
							size = 4
						} else if tmp[0] >= 0xe0 {
							size = 3
						}
						charBytes := make([]byte, size)
						charBytes[0] = tmp[0]
						for i := 1; i < size; i++ {
							one := make([]byte, 1)
							os.Stdin.Read(one)
							charBytes[i] = one[0]
						}
						if r, _ := utf8.DecodeRune(charBytes); r != utf8.RuneError {
							insertRune(&buf, &cursor, r)
							redraw()
						}
					}
				}
			}
		}
	}
}

func insertRune(buf *[]rune, cursor *int, r rune) {
	*buf = append(*buf, 0)
	copy((*buf)[*cursor+1:], (*buf)[*cursor:])
	(*buf)[*cursor] = r
	*cursor++
}
