package service

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"meme/analysis"
	"meme/model"
)

var addr = flag.String("addr", os.Getenv("API_Domain"), "http service address")
var Current_sub model.Value_StreamId
var Buffer = make(chan []byte)
var ReqFlg = make(chan bool)
var ExtFlg = make(chan bool)
var recvErrFlg = false
var waitStartTime = ""
var failedGuilds = map[int]model.FailedGuild{}
var failedGuildsKey = 0
var failedRecords = map[int]model.FailedRecord{}
var failedRecordsKey = 0
var WsGuild = make(chan model.FailedGuild)
var WsRecord = make(chan model.FailedRecord)
var WsFlg = false
var StreamConf model.StreamConf

func SetCurrentSub() {
	StreamConf = GetStreamConf()

	country_code := StreamConf.Country

	world_str := StrJoin(6, "000", StreamConf.World)
	world_str = world_str[len(world_str)-3:]

	group, _ := strconv.Atoi(StreamConf.Group)
	class, _ := strconv.Atoi(StreamConf.Class)
	block, _ := strconv.Atoi(StreamConf.Block)
	castle, _ := strconv.Atoi(StreamConf.Castle)

	if group != 0 && class != 0 {
		country_code = "0"
	}

	world, _ := strconv.Atoi(StrJoin(4, country_code, world_str))

	Current_sub = model.Value_StreamId{
		WorldId:  world,
		GroupId:  group,
		Class:    class,
		Block:    block,
		CastleId: castle,
	}
}

func GetCountryCode(country_name string) int {
	country_code := 0
	switch country_name {
	case "Japan":
		country_code = 1
	case "Korea":
		country_code = 2
	case "Asia":
		country_code = 3
	case "North America":
		country_code = 4
	case "Europe":
		country_code = 5
	case "Global":
		country_code = 6
	}

	return country_code
}

func GetCountryName(country_code int) string {
	country_name := ""
	switch country_code {
	case 1:
		country_name = "Japan"
	case 2:
		country_name = "Korea"
	case 3:
		country_name = "Asia"
	case 4:
		country_name = "North America"
	case 5:
		country_name = "Europe"
	case 6:
		country_name = "Global"
	}

	return country_name
}

func GetBuffer() []byte {
	return analysis.PackStreamId(Current_sub)
}

func Gvg() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/gvg"}
	fmt.Printf("connecting to %s", u.String())

	c, res, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err, "\n", res.StatusCode)
	}

	defer func() {
		fmt.Println("exit [Gvg]")
		c.Close()

		if recvErrFlg {
			recvErrFlg = false

			go Gvg()

			Buffer <- GetBuffer()
			<-ReqFlg
		}
	}()

	done := make(chan struct{})

	go recv(c, done)
	wait(c, done)
}

func RegisterGuild(guilds map[int]*model.Value_GuildId) {
	query := `
	insert into guilds (
		world_id,
		guild_id,
		guild_name
	) values %s
	on duplicate key update
		update_date = current_timestamp();
	`

	values := make([]string, 0, 20)
	for k, guild := range guilds {
		if guilds[k].Changed {

			value := `
			(
				%d,
				%d,
				'%s'
			)
			`
			value = fmt.Sprintf(value,
				guild.StreamId.WorldId,
				guild.GuildId,
				guild.GuildName)
			values = append(values, value)

			failed := model.FailedGuild{
				WorldId:    guild.StreamId.WorldId,
				GuildId:    guild.GuildId,
				GuildName:  guild.GuildName,
				CreateDate: time.Now(),
			}

			if WsFlg {
				WsGuild <- failed
			}

			failedGuilds[failedGuildsKey] = failed
			failedGuildsKey++

			guilds[k].Changed = false
		}
	}

	if len(values) > 0 {
		query = fmt.Sprintf(query, strings.Join(values, ","))
		query = strings.Join(strings.Fields(query), " ")

		err := ExecQuery(query)

		if err != nil {
			LogPrint("red", "exec_query", err)
		} else {
			for i := 0; i < len(guilds); i++ {
				failedGuildsKey--
				delete(failedGuilds, failedGuildsKey)
			}
		}
	}
}

