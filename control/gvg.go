package control

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"meme/model"
	"meme/service"
)

// processing of GET request for grand battle page
func Ggvg(c *gin.Context) {
	path, _ := c.Get("FullPath")

	c.HTML(http.StatusOK, "ggvg.html.tmpl", gin.H{
		"Path":   path,
		"Region": "ggvg",
	})
}

// processing of GET request for guild battle page
func Gvg(c *gin.Context) {
	path, _ := c.Get("FullPath")

	c.HTML(http.StatusOK, "gvg.html.tmpl", gin.H{
		"Path":   path,
		"Region": "gvg",
	})
}

// processing of POST request for guild battle
func GvgStream(c *gin.Context) {
	//bind received command arguments to model
	var cmd model.Cmd
	c.Bind(&cmd)

	var res gin.H

	switch cmd.App {
	case "check":
		//confirmation of stream_id or websocket
		what := "--id"
		if len(cmd.Args) > 0 {
			what = cmd.Args[0]
		}

		if what == "--id" {
			//stream_id
			res = check_stream(cmd.Args)

		} else if what == "--wss" {
			//websocket
			res = gin.H{
				"ws": "",
			}
		}
	case "set":
		//set stream_id
		res = set_stream(cmd.Args)
		service.FmtPrintln("blue", "set stream_id")

	case "start":
		//start streaming

		//set current stream_id
		service.SetCurrentSub()
		service.FmtPrintln("blue", "set current stream_id")

		//open websocket
		go service.Gvg()
		service.FmtPrintln("blue", "open websocket")

		//send stream_id to start streaming
		service.Buffer <- service.GetBuffer()
		<-service.ReqFlg

		//change StreamConf to streaming
		var err error
		country, err := strconv.Atoi(service.StreamConf.Country)
		if err != nil {
			service.LogPrintln("red", "Atoi", err)
		}

		world, err := strconv.Atoi(service.StreamConf.World)
		if err != nil {
			service.LogPrintln("red", "Atoi", err)
		}

		group, err := strconv.Atoi(service.StreamConf.Group)
		if err != nil {
			service.LogPrintln("red", "Atoi", err)
		}

		class, err := strconv.Atoi(service.StreamConf.Class)
		if err != nil {
			service.LogPrintln("red", "Atoi", err)
		}

		block, err := strconv.Atoi(service.StreamConf.Block)
		if err != nil {
			service.LogPrintln("red", "Atoi", err)
		}

		castle, err := strconv.Atoi(service.StreamConf.Castle)
		if err != nil {
			service.LogPrintln("red", "Atoi", err)
		}

		err = service.SetStreamConf(country, world, group, class, block, castle, true)

		if err != nil {
			service.LogPrintln("red", "SetStreamConf", err)
		}

		service.FmtPrintln("blue", "start streaming")

		res = gin.H{
			"start": [...]string{
				service.StrJoin(22, "country: ", service.GetCountryName(service.Current_sub.WorldId/1000)),
				service.StrJoin(22, "world: ", strconv.Itoa(service.Current_sub.WorldId)),
			},
		}
	default:
		//others are not accepted
		service.FmtPrintln("blue", service.StrJoin(34, "disapproval command: ", cmd.App))
		res = gin.H{
			"failed": [...]string{
				service.StrJoin(128, "'", cmd.App, "' is not recognized as an internal or external command,"),
				"operable program or batch file",
			},
		}
	}

	c.JSON(http.StatusOK, res)
}

// migration to server
func GvgTranRegist(c *gin.Context) {
	//bind received command arguments to model
	var cmd model.Cmd
	c.Bind(&cmd)

	if cmd.App == "upload" {
		ins := "INSERT INTO gvg_records(world_id,group_id,class,block,castle_id,def_guild_id,atk_guild_id,utc_fallen_time_stamp,def_count,atk_count,state,create_date,update_date)VALUES"
		query := service.StrJoin(135168, ins, cmd.Args[0])

		err := service.ExecQuery(query)

		if err != nil {
			service.FmtPrintln("red", query)
			service.LogPrint("red", "TranRegist", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"failed": [1]string{err.Error()},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"": service.StrJoin(24, "process completed", ": ", cmd.Args[1]),
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"app":  cmd.App,
			"args": cmd.Args,
		})
	}
}

