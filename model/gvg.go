package model

import "time"

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
	Changed   bool
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
	Changed            bool
}

// https://itpfdoc.hitachi.co.jp/manuals/3000/30003F5600/EEXR0016.HTM
type FailedGuild struct {
	WorldId    int
	GuildId    int
	GuildName  string
	CreateDate time.Time
}

type FailedRecord struct {
	WorldId            int
	GroupId            int
	Class              int
	Block              int
	CastleId           int
	GuildId            int
	AttackerGuildId    int
	UtcFallenTimeStamp int
	DefensePartyCount  int
	AttackPartyCount   int
	GvgCastleState     int
	CreateDate         time.Time
}

type StreamConf struct {
	Country string
	World   string
	Group   string
	Class   string
	Block   string
	Castle  string
}
