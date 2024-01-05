package service

import (
	"bytes"
	"fmt"
	"log"
)

// fmt.Printf
//
// black,
// red,
// green,
// yellow,
// blue,
// magenta,
// cyan,
// white,
// bright_black_gray,
// bright_red,
// bright_green,
// bright_yellow,
// bright_blue,
// bright_magenta,
// bright_cyan,
// bright_white
func FmtPrint(color string, val any) {
	fmt.Printf("\x1b[%dm%v\x1b[0m", color_code(color), val)
}
func FmtPrintln(color string, val any) {
	fmt.Printf("\x1b[%dm%v\x1b[0m\n", color_code(color), val)
}

// log.Printf
//
// black,
// red,
// green,
// yellow,
// blue,
// magenta,
// cyan,
// white,
// bright_black_gray,
// bright_red,
// bright_green,
// bright_yellow,
// bright_blue,
// bright_magenta,
// bright_cyan,
// bright_white
func LogPrint(color string, title string, err error) {
	log.Printf("\x1b[%dm%s: %v\x1b[0m", color_code(color), title, err)
}
func LogPrintln(color string, title string, err error) {
	log.Printf("\x1b[%dm%s: %v\x1b[0m\n", color_code(color), title, err)
}

// return color code
func color_code(color string) int {
	code := 0

	//Terminal Text Color Codes
	switch color {
	case "black":
		code = 30
	case "red":
		code = 31
	case "green":
		code = 32
	case "yellow":
		code = 33
	case "blue":
		code = 34
	case "magenta":
		code = 35
	case "cyan":
		code = 36
	case "white":
		code = 37
	case "bright_black_gray":
		code = 90
	case "bright_red":
		code = 91
	case "bright_green":
		code = 92
	case "bright_yellow":
		code = 93
	case "bright_blue":
		code = 94
	case "bright_magenta":
		code = 95
	case "bright_cyan":
		code = 96
	case "bright_white":
		code = 97
	}

	return code
}

func StrJoin(cap int, strs ...string) string {
	b := bytes.NewBuffer(make([]byte, 0, cap))

	for _, str := range strs {
		b.WriteString(str)
	}

	return b.String()
}

func StrJoin_2(cap int, str0 string, strs ...string) string {
	b := bytes.NewBuffer(make([]byte, 0, cap))

	b.WriteString(str0)
	for _, str := range strs {
		b.WriteString(" ")
		b.WriteString(str)
	}

	return b.String()
}