// check stream_id
func check_stream(args []string) gin.H {
	//get from dynamoDb
	sc := service.GetStreamConf()

	country, err := strconv.Atoi(sc.Country)

	if err != nil {
		service.LogPrintln("red", "Atoi", err)
	}

	return gin.H{
		"Country": service.GetCountryName(country),
		"World":   sc.World,
		"Group":   sc.Group,
		"Class":   sc.Class,
		"Block":   sc.Block,
		"Castle":  sc.Castle,
		"Status":  sc.Status,
	}
}

// handling invalid option
func invalid_option(reserved_words []string, arg string) gin.H {
	//presenting the maybe
	appearance_count := service.AppearanceCount(arg, reserved_words...)
	pl := service.SortMapValue_StrInt(appearance_count)
	i := len(reserved_words) - 1

	service.FmtPrintln("bluee", service.StrJoin(128, "\"", arg, "\"", " is not a valid argument\nmaybe: \"", pl[i].Key, "\""))

	return gin.H{
		service.StrJoin(64, "\"", arg, "\"", " is not a valid argument"): "",
		service.StrJoin(64, "maybe: \"", pl[i].Key, "\""):                "",
	}
}

// set stream_id
func set_stream(args []string) gin.H {
	country := 1
	world := 1
	group := 0
	class := 0
	block := 0
	castle := 0
	status := false
	var err error

	//check options
	service.FmtPrintln("blue", "check options")
	for i, arg := range args {
		val := args[i+1]

		//put into the corresponding variable
		switch arg {
		case "--country":
			if i < len(args)-1 {
				if val[0] != 45 {
					if len(val) != 1 {
						//get country code
						country = service.GetCountryCode(val)

						//if country name is wrong
						if country == 0 {
							//reserved words
							reserved_words := []string{
								"Japan",
								"Korea",
								"Asia",
								"North America",
								"Europe",
								"Global",
							}

							return invalid_option(reserved_words, val)
						}
					}
				}
			}
		case "--world":
			if i < len(args)-1 {
				if val[0] != 45 {
					world, err = strconv.Atoi(val)

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		case "--group":
			if i < len(args)-1 {
				if val[0] != 45 {
					group, err = strconv.Atoi(val)

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		case "--class":
			if i < len(args)-1 {
				if val[0] != 45 {
					class, err = strconv.Atoi(val)

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		case "--block":
			if i < len(args)-1 {
				if val[0] != 45 {
					block, err = strconv.Atoi(val)

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		case "--castle":
			if i < len(args)-1 {
				if val[0] != 45 {
					castle, err = strconv.Atoi(val)

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		case "--status":
			if i < len(args)-1 {
				if val[0] != 45 {
					status, err = strconv.ParseBool(val)

					if err != nil {
						service.LogPrintln("red", "Atoi", err)
					}
				}
			}
		default:
			//invalid option

			//reserved words
			reserved_words := []string{
				"--Country",
				"--World",
				"--Group",
				"--Class",
				"--Block",
				"--Castle",
				"--Status",
			}

			return invalid_option(reserved_words, arg)
		}
	}

	//update dynamoDb
	err = service.SetStreamConf(country, world, group, class, block, castle, status)

	if err != nil {
		service.LogPrintln("red", "SetStreamConf", err)
	}

	return gin.H{
		"Country": service.GetCountryName(country),
		"World":   world,
		"Group":   group,
		"Class":   class,
		"Block":   block,
		"Castle":  castle,
		"Status":  status,
	}
}
