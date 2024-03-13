package service

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"meme/model"
)

// probability of being a reserved werd
func AppearanceCount(arg string, reserved_words ...string) map[string]int {
	//reserved word composition
	reserved_words_component := map[string]map[rune]struct{}{}
	//appearance count
	appearance_count := map[string]int{}

	FmtPrintln("blue", StrJoin(32, "arg: ", arg))

	//check the list of reserved words one by one
	for _, word := range reserved_words {
		component := map[rune]struct{}{}
		for _, r := range word {
			component[r] = struct{}{}
		}

		//record composition
		reserved_words_component[word] = component
		//initialize counter
		appearance_count[word] = 0
	}

	FmtPrint("blue", "components of reserved words: ")
	FmtPrintln("blue", reserved_words_component)

	//check arg one letter at a time
	for _, r := range arg {
		for word, component := range reserved_words_component {
			_, ok := component[r]
			if ok {
				appearance_count[word]++
			}
		}
	}

	FmtPrint("blue", "aggregate results: ")
	FmtPrintln("blue", appearance_count)
	return appearance_count
}

// convert map to struct
func ConvertMapToStruct(record map[string]string) any {
	len := len(record)
	dynamicFields := make([]reflect.StructField, len)

	//retrieve key
	keys := make([]string, len)
	{
		i := 0
		for k := range record {
			keys[i] = k
			i++
		}
	}

	//sort
	sort.Strings(keys)

	//define struct and retrieve val
	vals := make([]string, len)
	for i, v := range keys {
		dynamicFields[i] = reflect.StructField{
			Name: SnakeToPascal(v),
			Type: reflect.TypeOf(""),
			Tag:  reflect.StructTag(StrJoin(64, `json:"`, v, `"`)),
		}

		vals[i] = record[v]
	}

	//convert to struct
	structDef := reflect.New(reflect.StructOf(dynamicFields))
	structVal := structDef.Elem()

	//set val
	for i, v := range vals {
		structVal.Field(i).SetString(v)
	}

	return structVal.Interface()
}

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

// convert pascal case to snake case
func PascalToSnake(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([0-9a-z])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}

// convert snake case to pascal case
func SnakeToPascal(str string) string {
	r := make([]rune, len(str))
	flg := false

	for i, v := range str {
		if i == 0 {
			//A: 65, Z: 90, a: 97, z: 122
			if v >= 97 && v <= 122 {
				v -= 32
			}
		} else {
			if flg {
				if v >= 97 && v <= 122 {
					v -= 32
					flg = false
				}
			} else {
				//_: 95
				if v == 95 {
					flg = true
				}
			}
		}

		r[i] = v
	}

	return strings.Replace(string(r), "_", "", -1)
}

// sort [string]int map by value
func SortMapValue_StrInt(targe_map map[string]int) model.PairList_StrInt {
	pl := make(model.PairList_StrInt, len(targe_map))
	i := 0

	for k, v := range targe_map {
		pl[i] = model.Pair_StrInt{Key: k, Value: v}
		i++
	}

	sort.Sort(pl)

	return pl
}

// combine strings
//
// pass characters one at a time or pass a single array
func StrJoin(cap int, strs ...string) string {
	b := bytes.NewBuffer(make([]byte, 0, cap))

	for _, str := range strs {
		b.WriteString(str)
	}

	return b.String()
}

// combine strings
//
// the firs is a string, followed by the characters one by one or a single array
func StrJoin_2(cap int, str0 string, strs ...string) string {
	b := bytes.NewBuffer(make([]byte, 0, cap))

	b.WriteString(str0)
	for _, str := range strs {
		b.WriteString(" ")
		b.WriteString(str)
	}

	return b.String()
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
