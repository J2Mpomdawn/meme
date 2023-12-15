package model

type StreamId struct {
	Value  Value_StreamId
	Offset int
}

type Value_StreamId struct {
	WorldId  int
	GroupId  int
	Class    int
	Block    int
	CastleId int
}

type GuildId struct {
	Value  Value_GuildId
	Offset int
}

type Value_GuildId struct {
	StreamId  Value_StreamId
	GuildId   int
	GuildName string
}

type CastleId struct {
	Value  Value_CastleId
	Offset int
}

type Value_CastleId struct {
	StreamId           Value_StreamId
	GuildId            int
	AttackerGuildId    int
	UtcFallenTimeStamp int
	DefensePartyCount  int
	AttackPartyCount   int
	GvgCastleState     int
}

type Hoge struct {
	Status    int `json:"status"`
	Timestamp int `json:"timestamp"`
	Data      struct {
		WorldID int `json:"world_id`
		Castles []struct {
			CastleId int `json:"CastleId"`
			GuildId  int `json:"GuildId"`
		} `json:"castles"`
	} `json:"data"`
}

var Jsonstr = `
{
	"status":200,
	"timestamp":1702046684,
	"data":{
		"world_id":1099,
		"castles":[
			{
				"CastleId":1,
				"GuildId":199652669099
			},
			{
				"CastleId":2,
				"GuildId":199652669099
			}
		]
	}
}`

/*

	"status":200,
	"timestamp":1702046684,
	"data":{
	   "world_id":1099,
	   "castles":[
		  {
			 "CastleId":1,
			 "GuildId":199652669099,{
			 "AttackerGuildId":394114239099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":46,
			 "GvgCastleState":0,
			 "UtcFallenTimeStamp":0
		  },
		  {
			 "CastleId":2,
			 "GuildId":494634944099,
			 "AttackerGuildId":908515237099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":153,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702069873102
		  },
		  {
			 "CastleId":3,
			 "GuildId":394114239099,
			 "AttackerGuildId":199652669099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":159,
			 "GvgCastleState":0,
			 "UtcFallenTimeStamp":0
		  },
		  {
			 "CastleId":4,
			 "GuildId":129044010099,
			 "AttackerGuildId":908515237099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":2,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702070217101
		  },
		  {
			 "CastleId":5,
			 "GuildId":199652669099,
			 "AttackerGuildId":494634944099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":49,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702070041100
		  },
		  {
			 "CastleId":6,
			 "GuildId":494634944099,
			 "AttackerGuildId":741934801099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":28,
			 "GvgCastleState":0,
			 "UtcFallenTimeStamp":0
		  },
		  {
			 "CastleId":7,
			 "GuildId":494634944099,
			 "AttackerGuildId":530715127099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":531,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702069650096
		  },
		  {
			 "CastleId":8,
			 "GuildId":530715127099,
			 "AttackerGuildId":768981549099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":3,
			 "GvgCastleState":4,
			 "UtcFallenTimeStamp":1702069273100
		  },
		  {
			 "CastleId":9,
			 "GuildId":768981549099,
			 "AttackerGuildId":530715127099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":10,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702070752100
		  },
		  {
			 "CastleId":10,
			 "GuildId":908515237099,
			 "AttackerGuildId":494634944099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":45,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702069866101
		  },
		  {
			 "CastleId":11,
			 "GuildId":720141187099,
			 "AttackerGuildId":669198682099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":452,
			 "GvgCastleState":0,
			 "UtcFallenTimeStamp":0
		  },
		  {
			 "CastleId":12,
			 "GuildId":394114239099,
			 "AttackerGuildId":816626613099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":107,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702070014102
		  },
		  {
			 "CastleId":13,
			 "GuildId":816626613099,
			 "AttackerGuildId":669198682099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":221,
			 "GvgCastleState":0,
			 "UtcFallenTimeStamp":0
		  },
		  {
			 "CastleId":14,
			 "GuildId":394114239099,
			 "AttackerGuildId":816626613099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":99,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702070019096
		  },
		  {
			 "CastleId":15,
			 "GuildId":394114239099,
			 "AttackerGuildId":719463286099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":65,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702070056096
		  },
		  {
			 "CastleId":16,
			 "GuildId":719463286099,
			 "AttackerGuildId":129044010099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":68,
			 "GvgCastleState":0,
			 "UtcFallenTimeStamp":0
		  },
		  {
			 "CastleId":17,
			 "GuildId":404400524099,
			 "AttackerGuildId":719463286099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":15,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702070170101
		  },
		  {
			 "CastleId":18,
			 "GuildId":482173199099,
			 "AttackerGuildId":0,
			 "AttackPartyCount":0,
			 "DefensePartyCount":235,
			 "GvgCastleState":0,
			 "UtcFallenTimeStamp":0
		  },
		  {
			 "CastleId":19,
			 "GuildId":908515237099,
			 "AttackerGuildId":129044010099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":176,
			 "GvgCastleState":0,
			 "UtcFallenTimeStamp":0
		  },
		  {
			 "CastleId":20,
			 "GuildId":196765134099,
			 "AttackerGuildId":741934801099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":603,
			 "GvgCastleState":0,
			 "UtcFallenTimeStamp":0
		  },
		  {
			 "CastleId":21,
			 "GuildId":741934801099,
			 "AttackerGuildId":196765134099,
			 "AttackPartyCount":0,
			 "DefensePartyCount":17,
			 "GvgCastleState":2,
			 "UtcFallenTimeStamp":1702070527101
		  }
	   ],
	   "guilds":{
		  "129044010099":"星喰",
		  "196765134099":"メメ温泉99番地",
		  "199652669099":"サリンジャー",
		  "394114239099":"REsKend",
		  "404400524099":"引きちぎられたヒレ",
		  "482173199099":"黒魔女のお茶会",
		  "494634944099":"凛として時雨",
		  "530715127099":"9/POSTER",
		  "664163737099":"アリウススクワッド",
		  "669198682099":"初心者魔女の喫茶店",
		  "719463286099":"ウロボロス鎮魂歌",
		  "720141187099":"魔法の琥珀",
		  "741934801099":"電気羊の夢を見る",
		  "768981549099":"カルデア",
		  "769575323099":"ああああ",
		  "816626613099":"独立国家わちゃわちゃ",
		  "908515237099":"戒めの赤袋爺"
	   }
	}
 }
*/
