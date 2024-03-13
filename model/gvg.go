package model

import "time"

//castle_id received via webocket
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
	Changed            bool
}

//guild_id received via webocket
type GuildId struct {
	Value  Value_GuildId
	Offset int
}
type Value_GuildId struct {
	StreamId  Value_StreamId
	GuildId   int
	GuildName string
	Changed   bool
}

//stream_id received via webocket
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

// https://itpfdoc.hitachi.co.jp/manuals/3000/30003F5600/EEXR0016.HTM
//registration failed guild / for websocket display
type FailedGuild struct {
	WorldId    int       `json:"world_id"`
	GuildId    int       `json:"guild_id"`
	GuildName  string    `json:"guild_name"`
	CreateDate time.Time `json:"create_date"`
}

//registration failed record / for websocket display
type FailedRecord struct {
	WorldId            int       `json:"world_id"`
	GroupId            int       `json:"group_id"`
	Class              int       `json:"class"`
	Block              int       `json:"block"`
	CastleId           int       `json:"castle_id"`
	GuildId            int       `json:"guild_id"`
	AttackerGuildId    int       `json:"atk_guild_id"`
	UtcFallenTimeStamp int       `json:"fallen_time"`
	DefensePartyCount  int       `json:"def_pty_count"`
	AttackPartyCount   int       `json:"atk_pty_count"`
	GvgCastleState     int       `json:"state"`
	CreateDate         time.Time `json:"create_date"`
}

//current setting for streaming
type StreamConf struct {
	Country string
	World   string
	Group   string
	Class   string
	Block   string
	Castle  string
}