func RegisterRecord(castles map[int]*model.Value_CastleId) {

	for k, castle := range castles {

		if castle.StreamId.WorldId == 0 {
			continue
		}

		if castles[k].Changed {
			query := `
			call registering_record(
				%d,
				%d,
				%d,
				%d,
				%d,
				%d,
				%d,
				%d,
				%d,
				%d,
				%d,
				'%s'
			);
			`

			now := time.Now()

			query = fmt.Sprintf(query,
				castle.StreamId.WorldId,
				castle.StreamId.GroupId,
				castle.StreamId.Class,
				castle.StreamId.Block,
				castle.StreamId.CastleId,
				castle.GuildId,
				castle.AttackerGuildId,
				castle.UtcFallenTimeStamp,
				castle.DefensePartyCount,
				castle.AttackPartyCount,
				castle.GvgCastleState,
				now.Format("2006/01/02 15:04:05"))

			query = strings.Join(strings.Fields(query), " ")

			failed := model.FailedRecord{
				WorldId:            castle.StreamId.WorldId,
				GroupId:            castle.StreamId.GroupId,
				Class:              castle.StreamId.Class,
				Block:              castle.StreamId.Block,
				CastleId:           castle.StreamId.CastleId,
				GuildId:            castle.GuildId,
				AttackerGuildId:    castle.AttackerGuildId,
				UtcFallenTimeStamp: castle.UtcFallenTimeStamp,
				DefensePartyCount:  castle.DefensePartyCount,
				AttackPartyCount:   castle.AttackPartyCount,
				GvgCastleState:     castle.GvgCastleState,
				CreateDate:         now,
			}

			if WsFlg {
				WsRecord <- failed
			}

			err := ExecQuery(query)

			if err != nil {
				LogPrint("red", "exec_query", err)

				failedRecords[failedRecordsKey] = failed
				failedGuildsKey++
			}

			castles[k].Changed = false
		}
	}
}

func recv(c *websocket.Conn, done chan struct{}) {
	defer func() {
		fmt.Println("exit [recv]")
		close(done)
	}()

	go func() {
		defer log.Println("exit [recv-for]")
		for {
			_, message, err := c.ReadMessage()

			if err != nil {
				log.Println("read:", err)

				recvErrFlg = true
				return
			}

			analysis.GvgAnalysis(message, Current_sub)
			RegisterGuild(analysis.Guilds)
			RegisterRecord(analysis.Castles)

			waitStartTime = time.Now().Format("2006/01/02 15:04:05")
		}
	}()
	<-ExtFlg
}

func wait(c *websocket.Conn, done chan struct{}) {
	/*
	 * Time allowed to read the next pong message from the peer
	 * pongWait = 60
	 *
	 * Send pings to peer with this period. Must be less than pongWait.
	 * pingPeriod = (pongWait * 9) / 10
	 *            = (60 * 9) / 10
	 *            = 540 / 10
	 *            = 54
	 */
	ticker := time.NewTicker(50 * time.Second)
	defer ticker.Stop()
	defer fmt.Println("exit [wait]")

	for {
		select {
		case <-done:
			fmt.Println("done")
			return
		case b := <-Buffer:
			w, err := c.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Println("writer b:", err)
			}
			_, err = w.Write(b)
			w.Close()
			if err != nil {
				log.Println("write b:", err)
				return
			}

			fmt.Println(b)
			ReqFlg <- false
		case <-ticker.C:
			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("write t:", err)
				return
			}
			fmt.Printf("\r%s -> %s", waitStartTime, time.Now().Format("2006/01/02 15:04:05"))
		}
	}
}
