package analysis

import (
	"meme/model"
)

// guild objects
var Guilds = map[int]*model.Value_GuildId{}

// record objects
var Castles = map[int]*model.Value_CastleId{}

// parsing webhook binary data
func GvgAnalysis(u8arr []byte, current_sub model.Value_StreamId) {
	offset := 0

	for offset < len(u8arr) {
		//extract stream
		res_stream := unpack_stream_id(u8arr, offset)
		stream_id := res_stream.Value
		offset = res_stream.Offset

		if stream_id.CastleId == 0 {
			//extract guild
			res_guild := unpack_guild(u8arr, offset, stream_id.WorldId)
			guild := res_guild.Value
			guild.StreamId = stream_id
			offset = res_guild.Offset

			if !check_stream_compat(stream_id, current_sub) {
				continue
			}

			if guild.GuildId == 0 {
				//init list of guilds
				Guilds = map[int]*model.Value_GuildId{}
			} else {
				//register guild
				Guilds[guild.GuildId] = &guild
			}
		} else {
			//extract castle
			res_castle := unpack_castle(u8arr, offset, stream_id.WorldId)
			castle := res_castle.Value
			castle.StreamId = stream_id
			offset = res_castle.Offset

			if !check_stream_compat(stream_id, current_sub) {
				continue
			}

			//register castle
			Castles[stream_id.CastleId-1] = &castle
		}
	}
}

// packing stream_id for websocket
func PackStreamId(current_sub model.Value_StreamId) []byte {
	return pack_stream_id(0, current_sub)
}

// check if subject has changed
func check_stream_compat(x model.Value_StreamId, y model.Value_StreamId) bool {
	if (x.GroupId == 0 && x.Class == 0 && x.Block == 0) ||
		(y.GroupId == 0 && y.Class == 0 && y.Block == 0) {
		return x.WorldId == y.WorldId
	} else {
		return x.GroupId == y.GroupId && x.Class == y.Class && x.Block == y.Block
	}
}

// pack stream_id into binary
func pack_le32(offset int, sid int) (u8arr []byte) {
	u8arr = make([]byte, 4)
	u8arr[offset] = byte(sid & 0xFF)
	u8arr[offset+1] = byte((sid >> 8) & 0xFF)
	u8arr[offset+2] = byte((sid >> 16) & 0xFF)
	u8arr[offset+3] = byte((sid >> 24) & 0xFF)

	return u8arr
}

// extract stream_id from binary data
func pack_stream_id(offset int, stream_id model.Value_StreamId) []byte {
	sid := stream_id.WorldId<<19 | stream_id.Class<<16 | stream_id.GroupId<<8 | stream_id.Block<<5 | stream_id.CastleId

	return pack_le32(offset, sid)
}

// extract castle_id from binary data
func unpack_castle(u8arr []byte, offset int, world int) model.CastleId {
	def_guild := unpack_le32(u8arr, offset)
	atk_guild := unpack_le32(u8arr, offset+4)
	fallen := unpack_le32(u8arr, offset+8) * 1000
	def_count := unpack_le16(u8arr, offset+12)
	atk_count := unpack_le16(u8arr, offset+14)
	state := u8arr[offset+16]

	return model.CastleId{
		Value: model.Value_CastleId{
			GuildId:            def_guild*1000 + world%1000,
			AttackerGuildId:    atk_guild*1000 + world%1000,
			UtcFallenTimeStamp: fallen * 1000,
			DefensePartyCount:  def_count,
			AttackPartyCount:   atk_count,
			GvgCastleState:     int(state),
			Changed:            true,
		},
		Offset: offset + 20,
	}
}

// extract guild_id from binary data
func unpack_guild(u8arr []byte, offset int, world int) model.GuildId {
	gid := unpack_le32(u8arr, offset)
	len := int(u8arr[offset+4])
	str := string(u8arr[offset+5 : offset+5+len])

	return model.GuildId{
		Value: model.Value_GuildId{
			GuildId:   gid*1000 + world%1000,
			GuildName: str,
			Changed:   true,
		},
		Offset: offset + 5 + len,
	}
}

// extraction processing core 16
func unpack_le16(u8arr []byte, offset int) int {
	return int(u8arr[offset]) | int(u8arr[offset+1])<<8
}

// extraction processing core 32
func unpack_le32(u8arr []byte, offset int) int {
	return int(u8arr[offset]) | int(u8arr[offset+1])<<8 | int(u8arr[offset+2])<<16 | int(u8arr[offset+3])<<24
}

// extract stream_id from binary data
func unpack_stream_id(u8arr []byte, offset int) model.StreamId {
	s := unpack_le32(u8arr, offset)

	return model.StreamId{
		Value: model.Value_StreamId{
			WorldId:  s >> 19,
			GroupId:  (s >> 8) & 0xFF,
			Class:    (s >> 16) & 0x7,
			Block:    (s >> 5) & 0x7,
			CastleId: s & 0x1F,
		},
		Offset: offset + 4,
	}
}
