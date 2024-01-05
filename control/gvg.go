package control

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"meme/model"
	"meme/service"
)

func Gvg(c *gin.Context) {
	path, _ := c.Get("FullPath")

	c.HTML(http.StatusOK, "gvg_.html.tmpl", gin.H{
		"Path":   path,
		"Region": "gvg",
	})
}

func GvgStream(c *gin.Context) {
	var cmd model.Cmd
	c.Bind(&cmd)

	var res gin.H

	switch cmd.App {
	case "check":
		what := "--id"
		if len(cmd.Args) > 0 {
			what = cmd.Args[0]
		}

		if what == "--id" {
			res = check_stream(cmd.Args)

		} else if what == "--wss" {
			//
		}
	case "set":
		res = set_stream(cmd.Args)
	}

	c.JSON(http.StatusOK, res)
}

func check_stream(args []string) gin.H {
	sc := service.GetStreamConf()

	return gin.H{
		"Country": sc.Country,
		"World":   sc.World,
		"Group":   sc.Group,
		"Class":   sc.Class,
		"Block":   sc.Block,
		"Castle":  sc.Castle,
	}
}

func set_stream(args []string) gin.H {
	country := "Japan"
	world := 0
	group := 0
	class := 0
	block := 0
	castle := 0
	var err error

	for i, arg := range args {
		switch arg {
		case "--country":
			if i < len(args)-1 {
				if args[i+1][0] != 45 {
					country = args[i+1]
				}
			}
		case "--world":
			if i < len(args)-1 {
				if args[i+1][0] != 45 {
					world, err = strconv.Atoi(args[i+1])

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		case "--group":
			if i < len(args)-1 {
				if args[i+1][0] != 45 {
					group, err = strconv.Atoi(args[i+1])

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		case "--class":
			if i < len(args)-1 {
				if args[i+1][0] != 45 {
					class, err = strconv.Atoi(args[i+1])

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		case "--block":
			if i < len(args)-1 {
				if args[i+1][0] != 45 {
					block, err = strconv.Atoi(args[i+1])

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		case "--castle":
			if i < len(args)-1 {
				if args[i+1][0] != 45 {
					castle, err = strconv.Atoi(args[i+1])

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		default:
			//
			appearance_count := appearance_count(arg)
			pl := service.SortMapValue(appearance_count)

			service.FmtPrintln("red", service.StrJoin(64, "\"", arg, "\"", " is not a valid argument\nmaybe: \"--", pl[5].Key, "\""))

			return gin.H{
				service.StrJoin(32, "\"", arg, "\"", " is not a valid argument"): "",
				service.StrJoin(32, "maybe: \"--", pl[5].Key, "\""):              "",
			}
		}
	}

	service.SetStreamConf(country, world, group, class, block, castle)

	return gin.H{
		"Country": country,
		"World":   world,
		"Group":   group,
		"Class":   class,
		"Block":   block,
		"Castle":  castle,
	}
}

func appearance_count(arg string) map[string]int {
	country := map[rune]struct{}{45: {}, 99: {}, 110: {}, 111: {}, 116: {}, 117: {}, 121: {}}
	world := map[rune]struct{}{45: {}, 100: {}, 108: {}, 111: {}, 114: {}, 119: {}}
	group := map[rune]struct{}{45: {}, 103: {}, 111: {}, 112: {}, 114: {}, 117: {}}
	class := map[rune]struct{}{45: {}, 97: {}, 99: {}, 108: {}, 115: {}}
	block := map[rune]struct{}{45: {}, 98: {}, 99: {}, 107: {}, 108: {}, 111: {}}
	castle := map[rune]struct{}{45: {}, 97: {}, 99: {}, 101: {}, 108: {}, 115: {}, 116: {}}

	appearance_count := map[string]int{}
	appearance_count["country"] = 0
	appearance_count["world"] = 0
	appearance_count["group"] = 0
	appearance_count["class"] = 0
	appearance_count["block"] = 0
	appearance_count["castle"] = 0

	for _, b := range arg {
		_, ok := country[b]
		if ok {
			appearance_count["country"]++
		}
		_, ok = world[b]
		if ok {
			appearance_count["world"]++
		}
		_, ok = group[b]
		if ok {
			appearance_count["group"]++
		}
		_, ok = class[b]
		if ok {
			appearance_count["class"]++
		}
		_, ok = block[b]
		if ok {
			appearance_count["block"]++
		}
		_, ok = castle[b]
		if ok {
			appearance_count["castle"]++
		}
	}

	return appearance_count
}
